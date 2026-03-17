package arena

// tryWrite attempts BeginWrite once. If it fails, it reloads the active
// arena and tries exactly one more time.
//
// This is a convenience helper for callers who want a simple
// "try once, rotate may have happened, try again" pattern.
//
// It does NOT loop indefinitely and does NOT block.
func (m *RawLogger) tryWrite(n uint32) (WriteRegion, error) {
	// First attempt.
	region, canWrite := m.beginWrite(n)
	if canWrite == nil {
		return region, nil
	}

	// Reload active arena — rotation may have occurred.
	// Second attempt.
	return m.beginWrite(n)
}

// beginWrite attempts to reserve n bytes in the current active arena.
//
// On success:
//   - writers-in-flight is incremented
//   - a region is returned
//   - caller MUST call EndWrite
//
// On failure:
//   - writers-in-flight is decremented
//   - reservation if reversed
//   - rollback counter is incremented
//   - ok == false
func (m *RawLogger) beginWrite(n uint32) (WriteRegion, error) {
	arena := m.active.Load()
	if arena == nil {
		return WriteRegion{},
			ErrWriteNoActiveArena
	}

	// Enter BEFORE reserving, but we must validate we are still on the
	// active arena. A rotation could have happened between Load and Enter (TOCTOU).
	arena.Enter()

	// Re-check: if the active arena changed after we entered, this arena
	// is now sealed. Leave immediately — the cursor may be reset under us.
	if m.active.Load() != arena {
		arena.Leave()

		return WriteRegion{},
			ErrWriteActiveArenaMismatch
	}

	// Reserve space.
	offset := arena.Reserve(n)

	// Check for overflow.
	if offset+n > m.arenaSize {
		// undo reservation
		arena.cursor.Add(int32(-n)) //nolint:gosec

		arena.AddRollback()
		arena.Leave()

		m.signalFlush()

		return WriteRegion{},
			ErrWriteMessageTooLarge
	}

	return WriteRegion{
			arena:  arena,
			offset: offset,
			size:   n,
		},
		nil
}

// write attempts to write n bytes into the active arena.
// The caller provides a function that writes into the reserved buffer.
//
// The write function receives a byte slice of length n and must fill it.
func (m *RawLogger) write(n uint32, fn func(dst []byte)) error {
	// Try to region space (with one retry).
	region, canWrite := m.tryWrite(n)
	if canWrite != nil {
		return canWrite
	}

	// Mark write complete.
	defer m.endWrite(region)

	// Write into the reserved region.
	fn(region.Buf())

	return nil
}

func (m *RawLogger) Write(payload []byte) (int, error) {
	if len(payload) == 0 {
		return 0, nil
	}

	// Fast path: try once
	errWrite := m.write(
		uint32(len(payload)), //nolint:gosec
		func(dst []byte) {
			copy(dst, payload)
		},
	)
	if errWrite != nil {
		return 0, errWrite
	}

	return len(payload), nil
}
