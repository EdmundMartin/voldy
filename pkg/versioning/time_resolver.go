package versioning

type TimeBasedResolver[T any] struct{}

func (t TimeBasedResolver[T]) ResolveConflicts(items []*Versioned[T]) []*Versioned[T] {
	if len(items) <= 1 {
		return items
	}
	maxItem := items[0]
	maxTime := items[0].Version.timestamp
	maxClock := items[0].Version
	for _, item := range items {
		clock := item.Version
		if clock.timestamp > maxTime {
			maxItem = item
			maxTime = item.Version.timestamp
		}
		maxClock = maxClock.Merge(clock)
	}
	maxClockVersioned := &Versioned[T]{
		Version:  maxClock,
		Contents: maxItem.Contents,
	}
	return []*Versioned[T]{
		maxClockVersioned,
	}
}
