package safewriter

import (
	"io"
	"sync"
)

type SafeWriter struct {
	writeTo io.Writer
	mu      sync.Mutex
}

var _ io.Writer = &SafeWriter{}

func NewSafeWriter(writer io.Writer) *SafeWriter {
	return &SafeWriter{
		writeTo: writer,
	}
}

func (safe *SafeWriter) Write(payload []byte) (n int, err error) {
	safe.mu.Lock()
	defer safe.mu.Unlock()

	return safe.writeTo.Write(payload)
}
