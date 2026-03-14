package arena

import (
	"bytes"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Test Case 3: Seal During Active Writes

// Test: Consumer seals arena while producers are in middle of writing
// Verifies: In-flight writes complete successfully, no writes to sealed arena
func TestSealDuringActiveWrites(t *testing.T) {
	var out bytes.Buffer

	manager := NewManager(1024, &out)

	var wg sync.WaitGroup

	chWritesStarted := make(chan struct{})

	noProducers := 5 // Start producers that write slowly

	for range noProducers {
		wg.Add(1)

		go func() {
			defer wg.Done()

			// Slow write that takes time
			reserve, couldWrite := manager.BeginWrite(100)
			if !couldWrite {
				return
			}

			chWritesStarted <- struct{}{}

			// Simulate slow write (50ms)
			time.Sleep(50 * time.Millisecond)

			// Write data
			copy(reserve.Buf(), bytes.Repeat([]byte("x"), 100))

			manager.EndWrite(reserve)
		}()
	}

	// Wait for all producers to start writing
	for range noProducers {
		<-chWritesStarted
	}

	// Seal arena while writes are in progress
	sealedArena := manager.rotate()
	require.NotNil(t, sealedArena)

	// Try to write to active arena (should be new one)
	reserve, couldWrite := manager.BeginWrite(10)
	require.True(t, couldWrite)
	require.Equal(t, manager.a1, reserve.a) // Should be other arena
	manager.EndWrite(reserve)

	// Wait for all slow writes to complete
	wg.Wait()

	// Verify: Sealed arena has writers=0
	require.Equal(t,
		int64(0),
		sealedArena.writers.Load(),
	)

	// Now safe to flush
	used := sealedArena.cursor.Load()
	require.Equal(t,
		used,
		out.Len(),
	)

	manager.flushArena(sealedArena)
	manager.resetArena(sealedArena)

	// Verify: All 5 writes were flushed
	require.Equal(t, 5*100, len(out.String()))
}
