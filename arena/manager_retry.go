package arena

// TryWrite attempts BeginWrite once. If it fails, it reloads the active
// arena and tries exactly one more time.
//
// This is a convenience helper for callers who want a simple
// "try once, rotate may have happened, try again" pattern.
//
// It does NOT loop indefinitely and does NOT block.
func (m *Manager) TryWrite(n int64) (WriteRegion, bool) {
	// First attempt.
	r, ok := m.BeginWrite(n)
	if ok {
		return r, true
	}

	// Reload active arena — rotation may have occurred.
	// Second attempt.
	return m.BeginWrite(n)
}
