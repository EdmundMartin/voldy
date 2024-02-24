package server

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"time"
	"voldy/pkg/cluster"
	"voldy/pkg/protocol"
	"voldy/pkg/routing"
	"voldy/pkg/store"
	"voldy/pkg/store/bengine"
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

func (g *GRPCServer) CreateTable(ctx context.Context, request *protocol.CreateTableRequest) (*protocol.CreateTableResponse, error) {

	nodes := g.strategy.GetNodes()
	for _, node := range nodes {
		if g.Config.Host == node.Host && g.Config.Port == node.GrpcPort {
			if err := g.storageEngine.CreateTable(request.TableName); err != nil {
				return nil, err
			}
		}
		// TODO - Communicate with other nodes -- all nodes most accept a create table request
	}
	return &protocol.CreateTableResponse{TableName: request.TableName}, nil
}

// TODO - Implement read-repair
func (g *GRPCServer) GetKey(ctx context.Context, request *protocol.GetKeyRequest) (*protocol.GetResponse, error) {

	nodes := g.strategy.RouteRequest(request.HashKey)

	var results []*versioning.Versioned
	for _, node := range nodes {
		if g.Config.Host == node.Host && g.Config.Port == node.GrpcPort {
			ver, err := g.storageEngine.Get(request.TableName, request.HashKey, request.SortKey)
			if err != nil {
				// TODO - Maybe we should gather and ignore errors
				return nil, err
			}
			results = append(results, ver)
		}
	}
	if len(results) == 0 {
		return &protocol.GetResponse{
			Key:   request.HashKey,
			Value: []byte{},
		}, nil
	}

	result := results[len(results)-1]

	return &protocol.GetResponse{
		Key:     request.HashKey,
		Value:   result.Contents,
		Version: result.Version.ToBytes(),
	}, nil
}

func (g *GRPCServer) Put(ctx context.Context, request *protocol.PutRequest) (*protocol.PutResponse, error) {
	nodes := g.strategy.RouteRequest(request.HashKey)

	res, err := g.storageEngine.Get(request.TableName, request.HashKey, request.SortKey)
	if err != nil {
		return nil, err
	}
	var vectorClock *versioning.VectorClock
	if res == nil {
		vectorClock = versioning.NewEmptyClock()
		if err := vectorClock.IncrementVersion(g.ourNode.Id, time.Now().Unix()); err != nil {
			return nil, err
		}
	} else {
		vectorClock = res.Version
	}
	if request.Version != nil {
		otherClock := versioning.VectorClockFromBytes(request.Version)
		occurred, err := vectorClock.Compare(otherClock)
		if err != nil {
			return nil, err
		}
		// we have obsolete version stored on this node we should increment the greater vector clock
		if occurred == versioning.AFTER {
			return nil, errors.New("obsolete version of object presented for saving")
		}
	}
	if err := vectorClock.IncrementVersion(1, time.Now().Unix()); err != nil {
		return nil, err
	}
	for _, node := range nodes {
		if g.Config.Host == node.Host && g.Config.Port == node.GrpcPort {
			err := g.storageEngine.Put(request.TableName, request.HashKey, request.SortKey, versioning.NewVersionedBytes(request.Value, vectorClock))
			if err != nil {
				return nil, err
			}
		} else {
			// TODO Try and put on other servers - if we fail because of other version being hire - fail request
			fmt.Println("Should attempt to send request to other server")
		}
	}
	return &protocol.PutResponse{
		Key:     request.HashKey,
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
	engine, err := bengine.NewStorageEngine("Demo.db", []byte("#"))
	if err != nil {
		return nil, err
	}
	server.storageEngine = engine

	return server, nil
}
