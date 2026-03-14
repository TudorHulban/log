package arena

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
func (m *Manager) flushArena(a *Arena) {
	if a == nil {
		return
	}

	used := a.cursor.Load()
	if used <= 0 {
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
func (m *Manager) flushOnShutdown(flush func(a *Arena, used int64)) {
	// Seal active arena to stop new writes.
	sealed := m.rotate()

	// Flush the arena that was active (now sealed).
	if sealed != nil {
		m.waitForWriters(sealed)
		used := sealed.cursor.Load()
		if used > 0 {
			flush(sealed, used)
		}
		m.resetArena(sealed)
	}

	// Flush the other arena only if it has unflushed data
	// (i.e. it was sealed by tick but not yet flushed, which shouldn't
	// happen since tick flushes inline — but guard anyway).
	other := m.sealed.Load()
	if other != nil && other != sealed {
		m.waitForWriters(other)
		used := other.cursor.Load()
		if used > 0 {
			flush(other, used)
		}
	}
}
