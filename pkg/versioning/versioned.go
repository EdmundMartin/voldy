package versioning

import "bytes"

type Versioned struct {
	Version  *VectorClock
	Contents []byte
}

func NewVersionedBytes(contents []byte, version *VectorClock) *Versioned {
	if version == nil {
		version = NewEmptyClock()
	}
	return &Versioned{
		Version:  version,
		Contents: contents,
	}
}

func (v *Versioned) VersionedToBytes() []byte {
	buf := bytes.Buffer{}
	sizeContents := len(v.Contents)
	buf.Write(uint16ToBytes(uint16(sizeContents)))
	buf.Write(v.Contents)
	buf.Write(v.Version.ToBytes())
	return buf.Bytes()
}

func VersionedFromBytes(byteArr []byte) *Versioned {
	v := &Versioned{}
	sizeContents := readUint16(byteArr)
	endOffset := 2 + sizeContents
	v.Contents = byteArr[2:endOffset]
	v.Version = VectorClockFromBytes(byteArr[endOffset:])
	return v
}

func (v *Versioned) HappenedBefore(other *Versioned) (int, error) {
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

type VersionedCollection []*Versioned

func (v VersionedCollection) Len() int {
	return len(v)
}

func (v VersionedCollection) Less(i, j int) bool {
	res, _ := v[i].HappenedBefore(v[j])
	return res < 0
}

func (v VersionedCollection) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}
