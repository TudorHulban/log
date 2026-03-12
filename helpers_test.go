package log

import (
	"sync/atomic"
)

// real writer that has observable side effects
// thus the compiler cannot eliminate code.
type countWriter struct {
	n atomic.Int64
}

func (w *countWriter) Write(p []byte) (int, error) {
	w.n.Add(int64(len(p)))

	return len(p), nil
}
