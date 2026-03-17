package arena

// shouldSeal determines whether the active arena should be sealed.
//
// This is a simple heuristic combining:
//   - cursor threshold (almost full)
//   - rollback pressure (many failed reservations)
//
// The exact thresholds can be tuned later.
func (m *RawLogger) shouldSeal(a *arena) bool {
	used := uint32(a.cursor.Load()) //nolint:gosec

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
	if a.rollbackCounter.Load() > 0 {
		return true
	}

	return false
}
