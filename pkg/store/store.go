package store

import "voldy/pkg/versioning"

type TransformFunction func([]byte) []byte

type Pair struct {
	Key   []byte
	Value []*versioning.Versioned[[]byte]
}

type Store interface {
	Get(key []byte, transform TransformFunction) ([]*versioning.Versioned[[]byte], error)
	GetAll(keys [][]byte, transform TransformFunction) (map[string][]*versioning.Versioned[[]byte], error)
	Put(key []byte, versioned *versioning.Versioned[[]byte], transform TransformFunction) error
	Delete(key []byte, version *versioning.VectorClock) (bool, error)
	GetName() string
	Close() error
	GetCapability() // TODO - Types
	GetVersions(key []byte) ([]versioning.Version, error)
}
