package arena

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test Case 2: Reservation Exactly at Arena Boundary

// Test: Producer reserves last bytes exactly at arenaSize-1
// Verifies: Bounds checking works, no off-by-one errors
func TestReservationAtBoundary(t *testing.T) {
	m := NewManager(100, &bytes.Buffer{})
	a := m.active.Load()

	// Fill arena to 90 bytes
	a.cursor.Store(90)

	// Producer 1: Reserve 10 bytes (should fit exactly)
	r1, ok1 := m.BeginWrite(10)
	require.True(t, ok1)
	require.Equal(t, int64(90), r1.offset)

	// Producer 2: Reserve 1 byte (should fail - overflow)
	r2, ok2 := m.BeginWrite(1)
	require.False(t, ok2)
	require.Equal(t, int64(0), a.rollback.Load())

	// Complete first write
	m.EndWrite(r1)

	// Verify: Final cursor at 100
	require.Equal(t, int64(100), a.cursor.Load())
}
