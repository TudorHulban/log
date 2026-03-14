package arena

// Write attempts to write n bytes into the active arena.
// The caller provides a function that writes into the reserved buffer.
//
// The write function receives a byte slice of length n and must fill it.
func (m *Manager) Write(n int64, fn func(dst []byte)) bool {
	// Try to region space (with one retry).
	region, canWrite := m.TryWrite(n)
	if !canWrite {
		return false
	}

	// Write into the reserved region.
	fn(region.Buf())

	// Mark write complete.
	m.EndWrite(region)

	return true
}
