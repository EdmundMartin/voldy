package cluster

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCluster(t *testing.T) {

	nodes := []*Node{
		NewNode(1, "localhost", 8000, 8080, 8081, []int{0}),
		NewNode(2, "localhost", 8001, 8090, 8091, []int{1}),
	}

	res, err := NewCluster("TestCluster", nodes, []*Zone{})
	assert.NotNil(t, res)
	assert.NoError(t, err)
}
