package store

import (
	"context"
	"voldy/pkg/versioning"
)

type Entry struct {
	HashKey []byte
	SortKey []byte
	Version *versioning.Versioned
}

// StorageEngine defines the methods that an implemented storage engine must support
type StorageEngine interface {
	CreateTable(ctx context.Context, tableName []byte) error
	// TODO - Replace with an entry as response
	Get(ctx context.Context, tableName, hashKey, sortKey []byte) (*versioning.Versioned, error)
	GetHashKey(ctx context.Context, tableName, hashKey []byte) ([]Entry, error)
	Put(ctx context.Context, tableName, hashKey, SortKey []byte, version *versioning.Versioned) error
	PutKeys(ctx context.Context, tableName []byte, entries []Entry) error
	// TODO - Delete
	// TODO - DropTable
	// TODO - Keys()
	// TODO - Iterator
}
