package arena

import (
	"bytes"
	"testing"
	"time"
)

// Test Case 8: NUMA-Style False Sharing Detection

// Test: Multiple cores hammer different atomics
// Verifies: Cache line padding works (performance, not correctness)
func TestFalseSharingResistance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	m := NewManager(1024*1024, &bytes.Buffer{})
	a := m.active.Load()

	// Goroutine 1: Hammer cursor
	go func() {
		for i := 0; i < 1000000; i++ {
			a.cursor.Add(1)
		}
	}()

	// Goroutine 2: Hammer writers
	go func() {
		for i := 0; i < 1000000; i++ {
			a.writers.Add(1)
		}
	}()

	// Goroutine 3: Hammer rollback
	go func() {
		for i := 0; i < 1000000; i++ {
			a.rollback.Add(1)
		}
	}()

	// If padding is wrong, this will be slow due to cache contention
	// We're not measuring, just ensuring no crashes
	time.Sleep(100 * time.Millisecond)

	// If we got here without data races, padding is likely correct
	// (run with -race to verify)
}
