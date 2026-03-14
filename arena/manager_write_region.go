package arena

// WriteRegion describes a reserved region inside an arena.
type WriteRegion struct {
	a      *Arena
	offset int64
	size   int64
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
//   - rollback counter is incremented
//   - ok == false
func (m *Manager) BeginWrite(n int64) (WriteRegion, bool) {
	a := m.active.Load()
	if a == nil {
		return WriteRegion{}, false
	}

	// Enter BEFORE reserving, but we must validate we're still on the
	// active arena. A rotation could have happened between Load and Enter.
	a.Enter()

	// Re-check: if the active arena changed after we entered, this arena
	// is now sealed. Leave immediately — the cursor may be reset under us.
	if m.active.Load() != a {
		a.Leave()
		return WriteRegion{}, false
	}

	// Reserve space.
	offset := a.Reserve(n)

	// Check for overflow.
	if offset < 0 || offset+n > m.arenaSize {
		// undo reservation
		a.cursor.Add(-n)

		a.AddRollback()
		a.Leave()

		return WriteRegion{}, false
	}

	return WriteRegion{
			a:      a,
			offset: offset,
			size:   n,
		},
		true
}

// Buf returns the writable slice for the reserved region.
func (r WriteRegion) Buf() []byte {
	return r.a.buf[r.offset : r.offset+r.size]
}

// EndWrite decrements writers-in-flight.
func (m *Manager) EndWrite(r WriteRegion) {
	r.a.Leave()
}
