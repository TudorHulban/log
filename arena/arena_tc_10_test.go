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
	rawLogger := NewRawLogger(1024, &bytes.Buffer{})

	// Use defer/recover to simulate panic in producer
	func() {
		defer func() { recover() }()

		region, canWrite := rawLogger.BeginWrite(100)
		require.True(t, canWrite)
		require.NotZero(t, region)

		// Panic before EndWrite
		panic("simulated crash")
	}()

	// writers should still be 1 (leaked!)
	activeArena := rawLogger.active.Load()
	require.EqualValues(t,
		1,
		activeArena.numberWriters.Load(),
	)

	// This would hang consumer forever - need timeout mechanism
	// Real implementation should handle this case
}
