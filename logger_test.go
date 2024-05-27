package log

import (
	"os"
	"testing"

	"github.com/TudorHulban/log/timestamp"
)

func TestLogger(t *testing.T) {
	l := NewLogger(
		&ParamsNewLogger{
			LoggerLevel:   LevelDEBUG,
			LoggerWriter:  os.Stdout,
			WithTimestamp: timestamp.TimestampNano,
		},
	)

	go l.PrintMessage("xxxx")

	l.PrintMessage("xxxx")
}
