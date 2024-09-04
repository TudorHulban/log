package log

import (
	"io"
	"testing"

	"log"
)

func Benchmark_StandardLogger(b *testing.B) {
	logger := log.New(
		io.Discard,
		"",
		log.LstdFlags,
	)

	b.ResetTimer()

	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				logger.Print("1")
			}
		},
	)
}
