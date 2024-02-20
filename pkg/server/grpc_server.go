package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"voldy/pkg/cluster"
	"voldy/pkg/protocol"
	"voldy/pkg/routing"
	"voldy/pkg/store"
	"voldy/pkg/store/memory"
	"voldy/pkg/versioning"
)

type GRPCServer struct {
	Config        GRPCServerConfig
	cluster       *cluster.Cluster
	storageEngine store.StorageEngine
	strategy      routing.Strategy
	ourNode       *cluster.Node
	otherServers  map[*cluster.Node]protocol.VoldyClient
	protocol.UnimplementedVoldyServer
}

type GRPCServerConfig struct {
	Host     string
	Port     int
	Replicas int
}

func (g *GRPCServer) Get(ctx context.Context, request *protocol.GetRequest) (*protocol.GetResponse, error) {

	nodes := g.strategy.RouteRequest(request.Key)

	var results []*versioning.Versioned[[]byte]
	for _, node := range nodes {
		if g.Config.Host == node.Host && g.Config.Port == node.GrpcPort {
			ver, err := g.storageEngine.Get(request.Key, nil)
			if err != nil {
				// TODO - Maybe we should gather and ignore errors
				return nil, err
			}
			results = append(results, ver...)
		}
	}
	if len(results) == 0 {
		return &protocol.GetResponse{
			Key:     request.Key,
			Message: []byte{},
		}, nil
	}
	result := results[len(results)-1]

	return &protocol.GetResponse{
		Key:     request.Key,
		Message: result.Contents,
		Version: result.Version.ToBytes(),
	}, nil
}

func (g *GRPCServer) Put(ctx context.Context, request *protocol.PutRequest) (*protocol.PutResponse, error) {
	nodes := g.strategy.RouteRequest(request.Key)

	for _, node := range nodes {
		if g.Config.Host == node.Host && g.Config.Port == node.GrpcPort {
			err := g.storageEngine.Put(request.Key, versioning.NewVersionedBytes(request.Value, nil), nil)
			if err != nil {
				return nil, err
			}
		} else {
			// TODO
			fmt.Println("Should attempt to send request to other server")
		}
	}
	return &protocol.PutResponse{
		Key:     request.Key,
		Value:   request.Value,
		Version: nil, // TODO - Remove
	}, nil
}

func (g *GRPCServer) Listen() error {
	serv := grpc.NewServer()
	protocol.RegisterVoldyServer(serv, g)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", g.Config.Port))
	if err != nil {
		return err
	}
	if err := serv.Serve(lis); err != nil {
		return err
	}

	return nil
}

func NewGRPCServer(config GRPCServerConfig, currentCluster *cluster.Cluster) (*GRPCServer, error) {

	server := &GRPCServer{
		cluster: currentCluster,
		Config:  config,
	}

	var allNodes []*cluster.Node
	for _, node := range server.cluster.NodesById {
		allNodes = append(allNodes, node)
		if config.Host == node.Host && config.Port == node.GrpcPort {
			server.ourNode = node
			continue
		}
		// TODO - Full dial of host and port
		conn, err := grpc.Dial(fmt.Sprintf(":%d", node.GrpcPort), grpc.WithInsecure())
		if err != nil {
			server.otherServers[node] = nil
			continue
		}
		client := protocol.NewVoldyClient(conn)
		server.otherServers[node] = client
	}
	// TODO - Make strategy configurable
	strat, err := routing.NewConsistentRoutingStrategy(allNodes, 0)
	if err != nil {
		return nil, err
	}
	server.strategy = strat

	// TODO - Make storage engine configurable
	server.storageEngine = memory.NewInMemoryStorageEngine("Test", server.ourNode.Id)

	return server, nil
}
