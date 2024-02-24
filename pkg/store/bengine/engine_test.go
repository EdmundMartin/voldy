package bengine

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
	"voldy/pkg/versioning"
)

func TestStoreEngine_BasicFlow(t *testing.T) {

	filepath := "test.db"
	defer cleanUp(filepath)

	s, err := NewStorageEngine(filepath, []byte("#"))
	require.NoError(t, err)

	tableName := []byte("TestTable")
	err = s.CreateTable(tableName)
	assert.NoError(t, err)

	hashKey := []byte("Tx_100")

	sortKey := []byte("2022-02-24")

	vectorC, err := versioning.NewEmptyClock().Incremented(1, time.Now().Unix())
	require.NoError(t, err)
	err = s.Put(tableName, hashKey, sortKey, versioning.Versioned{
		Version:  vectorC,
		Contents: []byte("Test transaction"),
	})
	require.NoError(t, err)

	result, err := s.Get(tableName, hashKey, sortKey)
	assert.NoError(t, err)

	assert.Equal(t, []byte("Test transaction"), result.Contents)
}

func cleanUp(filename string) {
	os.Remove(filename)
}
