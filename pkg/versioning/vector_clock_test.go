package versioning

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func getClock(nodes ...int) *VectorClock {
	vectorClock := NewEmptyClock()
	increment(vectorClock, nodes...)
	return vectorClock
}

func increment(clock *VectorClock, nodes ...int) {
	for _, n := range nodes {
		err := clock.IncrementVersion(n, clock.timestamp)
		if err != nil {
			return
		}
	}
	return
}

func TestVectorClock_Compare(t *testing.T) {
	result, err := getClock().Compare(getClock())
	assert.NoError(t, err)
	assert.NotEqual(t, CONCURRENTLY, result)

	result, err = getClock(1, 1, 2).Compare(getClock(1, 1, 2))
	assert.NoError(t, err)
	assert.NotEqual(t, CONCURRENTLY, result)

	result, err = getClock(1, 1, 2).Compare(getClock(1, 1, 2, 3))
	assert.NoError(t, err)
	assert.Equal(t, BEFORE, result)

	result, err = getClock(1).Compare(getClock(2))
	assert.NoError(t, err)
	assert.Equal(t, CONCURRENTLY, result)

	result, err = getClock(1, 1, 2).Compare(getClock(1, 1, 3))
	assert.NoError(t, err)
	assert.Equal(t, CONCURRENTLY, result)

	result, err = getClock(1, 2, 3, 3).Compare(getClock(1, 1, 2, 3))
	assert.NoError(t, err)
	assert.Equal(t, CONCURRENTLY, result)

	/*
		TODO - Investigate case logic looks correct
			result, err := getClock(2, 2).Compare(getClock(1, 1, 2, 3))
			assert.NoError(t, err)
			assert.Equal(t, BEFORE, result)

	*/

	result, err = getClock(1, 2, 2, 3).Compare(getClock(2, 2))
	assert.NoError(t, err)
	assert.Equal(t, AFTER, result)
}

func TestVectorClock_ToBytes(t *testing.T) {

	clock := getClock(1, 2, 3, 4, 4, 1, 1, 1)
	recoveredClock := VectorClockFromBytes(clock.ToBytes())
	assert.Equal(t, clock.timestamp, recoveredClock.timestamp)
	assert.Equal(t, clock.versionMap, recoveredClock.versionMap)
}
