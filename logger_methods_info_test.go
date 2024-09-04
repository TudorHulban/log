package log

import (
	"os"
	"testing"

	"github.com/TudorHulban/log/timestamp"
)

func TestInfo(t *testing.T) {
	l := NewLogger(
		&ParamsNewLogger{
			LoggerLevel:  LevelDEBUG,
			LoggerWriter: os.Stdout,

			WithCaller:    true,
			WithTimestamp: timestamp.TimestampYYYYMonth,
			WithColor:     true,
		},
	)

	l.Info("0")

	l.Infof("%d", 1)
}

// profile
// go test -bench=Benchmark_Info_Logger -run=^$ . -cpuprofile profile.out
func Benchmark_Info_Logger(b *testing.B) {
	l := NewLogger(
		&ParamsNewLogger{
			WithTimestamp: timestamp.TimestampNil,
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
