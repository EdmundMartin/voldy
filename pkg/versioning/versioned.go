package versioning

type Versioned[T any] struct {
	Version  *VectorClock
	Contents T
}

func NewVersionedBytes(contents []byte, version *VectorClock) *Versioned[[]byte] {
	if version == nil {
		version = NewEmptyClock()
	}
	return &Versioned[[]byte]{
		Version:  version,
		Contents: contents,
	}
}

func (v *Versioned[T]) HappenedBefore(other *Versioned[T]) (int, error) {
	result, err := v.Version.Compare(other.Version)
	if err != nil {
		return 0, err
	}
	if result == BEFORE {
		return -1, nil
	}
	if result == AFTER {
		return 1, nil
	}
	return 0, nil
}