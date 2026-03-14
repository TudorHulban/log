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

	manager := NewManager(1024, &out)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Start consumer with aggressive rotation
	go manager.ConsumerLoop(
		ctx,

		func(a *Arena, used int64) {
			manager.waitForWriters(a)
			manager.flushArena(a)
			manager.resetArena(a)
		},
	)

	var wg sync.WaitGroup

	writes := 10000
	successCount := atomic.Int64{}

	noProducers := 10

	for ix := range noProducers {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			for j := 0; j < writes/noProducers; j++ {
				payload := fmt.Sprintf(
					"producer-%d-%d\n",
					id,
					j,
				)

				couldWrite := manager.Write(
					int64(len(payload)),
					func(dst []byte) {
						copy(dst, []byte(payload))
					},
				)

				if couldWrite {
					successCount.Add(1)
				}

				// Small random delay to increase race probability
				time.Sleep(time.Duration(rand.Intn(10)) * time.Microsecond)
			}
		}(ix)
	}

	wg.Wait()
	cancel()

	// Verify: All successful writes appear in output
	output := out.String()
	require.NotEmpty(t, output)

	outputNoLines := strings.Split(output, "\n")
	require.EqualValues(t,
		len(outputNoLines),
		int(successCount.Load()),

		"number output lines: %d vs success count of %d",
		len(outputNoLines),
		int(successCount.Load()),
	)

	require.NotZero(t, successCount.Load())

	// Account for possible partial writes at end
	require.Equal(t,
		int(successCount.Load()),
		len(outputNoLines)-1,

		"success count is: %d",
		successCount.Load(),
	)
}
