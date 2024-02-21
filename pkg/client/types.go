package client

import "voldy/pkg/versioning"

type GetResponse struct {
	Key     []byte
	Value   []byte
	Version *versioning.VectorClock
}
