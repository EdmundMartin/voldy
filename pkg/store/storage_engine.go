package store

type StorageEngine interface {
	Store
	Entries() []Pair
	Keys() [][]byte
	KeysForPartition(partition int) ([][]byte, error)
	Truncate()
	IsPartitionAware() bool
	IsPartitionScanSupported() bool
	GetAndLock(key []byte) // TODO - Return type
	PutAndUnlock(key []byte)
	ReleaseLock(key []byte)
	EndBatchModifications() bool
}
