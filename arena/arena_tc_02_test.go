package arena

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test Case 2: Reservation Exactly at Arena Boundary

// Test: Producer reserves last bytes exactly at arenaSize-1
// Verifies:
// - bounds checking works, no off-by-one errors
// - rollback counter incremented
func TestReservationAtBoundary(t *testing.T) {
	rawLogger := NewRawLogger(100, &bytes.Buffer{})
	arena := rawLogger.active.Load()

	// Fill arena to 90 bytes
	arena.cursor.Store(90)

	// Producer 1: Reserve 10 bytes (should fit exactly)
	region10, couldReserve10 := rawLogger.BeginWrite(10)
	require.True(t, couldReserve10)
	require.Equal(t,
		int64(90),
		region10.offset,
	)

	// Producer 2: Reserve 1 byte (should fail - overflow)
	regionZero, couldReserveMore := rawLogger.BeginWrite(1)
	require.False(t, couldReserveMore)
	require.Zero(t, regionZero)
	require.Equal(t,
		int64(1),
		arena.rollbackCounter.Load(),
	)

	// Complete first write
	rawLogger.EndWrite(region10)

	// Verify: Final cursor at 100
	require.Equal(t,
		int64(100),
		arena.cursor.Load(),
	)
}
