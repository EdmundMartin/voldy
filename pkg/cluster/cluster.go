package cluster

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
)

type Cluster struct {
	// TODO - A lot of these fields should probably be private
	Name              string
	NumberPartitions  int
	NodesById         map[int]*Node
	ZonesById         map[int]*Zone
	NodesPerZone      map[*Zone][]int
	PartitionsPerZone map[*Zone][]int

	PartitionIdToZone      map[int]*Zone
	PartitionIdToNodeArray []*Node
	PartitionIdToNode      map[int]*Node
	PartitionIdToNodeId    map[int]int
	ShuffledNodes          []*Node
}

func NewCluster(name string, nodes []*Node, zones []*Zone) (*Cluster, error) {

	cl := &Cluster{
		Name: name,
	}
	// TODO place these in Cluster
	cl.PartitionsPerZone = map[*Zone][]int{}
	cl.NodesPerZone = map[*Zone][]int{}
	cl.PartitionIdToZone = map[int]*Zone{}
	cl.PartitionIdToNode = map[int]*Node{}
	cl.PartitionIdToNodeId = map[int]int{}

	partitionIdToNodeMap := map[int]*Node{}

	zonesById := map[int]*Zone{}
	if len(zones) != 0 {
		for _, zone := range zones {
			_, ok := zonesById[zone.ZoneId]
			if ok {
				return nil, errors.New(fmt.Sprintf("zone id %d appears twice in zone list", zone.ZoneId))
			}
			zonesById[zone.ZoneId] = zone
			cl.NodesPerZone[zone] = []int{}
			cl.PartitionsPerZone[zone] = []int{}
		}
	} else {
		// Add a default zone
		zone := NewDefaultZone()
		zonesById[zone.ZoneId] = zone
		cl.NodesPerZone[zone] = []int{}
		cl.PartitionsPerZone[zone] = []int{}
	}

	cl.NodesById = map[int]*Node{}

	for _, node := range nodes {
		_, ok := cl.NodesById[node.Id]
		if ok {
			return nil, errors.New(fmt.Sprintf("node id %d appears twice in node list", node.Id))
		}
		cl.NodesById[node.Id] = node

		zone, ok := zonesById[node.ZoneId]
		if !ok {
			return nil, errors.New("no zone exists associated with zone")
		}

		val, ok := cl.NodesPerZone[zone]
		if !ok {
			cl.NodesPerZone[zone] = []int{node.Id}
		} else {
			val = append(val, node.Id)
			cl.NodesPerZone[zone] = val
		}

		val, ok = cl.PartitionsPerZone[zone]
		if !ok {
			cl.PartitionsPerZone[zone] = node.Partitions
		} else {
			val = append(val, node.Partitions...)
			cl.PartitionsPerZone[zone] = val
		}

		for _, partitionId := range node.Partitions {
			_, ok = cl.PartitionIdToNodeId[partitionId]
			if ok {
				return nil, fmt.Errorf("partition id found on two nodes, partition id:%d , node id: %d", partitionId, node.Id)
			}
			cl.PartitionIdToZone[partitionId] = zone
			partitionIdToNodeMap[partitionId] = node
			cl.PartitionIdToNode[partitionId] = node
			cl.PartitionIdToNodeId[partitionId] = node.Id
		}
	}
	count, err := countPartitions(nodes)
	if err != nil {
		return nil, err
	}
	cl.NumberPartitions = count

	cl.PartitionIdToNodeArray = make([]*Node, cl.NumberPartitions)
	for i := 0; i < cl.NumberPartitions; i++ {
		cl.PartitionIdToNodeArray[i] = partitionIdToNodeMap[i]
	}

	var allNodes []*Node
	for _, v := range cl.NodesById {
		allNodes = append(allNodes, v)
	}
	rand.Shuffle(len(allNodes), func(i, j int) {
		allNodes[i], allNodes[j] = allNodes[j], allNodes[i]
	})
	cl.ShuffledNodes = allNodes
	return cl, nil
}

func countPartitions(nodes []*Node) (int, error) {
	tags := []int{}
	for _, node := range nodes {
		tags = append(tags, node.Partitions...)
	}
	sort.Ints(tags)

	for i := 0; i < len(tags); i++ {
		if tags[i] != i {
			return 0, errors.New("invalid partition assignment")
		}
	}
	return len(tags), nil
}

func (cl *Cluster) GetNodes() []*Node {
	result := make([]*Node, len(cl.NodesById))
	idx := 0
	for _, val := range cl.NodesById {
		result[idx] = val
		idx++
	}
	return result
}

func (cl *Cluster) GetNodeIds() []int {
	result := make([]int, len(cl.NodesById))
	idx := 0
	for k, _ := range cl.NodesById {
		result[idx] = k
		idx++
	}
	return result
}

func (cl *Cluster) GetZoneIds() []int {
	result := make([]int, len(cl.ZonesById))
	idx := 0
	for k, _ := range cl.ZonesById {
		result[idx] = k
		idx++
	}
	return result
}

func (cl *Cluster) GetZoneById(id int) (*Zone, error) {
	zone, ok := cl.ZonesById[id]
	if !ok {
		return nil, fmt.Errorf("no zone with id: %d", id)
	}
	return zone, nil
}

func (cl *Cluster) NumberOfZones() int {
	return len(cl.ZonesById)
}

// TODO - There are bunch more methods on the cluster required
func (cl *Cluster) NumberPartitionsInZone(id int) int {
	return 0
}
