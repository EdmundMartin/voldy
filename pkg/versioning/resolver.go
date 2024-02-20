package versioning

type Resolver[T any] interface {
	ResolveConflicts(items []*Versioned) []*Versioned
}
