package log

import (
	"bytes"
	"os"
	"testing"
)

func TestInfo(t *testing.T) {
	l := NewLogger(
		LevelDEBUG,
		os.Stdout,
		true,
	)

	l.Info("0")

	l.Infof("%d", 1)

	// <-l.w.ChStop
}

func TestWithCheckInfo(t *testing.T) {
	var output bytes.Buffer

	l := NewLogger(
		LevelDEBUG,
		&output,
		true,
	)

	l.Info("0")

	l.Infof("%d", 1)

	// <-l.w.ChStop
}

func Benchmark_Info_Logger(b *testing.B) {
	logger := NewLogger(
		LevelINFO,
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
