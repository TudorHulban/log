package arena

// WriteRegion describes a reserved region inside an arena.
type WriteRegion struct {
	arena *Arena

	offset uint32
	size   uint32
}

// BeginWrite attempts to reserve n bytes in the current active arena.
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
func (m *RawLogger) BeginWrite(n uint32) (WriteRegion, bool) {
	arena := m.active.Load()
	if arena == nil {
		return WriteRegion{}, false
	}

	// Enter BEFORE reserving, but we must validate we are still on the
	// active arena. A rotation could have happened between Load and Enter (TOCTOU).
	arena.Enter()

	// Re-check: if the active arena changed after we entered, this arena
	// is now sealed. Leave immediately — the cursor may be reset under us.
	if m.active.Load() != arena {
		arena.Leave()

		return WriteRegion{}, false
	}

	// Reserve space.
	offset := arena.Reserve(n)

	// Check for overflow.
	if offset < 0 || offset+n > m.arenaSize {
		// undo reservation
		arena.cursor.Add(int32(-n))

		arena.AddRollback()
		arena.Leave()

		return WriteRegion{}, false
	}

	return WriteRegion{
			arena:  arena,
			offset: offset,
			size:   n,
		},
		true
}

// Buf returns the writable slice for the reserved region.
func (r WriteRegion) Buf() []byte {
	return r.arena.buf[r.offset : r.offset+r.size]
}

// EndWrite decrements writers-in-flight.
func (m *RawLogger) EndWrite(r WriteRegion) {
	r.arena.Leave()
}
