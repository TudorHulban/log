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
	m := NewManager(1024, &out)

	// Start 5 producers that write slowly
	var wg sync.WaitGroup
	writesStarted := make(chan struct{})

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			// Slow write that takes time
			r, ok := m.BeginWrite(100)
			if !ok {
				return
			}

			writesStarted <- struct{}{}

			// Simulate slow write (50ms)
			time.Sleep(50 * time.Millisecond)

			// Write data
			copy(r.Buf(), bytes.Repeat([]byte("x"), 100))
			m.EndWrite(r)
		}()
	}

	// Wait for all producers to start writing
	for i := 0; i < 5; i++ {
		<-writesStarted
	}

	// Seal arena while writes are in progress
	sealed := m.rotate()
	require.NotNil(t, sealed)

	// Try to write to active arena (should be new one)
	r, ok := m.BeginWrite(10)
	require.True(t, ok)
	require.Equal(t, m.a1, r.a) // Should be other arena
	m.EndWrite(r)

	// Wait for all slow writes to complete
	wg.Wait()

	// Verify: Sealed arena has writers=0
	require.Equal(t, int64(0), sealed.writers.Load())

	// Now safe to flush
	used := sealed.cursor.Load()
	m.flushArena(sealed)
	m.resetArena(sealed)

	// Verify: All 5 writes were flushed
	require.Equal(t, 5*100, len(out.String()))
}
