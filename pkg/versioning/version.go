package versioning

type Version interface {
	Compare(v *VectorClock) (Occurred, error)
}
