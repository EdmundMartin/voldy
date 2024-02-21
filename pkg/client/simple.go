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

func (c *SimpleClient) Get(ctx context.Context, key []byte) (*GetResponse, error) {
	res, err := c.client.Get(ctx, &protocol.GetRequest{Key: key})
	if err != nil {
		return nil, err
	}
	response := &GetResponse{
		Key:   res.Key,
		Value: res.Value,
	}
	if len(res.Version) != 0 {
		response.Version = versioning.VectorClockFromBytes(res.Version)
	}

	return response, nil
}

func (c *SimpleClient) Put(ctx context.Context, key []byte, value []byte) error {

	result, err := c.Get(ctx, key)
	if err != nil {
		return nil
	}
	var vectorClock *versioning.VectorClock
	if result.Version == nil {
		vectorClock = versioning.NewEmptyClock()
	} else {
		vectorClock = result.Version
	}
	_, err = c.client.Put(ctx, &protocol.PutRequest{
		Key:     key,
		Value:   value,
		Version: vectorClock.ToBytes(),
	})
	return err
}

func (c *SimpleClient) Close() error {
	return c.close()
}
