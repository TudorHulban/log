package log

import (
	"os"
	"testing"

	"github.com/TudorHulban/log/timestamp"
)

func TestDebug(t *testing.T) {
	l := NewLogger(
		&ParamsNewLogger{
			LoggerLevel:  LevelDEBUG,
			LoggerWriter: os.Stdout,

			WithCaller:    true,
			WithTimestamp: timestamp.TimestampYYYYMonth,
		},
	)

	l.Debug("0")

	l.Debugf("%d", 1)
}

func Benchmark_Debug(b *testing.B) {
	l := NewLogger(
		&ParamsNewLogger{
			LoggerLevel: LevelDEBUG,

			WithCaller:    true,
			WithTimestamp: timestamp.TimestampYYYYMonth,
		},
	)

	b.ResetTimer()

	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				l.Info("1")
			}
		},
	)
}
