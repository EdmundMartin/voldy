package cluster

import (
	"fmt"
	"strings"
)

type Node struct {
	SerialVersion int
	Id            int
	Host          string
	GrpcPort      int
	SocketPort    int
	AdminPort     int
	ZoneId        int
	Partitions    []int
	RestPort      int
}

type NodeOptions func(node *Node)

func NewNode(id int, host string, grpcPort, socketPort, adminPort int, partitions []int, opts ...NodeOptions) *Node {

	node := &Node{
		SerialVersion: 1,
		Id:            id,
		Host:          host,
		GrpcPort:      grpcPort,
		SocketPort:    socketPort,
		AdminPort:     adminPort,
		ZoneId:        0,
		Partitions:    partitions,
		RestPort:      -1,
	}

	for _, opt := range opts {
		opt(node)
	}

	if adminPort == -1 {
		node.AdminPort = node.SocketPort + 1
	}

	return node
}

func (n *Node) NumberPartitions() int {
	return len(n.Partitions)
}

func (n *Node) GrpcUrl() string {
	return fmt.Sprintf("tcp://%s:%d", n.Host, n.GrpcPort)
}

func (n *Node) SocketUrl() string {
	return fmt.Sprintf("tcp://%s:%d", n.Host, n.SocketPort)
}

func (n *Node) String() string {
	return fmt.Sprintf("<Node Host: %s>, SocketPort: %d, Id: %d>", n.Host, n.SocketPort, n.Id)
}

func (n *Node) Equals(other *Node) bool {
	return n.Id == other.Id
}

func (n *Node) HasSameState(other *Node) bool {
	if n.Id != other.Id {
		return false
	}
	if strings.ToLower(n.Host) != strings.ToLower(other.Host) {
		return false
	}

	if n.GrpcPort != other.GrpcPort {
		return false
	}

	if n.SocketPort != other.SocketPort {
		return false
	}

	if n.AdminPort != other.AdminPort {
		return false
	}

	if n.ZoneId != other.ZoneId {
		return false
	}

	return true
}

type Nodes []*Node

func (n Nodes) Len() int {
	return len(n)
}

func (n Nodes) Less(i, j int) bool {
	return n[i].Id < n[j].Id
}

func (n Nodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}
