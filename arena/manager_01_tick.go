package arena

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

	if !m.shouldSeal(active) {
		return
	}

	sealed := m.rotate()
	if sealed == nil {
		return
	}

	m.waitForWriters(sealed)

	used := sealed.cursor.Load()
	if used > 0 {
		flush(sealed, used)
	}

	m.resetArena(sealed)
	m.sealed.Store(nil) // clear so flushOnShutdown doesn't re-flush this
}
