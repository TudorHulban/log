package log

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/log/bytearena"
	"github.com/tudorhulban/log/helpers"
	"github.com/tudorhulban/log/timestamp"
)

func TestNanoPrint(t *testing.T) {
	sink := helpers.CountWriter{}

	ingestor := bytearena.NewIngestor(
		bytearena.Size100K,
		&sink,
	)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:      ingestor,
			LoggerLevel:   LevelDEBUG,
			WithTimestamp: timestamp.TimestampNano,
		},
	)
	require.NoError(t, errCrLogger)

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
	sink := helpers.CountWriter{}

	ingestor := bytearena.NewIngestor(
		bytearena.Size100K,
		&sink,
	)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:      ingestor,
			LoggerLevel:   LevelDEBUG,
			WithTimestamp: timestamp.TimestampYYYYMonth,
		},
	)
	require.NoError(t, errCrLogger)

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
