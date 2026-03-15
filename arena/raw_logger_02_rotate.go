package arena

// rotate seals the current active arena and switches to the other one.
// It returns the sealed arena (the one that was active before the switch).
//
// This function does NOT wait for writers to drain and does NOT flush.
// Waiting for writers and flushing are handled by the consumer logic.
func (m *RawLogger) rotate() *Arena {
	// Load current active arena.
	current := m.active.Load()
	if current == nil {
		return nil
	}

	// Determine the next arena.
	var next *Arena
	if current == m.arenaFirst {
		next = m.arenaSecond
	} else {
		next = m.arenaFirst
	}

	// Mark current as sealed.
	m.sealed.Store(current)

	// Switch active to the next arena.
	m.active.Store(next)

	return current
}
