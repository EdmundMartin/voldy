package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"voldy/pkg/protocol"
)

func main() {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	defer conn.Close()

	srv := protocol.NewVoldyClient(conn)

	response, err := srv.Put(context.Background(), &protocol.PutRequest{
		Key:   []byte("Katie"),
		Value: []byte("Is Great"),
	})

	fmt.Println(response)
	fmt.Println(err)

	resp, err := srv.Get(context.Background(), &protocol.GetRequest{
		Key: []byte("Katie"),
	})

	fmt.Println(resp)
	fmt.Println(err)
}