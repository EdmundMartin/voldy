package routing

import (
	"errors"
	"hash/fnv"
	"voldy/pkg/cluster"
)

const (
	RouteConsistent = "consistent-routing"
)

type ConsistentRoutingStrategy struct {
	numberReplicas int
	nodes          []*cluster.Node
}

func NewConsistentRoutingStrategy(nodes []*cluster.Node, numReplicas int) (*ConsistentRoutingStrategy, error) {
	cs := &ConsistentRoutingStrategy{
		numberReplicas: numReplicas,
	}

	partToNode := map[int]*cluster.Node{}
	for _, node := range nodes {
		for _, part := range node.Partitions {
			_, ok := partToNode[part]
			if ok {
				return nil, errors.New("duplicate partition id")
			}
			partToNode[part] = node
		}
	}
	cs.nodes = make([]*cluster.Node, len(partToNode))
	for i := 0; i < len(partToNode); i++ {
		val, ok := partToNode[i]
		if !ok {
			return nil, errors.New("invalid configuration missing partition")
		}
		cs.nodes[i] = val
	}
	return cs, nil
}

func (c ConsistentRoutingStrategy) Type() string {
	return RouteConsistent
}

func (c ConsistentRoutingStrategy) RouteRequest(key []byte) []*cluster.Node {
	partitionList := c.GetPartitionList(key)

	if len(partitionList) == 0 {
		return []*cluster.Node{}
	}
	preferenceList := make([]*cluster.Node, len(partitionList))

	for idx, partition := range partitionList {
		preferenceList[idx] = c.nodes[partition]
	}
	return preferenceList
}

func (c ConsistentRoutingStrategy) GetPartitionList(key []byte) []int {
	idx := c.GetMasterPartition(key)
	return c.GetReplicatingPartitionList(idx)
}

func (c ConsistentRoutingStrategy) GetMasterPartition(key []byte) int {
	hash := fnv.New64()
	hash.Write(key)
	return int(hash.Sum64() % uint64(len(c.nodes)))
}

func (c ConsistentRoutingStrategy) GetReplicatingPartitionList(partition int) []int {
	preferenceList := make([]*cluster.Node, c.numberReplicas)
	replicationPartitionsList := make([]int, c.numberReplicas)

	if len(c.nodes) == 0 {
		return []int{}
	}

	for i := 0; i < len(c.nodes); i++ {
		if !listContains(preferenceList, c.nodes[partition]) {
			preferenceList = append(preferenceList, c.nodes[partition])
			replicationPartitionsList = append(replicationPartitionsList, partition)
		}

		if len(preferenceList) >= c.numberReplicas {
			return replicationPartitionsList
		}
		partition = (partition + 1) % len(c.nodes)
	}

	return replicationPartitionsList
}

func (c ConsistentRoutingStrategy) GetNodes() []*cluster.Node {
	result := make([]*cluster.Node, len(c.nodes))
	for i := 0; i < len(c.nodes); i++ {
		result[i] = c.nodes[i]
	}
	return result
}

func (c ConsistentRoutingStrategy) GetReplicaCount() int {
	return c.numberReplicas
}

func listContains(nodes []*cluster.Node, search *cluster.Node) bool {
	for _, n := range nodes {
		if n == search {
			return true
		}
	}
	return false
}
