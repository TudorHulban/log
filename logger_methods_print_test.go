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
			WithTimestamp: timestamp.TimestampNano,
		},
	)

	go l.PrintMessage("xxx1")
	go l.PrintMessage("xxx2")
	go l.PrintMessage("xxx3")

	l.Printw(
		"message:",
		[]string{
			"x1",
			"x2",
		},
		"x3",
	)

	time.Sleep(1 * time.Second)
}

func Benchmark_Print_Logger(b *testing.B) {
	l := NewLogger(
		&ParamsNewLogger{
			WithTimestamp: timestamp.TimestampNil,
		},
	)

	b.ResetTimer()

	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				l.PrintMessage("1")
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
