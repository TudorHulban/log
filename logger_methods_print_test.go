package log

import (
	"os"
	"testing"
	"time"

	"github.com/tudorhulban/log/timestamp"
)

func TestNanoPrint(t *testing.T) {
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

func TestYYYYPrint(t *testing.T) {
	l := NewLogger(
		&ParamsNewLogger{
			LoggerLevel:   LevelDEBUG,
			LoggerWriter:  os.Stdout,
			WithTimestamp: timestamp.TimestampYYYYMonth,
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
