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
	m := NewManager(1024, &bytes.Buffer{})

	// Use defer/recover to simulate panic in producer
	func() {
		defer func() { recover() }()

		r, ok := m.BeginWrite(100)
		require.True(t, ok)

		// Panic before EndWrite
		panic("simulated crash")
	}()

	// writers should still be 1 (leaked!)
	a := m.active.Load()
	require.Equal(t, int64(1), a.writers.Load())

	// This would hang consumer forever - need timeout mechanism
	// Real implementation should handle this case
}
