package bytearena

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test Case 10: Producer Panic Safety

// Test: Producer panics during write
// Verifies: writers counter is decremented even on panic
func TestWritePanicDoesNotLeak(t *testing.T) {
	rawLogger := NewIngestor(_Size1K, &bytes.Buffer{})

	func() {
		defer func() { _ = recover() }()

		rawLogger.write(
			100,
			func(dst []byte) {
				panic("some panic")
			},
		)
	}()

	active := rawLogger.active.Load()
	require.Zero(t, active.numberWriters.Load())
}

func TestProducerInWritePanic(t *testing.T) {
	rawLogger := NewIngestor(_Size1K, &bytes.Buffer{})

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
