package routing

import "voldy/pkg/cluster"

type Strategy interface {
	Type() string
	RouteRequest(key []byte) []*cluster.Node
	GetPartitionList(key []byte) []int
	GetMasterPartition(key []byte) int
	GetReplicatingPartitionList(partition int) []int
	GetNodes() []*cluster.Node
	GetReplicaCount() int
}
