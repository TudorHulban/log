package arena

import (
	"context"
	"io"
	"runtime"
	"sync/atomic"
	"time"
)

// RawLogger owns the two arenas and coordinates which one is active.
// It also handles the rotation and flush on context cancellation.
type RawLogger struct {
	writer io.Writer

	// Pointer to the currently active arena.
	// Producers read this atomically to know where to write.
	active atomic.Pointer[Arena]

	// The two arenas used in double-buffer rotation.
	arenaFirst  *Arena
	arenaSecond *Arena

	// The arena currently sealed and waiting to be flushed.
	// This is informational; consumer logic will use it.
	sealed atomic.Pointer[Arena]

	// Size of each arena (capacity of Arena.Buf).
	arenaSize int64
}

// NewRawLogger allocates two arenas of the given size and initializes
// the Manager with a0 as the active arena and a1 as the standby arena.
func NewRawLogger(arenaSize int64, w io.Writer) *RawLogger {
	// Allocate arena buffers.
	a0 := &Arena{
		buf: make([]byte, arenaSize),
	}

	result := RawLogger{
		arenaFirst: a0,
		arenaSecond: &Arena{
			buf: make([]byte, arenaSize),
		},

		arenaSize: arenaSize,
		writer:    w,
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
	chSignalStarted := make(chan struct{})

	go func() {
		defer close(chSignalStarted)

		m.consumerLoop(
			ctx,

			func(a *Arena, used int64) {
				m.flushArena(a)
			},
		)
	}()

	return chSignalStarted
}

// consumerLoop is the main consumer goroutine.
// It monitors the active arena, seals it when needed, waits for writers,
// flushes it, and resets it.
//
// This is only the skeleton — flushing and thresholds are implemented elsewhere.
func (m *RawLogger) consumerLoop(ctx context.Context, flusher func(a *Arena, used int64)) {
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
		}
	}
}

// resetArena clears the arena state so it can be reused after flushing.
// This does NOT reallocate the buffer.
func (m *RawLogger) resetArena(a *Arena) {
	a.cursor.Store(0)
	a.numberWriters.Store(0)
	a.rollbackCounter.Store(0)
}

// waitForWriters blocks until writers-in-flight reaches zero.
// should be used in tick.
func (m *RawLogger) waitForWriters(a *Arena) {
	writers := &a.numberWriters

	spin := 0

	for writers.Load() != 0 {
		spin++

		if spin < 64 {
			continue
		}

		spin = 0
		runtime.Gosched()
	}
}

func (m *RawLogger) waitForWritersCtx(ctx context.Context, a *Arena) bool {
	spin := 0

	for {
		if a.numberWriters.Load() == 0 {
			return true
		}

		spin++

		if spin < 50 {
			runtime.Gosched()
			continue
		}

		select {
		case <-ctx.Done():
			return false
		default:
		}

		time.Sleep(time.Microsecond)
	}
}
