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
	var out bytes.Buffer
	m := NewManager(1024, &out)

	// Start a write that never completes
	r, ok := m.BeginWrite(100)
	require.True(t, ok)

	// Don't call EndWrite() - simulate stuck producer

	// Rotate arena
	sealed := m.rotate()
	require.Equal(t, m.a0, sealed)

	// Start consumer with short-lived context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		m.ConsumerLoop(ctx, func(a *Arena, used int64) {
			m.flushArena(a)
		})
		close(done)
	}()

	// Should exit within timeout, not hang forever
	select {
	case <-done:
		// Success
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Consumer didn't exit after context cancel")
	}

	// Clean up stuck producer
	m.EndWrite(r)
}
