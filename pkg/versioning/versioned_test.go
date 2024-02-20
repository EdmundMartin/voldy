package versioning

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestVersioned_VersionedToBytes(t *testing.T) {
	timestamp := time.Now().UnixMilli()
	v := &Versioned{
		Version: &VectorClock{
			SerialVersionID: 0,
			versionMap:      map[uint16]uint64{10: 1},
			timestamp:       timestamp,
		},
		Contents: []byte("Hello World!"),
	}

	bytes := v.VersionedToBytes()
	otherV := VersionedFromBytes(bytes)

	assert.Equal(t, otherV.Version.timestamp, timestamp)
	oth, err := otherV.Version.versionMap[10]
	assert.True(t, err)
	assert.Equal(t, uint64(1), oth)
	assert.Equal(t, []byte("Hello World!"), otherV.Contents)
}
