package safewriter

import (
	"io"
)

type SafeWriter struct {
	chWrite chan []byte
	chStop  chan struct{}

	writer io.Writer
}

var _ io.Writer = &SafeWriter{}

func NewSafeWriter(writer io.Writer) *SafeWriter {
	result := SafeWriter{
		writer: writer,

		chWrite: make(chan []byte),
		chStop:  make(chan struct{}),
	}

	go result.listen()

	return &result
}

func (safe *SafeWriter) listen() {
	for {
		select {
		case <-safe.chStop:
			return

		case msg := <-safe.chWrite:
			_, _ = safe.writer.Write(msg)
		}
	}
}

func (safe *SafeWriter) Write(payload []byte) (int, error) {
	safe.chWrite <- payload

	return len(payload),
		nil
}

func (safe *SafeWriter) Stop() {
	safe.chStop <- struct{}{}
}
