package log

import (
	"os"
	"testing"
	"time"

	"github.com/TudorHulban/log/timestamp"
)

func TestPrint(t *testing.T) {
	l := NewLogger(
		&ParamsNewLogger{
			LoggerLevel:   LevelDEBUG,
			LoggerWriter:  os.Stdout,
			WithTimestamp: timestamp,
		},
	)

	go l.PrintLocal("xxx1")
	go l.PrintLocal("xxx2")
	go l.PrintLocal("xxx3")

	time.Sleep(1 * time.Second)
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

func Benchmark_Local_TimestampNano_Logger(b *testing.B) {
	logger := NewLogger(
		&ParamsNewLogger{
			WithTimestamp: timestamp.TimestampNano,
		},
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
