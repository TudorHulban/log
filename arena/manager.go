package arena

import (
	"context"
	"io"
	"sync/atomic"
	"time"
)

// Manager owns the two arenas and coordinates which one is active.
// It contains only data layout — rotation logic is implemented elsewhere.
type Manager struct {
	writer io.Writer

	// Pointer to the currently active arena.
	// Producers read this atomically to know where to write.
	active atomic.Pointer[Arena]

	// The two arenas used in double-buffer rotation.
	a0 *Arena
	a1 *Arena

	// The arena currently sealed and waiting to be flushed.
	// This is informational; consumer logic will use it.
	sealed atomic.Pointer[Arena]

	// Size of each arena (capacity of Arena.Buf).
	arenaSize int64
}

// NewManager allocates two arenas of the given size and initializes
// the Manager with a0 as the active arena and a1 as the standby arena.
func NewManager(arenaSize int64, w io.Writer) *Manager {
	// Allocate arena buffers.
	a0 := &Arena{
		buf: make([]byte, arenaSize),
	}

	a1 := &Arena{
		buf: make([]byte, arenaSize),
	}

	m := Manager{
		a0:        a0,
		a1:        a1,
		arenaSize: arenaSize,
		writer:    w,
	}

	// Set active arena to a0.
	m.active.Store(a0)

	// No sealed arena yet.
	m.sealed.Store(nil)

	return &m
}

// StartConsumer launches the consumer loop in a goroutine.
// The caller provides the flush function, which receives the
// raw bytes of each sealed arena.
func (m *Manager) StartConsumer(ctx context.Context) <-chan struct{} {
	chSignalStarted := make(chan struct{})

	go func() {
		defer close(chSignalStarted)

		m.ConsumerLoop(
			ctx,
			func(a *Arena, used int64) {
				m.flushArena(a)
			},
		)
	}()

	return chSignalStarted
}

// Write attempts to write n bytes into the active arena.
// The caller provides a function that writes into the reserved buffer.
//
// The write function receives a byte slice of length n and must fill it.
func (m *Manager) Write(n int64, fn func(dst []byte)) bool {
	// Try to reserve space (with one retry).
	r, ok := m.TryWrite(n)
	if !ok {
		return false
	}

	// Write into the reserved region.
	fn(r.Buf())

	// Mark write complete.
	m.EndWrite(r)

	return true
}

// Active returns the currently active arena.
func (m *Manager) Active() *Arena {
	return m.active.Load()
}

// Sealed returns the currently sealed arena (may be nil).
func (m *Manager) Sealed() *Arena {
	return m.sealed.Load()
}

// ConsumerLoop is the main consumer goroutine.
// It monitors the active arena, seals it when needed, waits for writers,
// flushes it, and resets it.
//
// This is only the skeleton — flushing and thresholds are implemented elsewhere.
func (m *Manager) ConsumerLoop(ctx context.Context, flush func(a *Arena, used int64)) {
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// Shutdown: flush both arenas best-effort.
			m.flushOnShutdown(flush)
			return

		case <-ticker.C:
			m.tick(flush)
		}
	}
}

// shouldSeal determines whether the active arena should be sealed.
//
// This is a simple heuristic combining:
//   - cursor threshold (almost full)
//   - rollback pressure (many failed reservations)
//
// The exact thresholds can be tuned later.
func (m *Manager) shouldSeal(a *Arena) bool {
	used := a.cursor.Load()

	// Hard threshold: near capacity.
	if used >= m.arenaSize {
		return true
	}

	// Soft threshold: "almost full".
	// Example: seal when 90% full.
	if used >= (m.arenaSize*9)/10 {
		return true
	}

	// Rollback pressure: many producers failed to reserve space.
	// This indicates high contention near the end of the arena.
	if a.rollback.Load() > 0 {
		return true
	}

	return false
}

// resetArena clears the arena state so it can be reused after flushing.
// This does NOT reallocate the buffer.
func (m *Manager) resetArena(a *Arena) {
	a.cursor.Store(0)
	a.writers.Store(0)
	a.rollback.Store(0)
}

// waitForWriters blocks until writers-in-flight reaches zero.
// Context cancellation is handled by the caller.
func (*Manager) waitForWriters(a *Arena) {
	for {
		if a.writers.Load() == 0 {
			return
		}

		time.Sleep(1 * time.Microsecond)
	}
}
