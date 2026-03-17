package arena

// tick performs one consumer iteration:
// - checks if active arena should be sealed
// - rotates if needed
// - drains writers
// - flushes sealed arena
func (m *RawLogger) tick(flusher flusher) {
	activeArena := m.active.Load()
	if activeArena == nil {
		return
	}

	if !m.shouldSeal(activeArena) {
		return
	}

	sealedArena := m.rotate()
	if sealedArena == nil {
		return
	}

	m.waitForWriters(sealedArena)

	used := sealedArena.cursor.Load()
	if used > 0 {
		flusher(sealedArena, used)
	}

	m.resetArena(sealedArena)
	m.sealed.Store(nil)
}
