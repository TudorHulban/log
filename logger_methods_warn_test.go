package log

import (
	"os"
	"testing"
)

func TestWarn(t *testing.T) {
	l := NewLogger(
		LevelDEBUG,
		os.Stdout,
		true,
	)

	l.Warn("0")

	l.Warnf("%d", 1)

	// <-l.w.ChStop
}

func BenchmarkLogger_Warn(b *testing.B) {
	logger := NewLogger(
		LevelDEBUG,
		nil,
		false,
	)

	b.ResetTimer()

	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				logger.Warn("1")
			}
		},
	)
}
