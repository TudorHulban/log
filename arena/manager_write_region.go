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

	// Mark producer as active.
	a.Enter()

	// Reserve space.
	offset := a.Reserve(n)

	// Check for overflow.
	if offset < 0 || offset+n > m.arenaSize {
		a.AddRollback()
		a.Leave()
		return WriteRegion{}, false
	}

	// Success.
	return WriteRegion{
		a:      a,
		offset: offset,
		size:   n,
	}, true
}

// Buf returns the writable slice for the reserved region.
func (r WriteRegion) Buf() []byte {
	return r.a.buf[r.offset : r.offset+r.size]
}

// EndWrite decrements writers-in-flight.
func (m *Manager) EndWrite(r WriteRegion) {
	r.a.Leave()
}
