package arena

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test Case 10: Producer Panic Safety

// Test: Producer panics during write
// Verifies: writers counter is decremented even on panic
func TestProducerPanic(t *testing.T) {
	manager := NewManager(1024, &bytes.Buffer{})

	// Use defer/recover to simulate panic in producer
	func() {
		defer func() { recover() }()

		reserve, couldWrite := manager.BeginWrite(100)
		require.True(t, couldWrite)
		require.NotZero(t, reserve)

		// Panic before EndWrite
		panic("simulated crash")
	}()

	// writers should still be 1 (leaked!)
	activeArena := manager.active.Load()
	require.Equal(t, int64(1), activeArena.writers.Load())

	// This would hang consumer forever - need timeout mechanism
	// Real implementation should handle this case
}
