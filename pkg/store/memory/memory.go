package memory

import (
	"sync"
	"time"
	"voldy/pkg/store"
	"voldy/pkg/versioning"
)

type InMemoryStorageEngine struct {
	Name    string
	NodeId  int
	storage map[string][]*versioning.Versioned[[]byte]
	clock   *versioning.VectorClock
	mu      sync.RWMutex
}

func NewInMemoryStorageEngine(name string, nodeId int) *InMemoryStorageEngine {
	return &InMemoryStorageEngine{
		Name:    name,
		NodeId:  nodeId,
		storage: make(map[string][]*versioning.Versioned[[]byte]),
		mu:      sync.RWMutex{},
		clock:   versioning.NewEmptyClock(),
	}
}

func (i *InMemoryStorageEngine) Get(key []byte, transform store.TransformFunction) ([]*versioning.Versioned[[]byte], error) {
	if !store.IsValidKey(key) {
		return nil, store.ErrInvalidKey
	}
	results, ok := i.storage[store.BytesToString(key)]
	if !ok {
		return []*versioning.Versioned[[]byte]{}, nil
	}
	return results, nil
}

func (i *InMemoryStorageEngine) GetAll(keys [][]byte, transform store.TransformFunction) (map[string][]*versioning.Versioned[[]byte], error) {
	result := map[string][]*versioning.Versioned[[]byte]{}
	for _, k := range keys {
		res, err := i.Get(k, transform)
		if err != nil {
			return nil, err
		}
		keyStr := store.BytesToString(k)
		result[keyStr] = res
	}
	return result, nil
}

func (i *InMemoryStorageEngine) Put(key []byte, versioned *versioning.Versioned[[]byte], transform store.TransformFunction) error {
	if !store.IsValidKey(key) {
		return store.ErrInvalidKey
	}
	// TODO - Client should present us a version based on the read. We then increment that clock as opposed to our
	// own clock.
	i.clock = i.clock.Clone()
	i.clock.IncrementVersion(i.NodeId, time.Now().UnixMilli())
	versioned.Version = i.clock

	keyStr := store.BytesToString(key)
	result, ok := i.storage[keyStr]
	if !ok {
		i.storage[store.BytesToString(key)] = []*versioning.Versioned[[]byte]{
			versioned,
		}
		return nil
	}
	var itemsToKeep []*versioning.Versioned[[]byte]
	for _, item := range result {
		occurred, err := versioned.Version.Compare(item.Version)
		if err != nil {
			return err
		}
		if occurred == versioning.BEFORE {
			return store.ErrObsoleteVersion
		}
		if occurred == versioning.CONCURRENTLY {
			itemsToKeep = append(itemsToKeep, item)
		}
	}
	itemsToKeep = append(itemsToKeep, versioned)
	i.storage[keyStr] = itemsToKeep
	return nil
}

func (i *InMemoryStorageEngine) Delete(key []byte, version *versioning.VectorClock) (bool, error) {
	if !store.IsValidKey(key) {
		return false, store.ErrInvalidKey
	}
	i.mu.Lock()
	defer i.mu.Unlock()
	keyStr := store.BytesToString(key)
	values, ok := i.storage[keyStr]
	if !ok {
		return false, nil
	}

	if version == nil {
		delete(i.storage, keyStr)
		return true, nil
	}

	itemDeleted := false

	var retained []*versioning.Versioned[[]byte]
	for _, item := range values {
		occurred, err := item.Version.Compare(version)
		if err != nil {
			return false, err
		}
		if occurred == versioning.BEFORE {
			itemDeleted = true
		} else {
			retained = append(retained, item)
		}
	}

	if len(retained) == 0 {
		delete(i.storage, keyStr)
	} else {
		i.storage[keyStr] = retained
	}
	return itemDeleted, nil
}

func (i *InMemoryStorageEngine) GetName() string {
	return i.Name
}

func (i *InMemoryStorageEngine) Close() error {
	return nil
}

func (i *InMemoryStorageEngine) GetCapability() {
	return
}

func (i *InMemoryStorageEngine) GetVersions(key []byte) ([]versioning.Version, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	res, err := i.Get(key, nil)
	if err != nil {
		return []versioning.Version{}, err
	}
	result := make([]versioning.Version, len(res))
	for _, item := range res {
		result = append(result, item.Version)
	}
	return result, nil
}

func (i *InMemoryStorageEngine) Entries() []store.Pair {
	i.mu.RLock()
	defer i.mu.RUnlock()
	result := make([]store.Pair, len(i.storage))
	idx := 0
	for key, value := range i.storage {
		result[idx] = store.Pair{
			Key:   store.StringToBytes(key),
			Value: value,
		}
		idx++
	}
	return result
}

func (i *InMemoryStorageEngine) Keys() [][]byte {
	//TODO implement me
	panic("implement me")
}

func (i *InMemoryStorageEngine) KeysForPartition(partition int) ([][]byte, error) {
	return nil, store.ErrUnsupportedForStorageType
}

func (i *InMemoryStorageEngine) Truncate() {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.storage = map[string][]*versioning.Versioned[[]byte]{}
}

func (i *InMemoryStorageEngine) IsPartitionAware() bool {
	return false
}

func (i *InMemoryStorageEngine) IsPartitionScanSupported() bool {
	return false
}

func (i *InMemoryStorageEngine) GetAndLock(key []byte) {
	// Unsupported for this storage engine
	return
}

func (i *InMemoryStorageEngine) PutAndUnlock(key []byte) {
	// Unsupported for this storage engine
	return
}

func (i *InMemoryStorageEngine) ReleaseLock(key []byte) {
	// Unsupported for this storage engine
	return
}

func (i *InMemoryStorageEngine) EndBatchModifications() bool {
	return false
}
