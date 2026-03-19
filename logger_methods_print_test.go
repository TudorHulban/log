package log

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/log/bytearena"
	"github.com/tudorhulban/log/helpers"
	"github.com/tudorhulban/log/timestamp"
)

func TestPrint(t *testing.T) {
	writer := bytes.Buffer{}

	ingestor := bytearena.NewIngestor(
		bytearena.Size100K,
		&writer,
	)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:      ingestor,
			LoggerLevel:   LevelDEBUG,
			WithTimestamp: timestamp.TimestampNano,
		},
	)
	require.NoError(t, errCrLogger)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	payload := "xxx"

	l.Print(payload)

	cancel()
	<-chIngestionEnd

	require.Contains(t, writer.String(), payload)
}

func TestNoTimestampPrint(t *testing.T) {
	writer := bytes.Buffer{}

	ingestor := bytearena.NewIngestor(
		bytearena.Size100K,
		&writer,
	)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelDEBUG,
		},
	)
	require.NoError(t, errCrLogger)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	l.PrintWithNoTimestamp("hi", 123, "world")

	cancel()
	<-chIngestionEnd

	require.Contains(t, writer.String(), "hi 123 world")

	fmt.Println(
		writer.String(),
	)
}

func TestNanoPrint(t *testing.T) {
	writer := helpers.CountWriter{}

	ingestor := bytearena.NewIngestor(
		bytearena.Size100K,
		&writer,
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
	writer := helpers.CountWriter{}

	ingestor := bytearena.NewIngestor(
		bytearena.Size100K,
		&writer,
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
