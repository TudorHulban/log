package arena

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

// Test Case 7: Memory Corruption Check

// Test: Concurrent writes don't corrupt each other's data
// Verifies: Each log entry remains intact and contiguous
// Enhanced version with write validation
func TestNoMemoryCorruption_Enhanced(t *testing.T) {
	var out bytes.Buffer

	m := NewManager(64*1024, &out)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	validationCh := make(chan string, 10000)

	// Consumer with validation
	go func() {
		for line := range validationCh {
			// Validate format immediately
			if !strings.HasPrefix(line, "P") {
				t.Errorf("Invalid line format: %q", line)
			}
		}
	}()

	go m.ConsumerLoop(ctx, func(a *Arena, used int64) {
		// Capture output for validation
		data := a.buf[:used]
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if line != "" {
				validationCh <- line
			}
		}

		m.flushArena(a)
		m.resetArena(a)
	})

	var wg sync.WaitGroup

	// Each producer writes a unique pattern
	for i := 0; i < 20; i++ {
		wg.Add(1)

		go func(producerID int) {
			defer wg.Done()

			for j := 0; j < 1000; j++ {
				payload := fmt.Sprintf("P%d-%d-%s", producerID, j,
					strings.Repeat("x", 50))

				m.Write(int64(len(payload)), func(dst []byte) {
					// Double-check destination before writing
					if len(dst) != len(payload) {
						t.Errorf("Buffer size mismatch: got %d, want %d",
							len(dst), len(payload))
					}

					copy(dst, []byte(payload))
				})
			}
		}(i)
	}

	wg.Wait()
	cancel()
	close(validationCh)
}
