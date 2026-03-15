package arena

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Test Case 5: Context Cancellation During Wait

// Test: Consumer context cancelled while waiting for writers
// Verifies: Shutdown happens promptly, no hangs
func TestContextCancelDuringWait(t *testing.T) {
	rawLogger := NewRawLogger(1024, &bytes.Buffer{})

	// Start a write that never completes
	region, couldWrite := rawLogger.BeginWrite(100)
	require.True(t, couldWrite)

	// Don't call EndWrite() - simulate stuck producer

	// Rotate arena
	sealed := rawLogger.rotate()
	require.Equal(t, rawLogger.arenaFirst, sealed)

	// Start consumer with short-lived context
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	chConsumerExit := make(chan struct{})

	go func() {
		rawLogger.ConsumerLoop(
			ctx,

			func(a *Arena, used int64) {
				rawLogger.flushArena(a)
			},
		)

		close(chConsumerExit)
	}()

	// Should exit within timeout, not hang forever
	select {
	case <-chConsumerExit:
		// Success

	case <-time.After(100 * time.Millisecond):
		t.Fatal("Consumer did not exit after context cancel")
	}

	// Clean up stuck producer
	rawLogger.EndWrite(region)
}
