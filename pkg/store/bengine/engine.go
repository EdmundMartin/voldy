package bengine

import (
	"bytes"
	"context"
	bolt "go.etcd.io/bbolt"
	"voldy/pkg/store"
	"voldy/pkg/versioning"
)

type StoreEngine struct {
	db        *bolt.DB
	separator []byte
}

func (s StoreEngine) GetHashKey(ctx context.Context, tableName, hashKey []byte) ([]store.Entry, error) {
	var result []store.Entry
	err := s.db.View(func(tx *bolt.Tx) error {
		cursor := tx.Bucket(tableName).Cursor()
		for k, v := cursor.Seek(hashKey); k != nil && bytes.HasPrefix(k, hashKey); k, v = cursor.Next() {
			value := versioning.VersionedFromBytes(v)
			pk, sk := splitKey(k, s.separator)
			result = append(result, store.Entry{
				HashKey: pk,
				SortKey: sk,
				Version: value,
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s StoreEngine) CreateTable(ctx context.Context, tableName []byte) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucket(tableName); err != nil {
			return err
		}
		return nil
	})
}

func (s StoreEngine) Get(ctx context.Context, tableName, HashKey, SortKey []byte) (*versioning.Versioned, error) {
	var result *versioning.Versioned
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(tableName)
		contents := bucket.Get(fullyJustifiedKey(HashKey, SortKey, s.separator))
		if len(contents) == 0 {
			return nil
		}
		result = versioning.VersionedFromBytes(contents)
		return nil
	})
	return result, err
}

func (s StoreEngine) Put(ctx context.Context, tableName, hashKey, sortKey []byte, version *versioning.Versioned) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(tableName)
		justifiedKey := fullyJustifiedKey(hashKey, sortKey, s.separator)
		return bucket.Put(justifiedKey, version.VersionedToBytes())
	})
}

func (s StoreEngine) PutKeys(ctx context.Context, tableName []byte, entries []store.Entry) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(tableName)
		for _, entry := range entries {
			justifiedKey := fullyJustifiedKey(entry.HashKey, entry.SortKey, s.separator)
			if err := bucket.Put(justifiedKey, entry.Version.VersionedToBytes()); err != nil {
				return err
			}
		}
		return nil
	})
}

func NewStorageEngine(filepath string, keySeparator []byte) (*StoreEngine, error) {
	db, err := bolt.Open(filepath, 0600, nil)
	if err != nil {
		return nil, err
	}
	s := &StoreEngine{
		db:        db,
		separator: keySeparator,
	}
	return s, nil
}

func fullyJustifiedKey(hashKey, sortKey, sep []byte) []byte {
	return bytes.Join([][]byte{
		hashKey,
		sortKey,
	}, sep)
}

func splitKey(hashKey, sep []byte) ([]byte, []byte) {
	result := bytes.Split(hashKey, sep)
	return result[0], result[1]
}
