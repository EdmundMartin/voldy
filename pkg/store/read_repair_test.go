package store

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"voldy/pkg/versioning"
)

func TestLatestVersion(t *testing.T) {

	now := time.Now().UnixMilli()

	cl := versioning.NewEmptyClock()
	cl.IncrementVersion(1, now)

	clTwo := versioning.NewEmptyClock()
	clTwo.IncrementVersion(1, now)
	clTwo.IncrementVersion(2, now)

	clThree := versioning.NewEmptyClock()
	clThree.IncrementVersion(1, now)
	clThree.IncrementVersion(2, now)
	clThree.IncrementVersion(2, now)

	res := LatestVersion([]*versioning.Versioned{
		{
			Version:  cl,
			Contents: []byte("First Clock"),
		},
		{
			Version:  clTwo,
			Contents: []byte("Second clock"),
		},
		{
			Version:  clThree,
			Contents: []byte("Third version"),
		},
	})
	assert.Equal(t, []byte("Third version"), res.Contents)
}
