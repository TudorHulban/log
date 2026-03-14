package arena

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Test Case 1: Concurrent Writes During Rotation

// Test: Multiple producers writing while consumer rotates arenas
// Verifies: No writes are lost, no panics, all logs eventually appear
func TestConcurrentWritesWithRotation(t *testing.T) {
	var out bytes.Buffer
	m := NewManager(1024, &out)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Start consumer with aggressive rotation
	go m.ConsumerLoop(ctx, func(a *Arena, used int64) {
		m.flushArena(a)
		m.resetArena(a)
	})

	var wg sync.WaitGroup
	writes := 10000
	successCount := atomic.Int64{}

	// 10 concurrent producers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < writes/10; j++ {
				payload := fmt.Sprintf("producer-%d-%d", id, j)
				ok := m.Write(int64(len(payload)), func(dst []byte) {
					copy(dst, []byte(payload))
				})
				if ok {
					successCount.Add(1)
				}
				// Small random delay to increase race probability
				time.Sleep(time.Duration(rand.Intn(10)) * time.Microsecond)
			}
		}(i)
	}

	wg.Wait()
	cancel()

	// Verify: All successful writes appear in output
	output := out.String()
	lines := strings.Split(output, "\n")
	// Account for possible partial writes at end
	require.Equal(t, int(successCount.Load()), len(lines)-1)
}
