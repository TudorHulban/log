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

	rawLogger := NewRawLogger(1024, &out)

	var wgProducers sync.WaitGroup

	chWritesStarted := make(chan struct{})

	noProducers := 5 // Start producers that write slowly

	wgProducers.Add(noProducers)

	for range noProducers {
		go func() {
			defer wgProducers.Done()

			// Slow write that takes time
			region, couldWrite := rawLogger.BeginWrite(100)
			if !couldWrite {
				return
			}

			chWritesStarted <- struct{}{}

			// Simulate slow write (50ms)
			time.Sleep(50 * time.Millisecond)

			// Write data
			copy(region.Buf(), bytes.Repeat([]byte("x"), 100))

			rawLogger.EndWrite(region)
		}()
	}

	// Wait for all producers to start writing
	for range noProducers {
		<-chWritesStarted
	}

	// Seal arena while writes are in progress
	sealedArena := rawLogger.rotate()
	require.NotNil(t, sealedArena)

	// no rollbacks should have occured
	require.Zero(t, sealedArena.rollback.Load())

	// Try to write to active arena (should be new one)
	region, couldWrite := rawLogger.BeginWrite(10)
	require.True(t, couldWrite)

	// Should be other arena
	require.Equal(t,
		rawLogger.arenaSecond,
		region.arena,
	)
	rawLogger.EndWrite(region)

	// Wait for all slow writes to complete
	wgProducers.Wait()

	// Verify: Sealed arena has writers=0
	require.Equal(t,
		int64(0),
		sealedArena.writers.Load(),
	)

	// Now safe to flush
	used := sealedArena.cursor.Load()

	rawLogger.flushArena(sealedArena)
	require.EqualValues(t,
		used,
		out.Len(),

		"used in sealed arena: %d different than what was written to writer: %d",
		used,
		out.Len(),
	)

	rawLogger.resetArena(sealedArena)

	// Verify: All 5 writes were flushed
	require.Equal(t, 5*100, len(out.String()))
}
