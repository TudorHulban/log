package arena

import "context"

// Flush sealed arena contents using the provided writer function.
//
// The writer receives:
//   - the arena pointer
//   - the slice of bytes to flush
//
// This function does NOT:
//   - rotate arenas
//   - wait for writers
//   - reset the arena
//   - handle errors
//
// Those responsibilities belong to the consumer loop.
func (m *RawLogger) flushArena(a *arena) {
	if a == nil {
		return
	}

	used := uint32(a.cursor.Load()) //nolint:gosec
	if used == 0 {
		return
	}

	// Cursor can exceed arenaSize because Reserve() does an unconditional
	// fetch-add before overflow is detected. Clamp to actual buffer length.
	if used > m.arenaSize {
		used = m.arenaSize
	}

	_, _ = m.writer.Write(a.buf[:used])
}

// flushOnShutdown flushes both arenas best-effort.
func (m *RawLogger) flushOnShutdown(ctx context.Context, flusher flusher) {
	// First rotation: seal whatever is currently active (call it A).
	firstSealed := m.rotate()

	// Second rotation: seal the other arena (B) which just became active.
	// Any producer that got bumped from A by the first rotate and retried
	// into B will be captured here.
	secondSealed := m.rotate()

	// Flush second-sealed first (it became active most recently,
	// producers who retried land here — wait for them first).
	if secondSealed != nil {
		m.waitForWritersCtx(ctx, secondSealed)

		used := secondSealed.cursor.Load()

		if used > 0 {
			flusher(secondSealed, used)
		}
	}

	// Flush first-sealed.
	if firstSealed != nil && firstSealed != secondSealed {
		m.waitForWritersCtx(ctx, firstSealed)

		used := firstSealed.cursor.Load()

		if used > 0 {
			flusher(firstSealed, used)
		}
	}
}

func (m *RawLogger) signalFlush() {
	select {
	case m.chFlush <- struct{}{}:

	default: // signal already pending, consumer will handle it
	}
}
