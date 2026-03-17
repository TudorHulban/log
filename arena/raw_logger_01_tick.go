package arena

type flusher func(a *arena, used int32)

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

	// Do NOT call resetArena here.
	// The flush callback owns the full lifecycle:
	// waitForWriters + flush + reset.
	m.sealed.Store(nil)
}
