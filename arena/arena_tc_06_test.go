package arena

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test Case 6: Race Between Reserve and Seal

// Test: Producer reserves space exactly as consumer seals
// Verifies: No writes to arena after it's sealed
func TestReserveVsSealRace(t *testing.T) {
	m := NewManager(1024, &bytes.Buffer{})

	// Channel to coordinate race
	ready := make(chan struct{})
	done := make(chan bool)

	// Producer goroutine
	go func() {
		<-ready // Wait for signal

		// Attempt to reserve
		r, ok := m.BeginWrite(100)
		if ok {
			// If we got a region, it must be in active arena
			if r.a != m.active.Load() {
				done <- false
				return
			}
			m.EndWrite(r)
		}
		done <- true
	}()

	// Consumer goroutine
	go func() {
		<-ready // Wait for same signal

		// Rotate arenas
		sealed := m.rotate()
		_ = sealed
	}()

	// Start both simultaneously
	close(ready)

	// Wait for result
	require.True(t, <-done)

	// Verify invariant: No writes to sealed arena
	sealed := m.sealed.Load()
	if sealed != nil {
		writers := sealed.writers.Load()
		require.True(t, writers == 0 || m.active.Load() == sealed)
	}
}
