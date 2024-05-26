package safewriter

import (
	"io"
)

type SafeWriter struct {
	chWrites chan []byte
	chStop   chan struct{}

	writer io.Writer
}

type SafeWriterInfo struct {
	Writer   *SafeWriter
	ChWrites chan []byte
	ChStop   chan struct{}
}

func NewSafeWriterInfo(writer io.Writer) *SafeWriterInfo {
	w := SafeWriter{
		writer:   writer,
		chWrites: make(chan []byte, 100),
		chStop:   make(chan struct{}),
	}

	return &SafeWriterInfo{
		Writer:   &w,
		ChWrites: w.chWrites,
		ChStop:   w.chStop,
	}
}

func (safe *SafeWriter) Listen() {
	for {
		select {
		case <-safe.chStop:
			return

		case msg := <-safe.chWrites:
			safe.writer.Write(msg)
		}
	}
}

func (safe *SafeWriter) write(payload []byte) {
	safe.chWrites <- payload
}
