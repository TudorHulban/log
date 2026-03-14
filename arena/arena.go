package arena

// arena.go

import (
	"sync/atomic"
)

// TODO: move to uint32 aatomics for less cache pressure, smaller.

// arena represents a single fixed-size logging buffer used in a
// double-buffered, lock-free producer/consumer setup.
//
// Methods are defined elsewhere; this file only defines the data layout.
type Arena struct { //nolint:govet
	// Hot atomics (each on its own cache line).

	// cursor is the current write position (in bytes) inside buf.
	// Producers use atomic fetch-add on this to reserve regions.
	cursor atomic.Int64
	_      [56]byte // pad to 64 bytes (typical cache line size)

	// writers tracks the number of producers currently writing into this arena.
	// The consumer waits for this to reach zero before flushing.
	writers atomic.Int64
	_       [56]byte // pad to 64 bytes

	// rollback counts failed reservations near the end of the arena.
	// Used by the consumer as a signal that the arena is under pressure.
	rollback atomic.Int64
	_        [56]byte // pad to 64 bytes

	// buf is the underlying byte storage for this arena.
	// Its capacity defines the arena size.
	buf []byte
}

// Reserve attempts to reserve N bytes inside the arena.
//
// It returns:
//
//	offset >= 0  → success, producer may write into buf[offset : offset+N]
//	offset < 0   → reservation failed (overflow)
//
// The caller must check (offset + N <= arenaSize).
// If not, the caller must treat this as a failed reservation and NOT write.
//
// This function does NOT roll back the cursor; rollback is logical only.
// The consumer will reset the arena when rotating.
func (a *Arena) Reserve(n int64) int64 {
	// Atomically reserve space by bumping the cursor.
	// offset = old cursor value
	offset := a.cursor.Add(n) - n

	return offset
}

// Enter increments the writers-in-flight counter.
// Producers must call this before attempting a reservation.
func (a *Arena) Enter() {
	a.writers.Add(1)
}

// Leave decrements the writers-in-flight counter.
// Producers must call this after finishing their write.
func (a *Arena) Leave() {
	a.writers.Add(-1)
}

// AddRollback increments the rollback counter.
// Producers call this when a reservation overflows the arena.
func (a *Arena) AddRollback() {
	a.rollback.Add(1)
}
