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
	rawLogger := NewRawLogger(_Size1K, &bytes.Buffer{})

	// Use defer/recover to simulate panic in producer
	func() {
		defer func() { recover() }()

		region, canWrite := rawLogger.beginWrite(100)
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

	require.Zero(t,
		activeArena.numberWriters.Load(),
		"blocked writer should be handled",
	)
}

func TestProducerInWritePanic(t *testing.T) {
	rawLogger := NewRawLogger(_Size1K, &bytes.Buffer{})

	// Use defer/recover to simulate panic in producer
	func() {
		defer func() { recover() }()

		payload := t.Name()

		rawLogger.write(
			uint32(len(payload)),

			func(dst []byte) {
				// Panic before EndWrite
				panic("simulated crash")
			},
		)
	}()

	// writers should be zero as the write finishes with an endwrite.
	activeArena := rawLogger.active.Load()
	require.Zero(t,
		activeArena.numberWriters.Load(),
		"blocked writer should be handled",
	)
}
