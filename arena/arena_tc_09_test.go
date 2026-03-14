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
	m := NewManager(100, &bytes.Buffer{})

	// Case 1: Write exactly arena size
	r, ok := m.BeginWrite(100)
	require.True(t, ok)
	require.Equal(t, int64(0), r.offset)
	m.EndWrite(r)
	require.Equal(t, int64(100), m.active.Load().cursor.Load())

	// Reset
	m.rotate()

	// Case 2: Write larger than arena (flooding)
	r, ok = m.BeginWrite(101)
	require.False(t, ok)

	// Rollback should increment
	a := m.active.Load()
	require.Equal(t, int64(1), a.rollback.Load())

	// But cursor should NOT move
	require.Equal(t, int64(0), a.cursor.Load())
}
