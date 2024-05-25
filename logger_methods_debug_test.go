package log

import (
	"os"
	"testing"
)

func TestDebug(t *testing.T) {
	l := NewLogger(
		LevelDEBUG,
		os.Stdout,
		true,
	)

	l.Debug("0")

	l.Debugf("%d", 1)

	// <-l.w.ChStop
}

func Benchmark_Debug(b *testing.B) {
	logger := NewLogger(
		LevelDEBUG,
		nil,
		false,
	)

	b.ResetTimer()

	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				logger.Info("1")
			}
		},
	)
}
