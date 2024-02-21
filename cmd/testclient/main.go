package main

import (
	"context"
	"fmt"
	"log"
	"voldy/pkg/client"
)

func main() {

	cl, err := client.NewClient(":9000")
	if err != nil {
		log.Fatal(err)
		return
	}

	err = cl.Put(context.Background(), []byte("Katie"), []byte("Is Great"))
	fmt.Println(err)

	err = cl.Put(context.Background(), []byte("Katie"), []byte("Is Good"))

	fmt.Println(err)

	resp, err := cl.Get(context.Background(), []byte("Katie"))

	fmt.Println(string(resp.Key))
	fmt.Println(string(resp.Value))
	fmt.Println(err)
}
