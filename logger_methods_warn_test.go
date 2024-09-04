package log

import (
	"os"
	"testing"

	"github.com/TudorHulban/log/timestamp"
)

func TestWarn(t *testing.T) {
	l := NewLogger(
		&ParamsNewLogger{
			LoggerLevel:   LevelDEBUG,
			LoggerWriter:  os.Stdout,
			WithTimestamp: timestamp.TimestampNano,
			WithColor:     true,
		},
	)

	l.Warn("0")

	l.Warnf("%d", 1)
}

func BenchmarkLogger_Warn(b *testing.B) {
	l := NewLogger(
		&ParamsNewLogger{
			WithTimestamp: timestamp.TimestampNil,
		},
	)

	b.ResetTimer()

	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				l.Warn("1")
			}
		},
	)
}
