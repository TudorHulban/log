package log

import (
	"os"
	"testing"

	"github.com/TudorHulban/log/timestamp"
)

func TestMix(t *testing.T) {
	l := NewLogger(
		&ParamsNewLogger{
			LoggerLevel:   LevelDEBUG,
			LoggerWriter:  os.Stdout,
			WithTimestamp: timestamp.TimestampNano,
		},
	)

	go l.Print("0")
	go l.Info("1")
	go l.Warn("2")
	l.Debug("3")
}
