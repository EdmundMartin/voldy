package bengine

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
	"voldy/pkg/store"
	"voldy/pkg/versioning"
)

func TestStoreEngine_BasicFlow(t *testing.T) {

	filepath := "test.db"
	defer cleanUp(filepath)
	ctx := context.Background()
	s, err := NewStorageEngine(filepath, []byte("#"))
	require.NoError(t, err)

	tableName := []byte("TestTable")
	err = s.CreateTable(ctx, tableName)
	assert.NoError(t, err)

	hashKey := []byte("Tx_100")

	sortKey := []byte("2022-02-24")

	vectorC, err := versioning.NewEmptyClock().Incremented(1, time.Now().Unix())
	require.NoError(t, err)
	err = s.Put(ctx, tableName, hashKey, sortKey, &versioning.Versioned{
		Version:  vectorC,
		Contents: []byte("Test transaction"),
	})
	require.NoError(t, err)

	result, err := s.Get(ctx, tableName, hashKey, sortKey)
	assert.NoError(t, err)

	assert.Equal(t, []byte("Test transaction"), result.Contents)
}

func TestStoreEngine_GetHashKey(t *testing.T) {

	filepath := "hashKey.db"
	defer cleanUp(filepath)
	ctx := context.Background()
	s, err := NewStorageEngine(filepath, []byte("#"))
	require.NoError(t, err)

	tableName := []byte("TestTable")
	err = s.CreateTable(ctx, tableName)
	require.NoError(t, err)

	hashKey := []byte("Tx_100")

	var entries []store.Entry
	for idx := 0; idx < 5; idx++ {
		sortKey := fmt.Sprintf("%d", idx)

		clock, _ := versioning.NewEmptyClock().Incremented(1, time.Now().Unix())
		entries = append(entries, store.Entry{
			HashKey: hashKey,
			SortKey: []byte(sortKey),
			Version: versioning.NewVersionedBytes([]byte(sortKey), clock),
		})
	}

	err = s.PutKeys(ctx, tableName, entries)
	assert.NoError(t, err)

	result, err := s.GetHashKey(ctx, tableName, hashKey)
	assert.Len(t, result, 5)
	for i := 0; i < len(result); i++ {
		val := []byte(fmt.Sprintf("%d", i))
		assert.Equal(t, val, result[i].Version.Contents)
		assert.Equal(t, val, result[i].SortKey)
	}

}

func cleanUp(filename string) {
	os.Remove(filename)
}
