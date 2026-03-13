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

	// Flush only the written portion.
	_, _ = m.writer.Write(a.buf[:used])
}
