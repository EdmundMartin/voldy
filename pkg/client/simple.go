package client

import (
	"context"
	"google.golang.org/grpc"
	"voldy/pkg/protocol"
	"voldy/pkg/versioning"
)

type Configuration struct {
}

type SimpleClient struct {
	client protocol.VoldyClient
	close  func() error
}

func NewClient(host string) (*SimpleClient, error) {
	c := &SimpleClient{}

	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	c.close = conn.Close
	srv := protocol.NewVoldyClient(conn)
	c.client = srv
	return c, nil
}

func (c *SimpleClient) CreateTable(ctx context.Context, tableName string) error {
	if _, err := c.client.CreateTable(ctx, &protocol.CreateTableRequest{TableName: []byte(tableName)}); err != nil {
		return err
	}
	return nil
}

func (c *SimpleClient) GetKey(ctx context.Context, tableName string, hashKey, sortKey []byte) (*versioning.Versioned, error) {

	res, err := c.client.GetKey(ctx, &protocol.GetKeyRequest{
		TableName: []byte(tableName),
		HashKey:   hashKey,
		SortKey:   sortKey,
	})
	if err != nil {
		return nil, err
	}

	return &versioning.Versioned{
		Version:  versioning.VectorClockFromBytes(res.Version),
		Contents: res.Value,
	}, nil
}

func (c *SimpleClient) Put(ctx context.Context, tableName string, hashKey, sortKey, value []byte) error {

	_, err := c.client.Put(ctx, &protocol.PutRequest{
		TableName: []byte(tableName),
		HashKey:   hashKey,
		SortKey:   sortKey,
		Value:     value,
		Version:   nil,
	})
	return err
}

func (c *SimpleClient) Close() error {
	return c.close()
}
