package store

import "voldy/pkg/versioning"

// StorageEngine will replace the existing interface for a storage engine
type StorageEngine interface {
	CreateTable(tableName []byte) error
	Get(tableName, hashKey, sortKey []byte) (*versioning.Versioned, error)
	GetHashKey(tableName, hashKey []byte) ([]*versioning.Versioned, error)
	Put(tableName, hashKey, SortKey []byte, version *versioning.Versioned) error
}
