package log

import (
	"os"
	"testing"
)

func TestPrint(t *testing.T) {
	l := NewLogger(
		LevelDEBUG,
		os.Stdout,
		true,
	)

	l.PrintMessage("0")

	// <-l.w.ChStop
}

func Benchmark_Print_Logger(b *testing.B) {
	logger := NewLogger(
		LevelINFO,
		nil,
		false,
	)

	b.ResetTimer()

	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				logger.PrintMessage("1")
			}
		},
	)
}
