package routing

import "voldy/pkg/cluster"

const (
	// RouteToAllPartitionId Use partition ID of for all keys in an implement to all strategy
	RouteToAllPartitionId = 0
	RouteToAllName        = "all-routing"
)

type RouteToAllStrategy struct {
	Nodes        []*cluster.Node
	PartitionIds []int
}

func NewRouteToAllStrategy(nodes []*cluster.Node) *RouteToAllStrategy {
	return &RouteToAllStrategy{Nodes: nodes, PartitionIds: []int{RouteToAllPartitionId}}
}

func (r RouteToAllStrategy) Type() string {
	return RouteToAllName
}

func (r RouteToAllStrategy) RouteRequest(key []byte) []*cluster.Node {
	return r.Nodes
}

func (r RouteToAllStrategy) GetPartitionList(key []byte) []int {
	return r.PartitionIds
}

func (r RouteToAllStrategy) GetMasterPartition(key []byte) int {
	return RouteToAllPartitionId
}

func (r RouteToAllStrategy) GetReplicatingPartitionList(partition int) []int {
	return r.PartitionIds
}

func (r RouteToAllStrategy) GetNodes() []*cluster.Node {
	return r.Nodes
}

func (r RouteToAllStrategy) GetReplicaCount() int {
	return len(r.Nodes)
}
