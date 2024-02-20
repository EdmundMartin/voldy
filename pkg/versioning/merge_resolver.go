package versioning

type Merger interface {
	Merge(t1, t2 []byte) []byte
}

type MergeResolver[T any] struct {
	Merger Merger
}

func (mr *MergeResolver[T]) ResolveConflicts(items []*Versioned) []*Versioned {

	if len(items) <= 1 {
		return items
	}

	current := items[0]
	merged := current.Contents
	clock := current.Version
	for i := 1; i < len(items); i++ {
		current = items[i]
		merged = mr.Merger.Merge(merged, current.Contents)
		clock = clock.Merge(current.Version)
	}
	return []*Versioned{
		{
			clock,
			merged,
		},
	}
}
