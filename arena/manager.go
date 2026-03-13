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

// rotate seals the current active arena and switches to the other one.
// It returns the sealed arena (the one that was active before the switch).
//
// This function does NOT wait for writers to drain and does NOT flush.
// Waiting for writers and flushing are handled by the consumer logic.
func (m *Manager) rotate() *Arena {
	// Load current active arena.
	current := m.active.Load()
	if current == nil {
		return nil
	}

	// Determine the next arena.
	var next *Arena
	if current == m.a0 {
		next = m.a1
	} else {
		next = m.a0
	}

	// Mark current as sealed.
	m.sealed.Store(current)

	// Switch active to the next arena.
	m.active.Store(next)

	return current
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

// tick performs one consumer iteration:
// - checks if active arena should be sealed
// - rotates if needed
// - drains writers
// - flushes sealed arena
func (m *Manager) tick(flush func(a *Arena, used int64)) {
	active := m.active.Load()
	if active == nil {
		return
	}

	// Check if we should seal the active arena.
	// Threshold logic is implemented elsewhere.
	if !m.shouldSeal(active) {
		return
	}

	// Rotate: active becomes sealed, other becomes active.
	sealed := m.rotate()
	if sealed == nil {
		return
	}

	// Wait for writers to finish.
	m.waitForWriters(sealed)

	// Flush sealed arena.
	used := sealed.cursor.Load()
	if used > 0 {
		flush(sealed, used)
	}

	// Reset sealed arena for reuse.
	m.resetArena(sealed)
}

// waitForWriters blocks until writers-in-flight reaches zero.
// Context cancellation is handled by the caller.
func (m *Manager) waitForWriters(a *Arena) {
	for {
		if a.writers.Load() == 0 {
			return
		}
		time.Sleep(1 * time.Microsecond)
	}
}

// flushOnShutdown flushes both arenas best-effort.
func (m *Manager) flushOnShutdown(flush func(a *Arena, used int64)) {
	// Flush active arena.
	if a := m.active.Load(); a != nil {
		used := a.cursor.Load()
		if used > 0 {
			flush(a, used)
		}
	}

	// Flush sealed arena.
	if s := m.sealed.Load(); s != nil {
		used := s.cursor.Load()
		if used > 0 {
			flush(s, used)
		}
	}
}
