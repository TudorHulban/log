package arena

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test Case 9: Exact Arena Size Edge Cases

// Test: Writes of exact arena size, writes larger than arena
// Verifies: Flood handling as described in Arena.md
func TestExactAndOversizedWrites(t *testing.T) {
	var sizeArena uint32 = 100

	rawLogger := NewRawLogger(sizeArena, &bytes.Buffer{})

	// Case 1: Write exactly arena size
	region, canWrite := rawLogger.BeginWrite(sizeArena)
	require.True(t, canWrite)
	require.Zero(t, region.offset)

	rawLogger.EndWrite(region)

	require.EqualValues(t,
		sizeArena,
		rawLogger.active.Load().cursor.Load(),
	)

	// Reset
	rawLogger.rotate()

	// Case 2: Write larger than arena (flooding)
	region, canWrite = rawLogger.BeginWrite(sizeArena + 1)
	require.False(t, canWrite)
	require.Zero(t, region)

	// Rollback should increment
	arenaActive := rawLogger.active.Load()
	require.EqualValues(t,
		1,
		arenaActive.rollbackCounter.Load(),
	)

	// But cursor should NOT move
	require.Zero(t, arenaActive.cursor.Load())
}
