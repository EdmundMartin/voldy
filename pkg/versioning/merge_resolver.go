package versioning

type Merger[T any] interface {
	Merge(t1, t2 T) T
}

type MergeResolver[T any] struct {
	Merger Merger[T]
}

func (mr *MergeResolver[T]) ResolveConflicts(items []*Versioned[T]) []*Versioned[T] {

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
	return []*Versioned[T]{
		{
			clock,
			merged,
		},
	}
}
