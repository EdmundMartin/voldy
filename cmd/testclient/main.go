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

	ctx := context.Background()
	cl.CreateTable(ctx, "TestTable")

	err = cl.Put(ctx, "TestTable", []byte("Tx_101"), []byte("2022-02-24"), []byte("Hello Transaction"))

	resp, err := cl.GetKey(ctx, "TestTable", []byte("Tx_101"), []byte("2022-02-24"))

	err = cl.Put(ctx, "TestTable", []byte("Tx_101"), []byte("2022-02-24"), []byte("Hello World"))

	resp, err = cl.GetKey(ctx, "TestTable", []byte("Tx_101"), []byte("2022-02-24"))

	fmt.Println(string(resp.Contents))
	fmt.Println(resp.Version)
}
