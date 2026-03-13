package arena

// ProducerCtx holds the bound function pointers for a producer to write
// into the currently active arena.
//
// This avoids interfaces and avoids reloading the active arena pointer
// on every write. The manager updates this context atomically when
// rotating arenas.
//
// All fields are direct function pointers and fully inlinable.
type ProducerCtx struct {
	// Bound arena pointer.
	a *Arena

	// Bound methods (monomorphic, no interface dispatch).
	enter   func()
	leave   func()
	reserve func(n int64) int64
}

// bindProducerCtx initializes a ProducerCtx for a given arena.
// This is called by the manager when switching active arenas.
func bindProducerCtx(a *Arena) ProducerCtx {
	return ProducerCtx{
		a: a,

		// Directly bind the arena methods.
		enter:   a.Enter,
		leave:   a.Leave,
		reserve: a.Reserve,
	}
}
