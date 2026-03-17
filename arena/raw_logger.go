package arena

import (
	"context"
	"io"
	"sync/atomic"
	"time"
)

// RawLogger owns the two arenas and coordinates which one is active.
// It also handles the rotation and flush on context cancellation.
type RawLogger struct {
	writer io.Writer

	chFlush chan struct{}

	// Pointer to the currently active arena.
	// Producers read this atomically to know where to write.
	active atomic.Pointer[arena]

	// The two arenas used in double-buffer rotation.
	arenaFirst  *arena
	arenaSecond *arena

	// The arena currently sealed and waiting to be flushed.
	// This is informational; consumer logic will use it.
	sealed atomic.Pointer[arena]

	// Size of each arena (capacity of Arena.Buf).
	arenaSize uint32
}

// NewRawLogger allocates two arenas of the given size and initializes
// the Manager with a0 as the active arena and a1 as the standby arena.
func NewRawLogger(arenaSize uint32, w io.Writer) *RawLogger {
	// Allocate arena buffers.
	a0 := &arena{
		buf: make([]byte, arenaSize),
	}

	result := RawLogger{
		arenaFirst: a0,
		arenaSecond: &arena{
			buf: make([]byte, arenaSize),
		},

		arenaSize: arenaSize,
		writer:    w,

		chFlush: make(chan struct{}, 1),
	}

	// Set active arena to a0.
	result.active.Store(a0)

	// No sealed arena yet.
	result.sealed.Store(nil)

	return &result
}

// StartIngestion launches the consumer loop in a goroutine.
// The caller provides the flush function, which receives the
// raw bytes of each sealed arena.
func (m *RawLogger) StartIngestion(ctx context.Context) <-chan struct{} {
	chIngestionEnd := make(chan struct{})

	go func() {
		defer close(chIngestionEnd)

		m.consumerLoop(
			ctx,

			func(a *arena, used int32) {
				m.flushArena(a)
			},
		)
	}()

	return chIngestionEnd
}

// consumerLoop is the main consumer goroutine.
// It monitors the active arena, seals it when needed, waits for writers,
// flushes it, and resets it.
//
// This is only the skeleton — flushing and thresholds are implemented elsewhere.
func (m *RawLogger) consumerLoop(ctx context.Context, flusher flusher) {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Shutdown: flush both arenas best-effort.
			m.flushOnShutdown(ctx, flusher)

			return

		case <-ticker.C:
			m.tick(flusher)

			// consumerLoop gets a third case:
		case <-m.chFlush:
			m.tick(flusher) // same seal/wait/flush/reset as ticker path
		}

	}
}

// resetArena clears the arena state so it can be reused after flushing.
// This does NOT reallocate the buffer.
func (*RawLogger) resetArena(a *arena) {
	a.cursor.Store(0)

	// numberWriters is intentionally NOT reset here.
	// waitForWriters guarantees it reaches zero before this arena
	// is reused. Resetting it here would race with in-flight writers
	// still holding Enter(), corrupting the count to -1 and hanging
	// the next waitForWriters call permanently.
	// a.numberWriters.Store(0)

	a.rollbackCounter.Store(0)
}
