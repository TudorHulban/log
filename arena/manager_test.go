package arena

import (
	"bytes"
	"context"
	"testing"
	"time"
)

func TestManagerSingleWrite(t *testing.T) {
	var out bytes.Buffer

	// Small arena for easy testing.
	manager := NewManager(1024, &out)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start consumer.
	manager.StartConsumer(ctx)

	// Perform one write.
	ok := manager.Write(5, func(dst []byte) {
		copy(dst, []byte("hello"))
	})
	if !ok {
		t.Fatalf("Write returned false")
	}

	// Give the consumer a moment to flush.
	time.Sleep(20 * time.Millisecond)

	// Verify output.
	if out.String() != "hello" {
		t.Fatalf("unexpected output: %q", out.String())
	}
}
