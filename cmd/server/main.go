package main

import (
	"log"
	"voldy/pkg/cluster"
	"voldy/pkg/server"
)

func main() {

	nodes := []*cluster.Node{
		cluster.NewNode(1, "localhost", 9000, 9000, 9000, []int{0}),
	}
	currentCluster, err := cluster.NewCluster("Test", nodes, []*cluster.Zone{})
	if err != nil {
		log.Fatal(err)
	}

	srv, err := server.NewGRPCServer(server.GRPCServerConfig{
		Host:     "localhost",
		Port:     9000,
		Replicas: 1,
	}, currentCluster)

	if err := srv.Listen(); err != nil {
		log.Fatal(err)
	}
}