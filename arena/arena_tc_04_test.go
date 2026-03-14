package arena

import (
	"bytes"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test Case 4: Overflow and Rollback Storm

// Test: Many producers simultaneously attempt writes near arena end
// Verifies: Rollback counter correctly tracks failures, no deadlocks
func TestRollbackStorm(t *testing.T) {
	m := NewManager(1000, &bytes.Buffer{})
	a := m.active.Load()

	// Fill arena near capacity
	a.cursor.Store(950)

	var wg sync.WaitGroup
	rollbacks := atomic.Int64{}
	successes := atomic.Int64{}

	// 100 concurrent producers each trying to write varying sizes
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				// Random size between 10-100 bytes
				size := int64(10 + rand.Intn(90))

				r, ok := m.BeginWrite(size)
				if ok {
					successes.Add(1)
					m.EndWrite(r)
				} else {
					rollbacks.Add(1)
				}
			}
		}()
	}

	wg.Wait()

	// Verify: Rollback counter matches failures
	require.Equal(t, rollbacks.Load(), a.rollback.Load())
	require.True(t, successes.Load() > 0 || rollbacks.Load() > 0)
	require.True(t, a.cursor.Load() <= 1000) // Never exceed
}
