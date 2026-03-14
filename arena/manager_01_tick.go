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

	// Do NOT call resetArena here.
	// The flush callback owns the full lifecycle: waitForWriters + flush + reset.
	m.sealed.Store(nil)
}
