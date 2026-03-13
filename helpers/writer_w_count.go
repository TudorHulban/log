package helpers

import "sync/atomic"

// real writer that has observable side effects
// thus the compiler cannot eliminate code.
type CountWriter struct {
	N atomic.Int64
}

func (w *CountWriter) Write(p []byte) (int, error) {
	w.N.Add(int64(len(p)))

	return len(p), nil
}
