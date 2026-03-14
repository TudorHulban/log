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
	manager := NewManager(100, &bytes.Buffer{})
	arena := manager.active.Load()

	// Fill arena to 90 bytes
	arena.cursor.Store(90)

	// Producer 1: Reserve 10 bytes (should fit exactly)
	reserve10, couldReserve10 := manager.BeginWrite(10)
	require.True(t, couldReserve10)
	require.Equal(t,
		int64(90),
		reserve10.offset,
	)

	// Producer 2: Reserve 1 byte (should fail - overflow)
	reserveZero, couldNotReserveMore := manager.BeginWrite(1)
	require.False(t, couldNotReserveMore)
	require.Zero(t, reserveZero)
	require.Equal(t,
		int64(0),
		arena.rollback.Load(),
	)

	// Complete first write
	manager.EndWrite(reserve10)

	// Verify: Final cursor at 100
	require.Equal(t, int64(100), arena.cursor.Load())
}
