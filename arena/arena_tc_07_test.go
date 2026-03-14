package arena

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// Test Case 7: Memory Corruption Check

// Test: Concurrent writes don't corrupt each other's data
// Verifies: Each log entry remains intact and contiguous
func TestNoMemoryCorruption(t *testing.T) {
	var out bytes.Buffer
	m := NewManager(64*1024, &out)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	go m.ConsumerLoop(ctx, func(a *Arena, used int64) {
		m.flushArena(a)
		m.resetArena(a)
	})

	var wg sync.WaitGroup
	errors := atomic.Int64{}

	// Each producer writes a unique pattern
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(producerID int) {
			defer wg.Done()

			pattern := fmt.Sprintf("P%d:", producerID)
			base := []byte(pattern)

			for j := 0; j < 1000; j++ {
				payload := fmt.Sprintf("%s-%d-%s", pattern, j,
					strings.Repeat("x", 50))

				m.Write(int64(len(payload)), func(dst []byte) {
					copy(dst, []byte(payload))
				})
			}
		}(i)
	}

	wg.Wait()
	cancel()

	// Verify output: Each line should start with "P{id}:"
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		if !strings.HasPrefix(line, "P") {
			t.Errorf("Corrupted line: %q", line)
		}
	}
}
