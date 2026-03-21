package log

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/log/bytearena"
	"github.com/tudorhulban/log/timestamp"
)

func TestLogger_Print(t *testing.T) {
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

func TestLogger_NoTimestampPrint(t *testing.T) {
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

func TestLogger_NanoPrint(t *testing.T) {
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

	time.Sleep(10 * time.Millisecond)

	cancel()
	<-chIngestionEnd

	fmt.Println(
		writer.String(),
	)
}

func TestLogger_YYYYPrint(t *testing.T) {
	writer := bytes.Buffer{}

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

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

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

	time.Sleep(10 * time.Millisecond)

	cancel()
	<-chIngestionEnd

	fmt.Println(
		writer.String(),
	)
}

func TestLogger_JSON_Print_With_Timestamp(t *testing.T) {
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

			WithJSON:  true,
			WithColor: true,
		},
	)
	require.NoError(t, errCrLogger)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	payload := "xxx"

	l.Printf("%s", payload)

	cancel()
	<-chIngestionEnd

	require.Contains(t, writer.String(), payload)

	fmt.Println(
		writer.String(),
	)
}

func TestLogger_JSON_Print_No_Timestamp(t *testing.T) {
	writer := bytes.Buffer{}

	ingestor := bytearena.NewIngestor(
		bytearena.Size100K,
		&writer,
	)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelDEBUG,

			WithJSON:  true,
			WithColor: true,
		},
	)
	require.NoError(t, errCrLogger)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	payload := "xxx"

	l.Printf("%s", payload)

	cancel()
	<-chIngestionEnd

	require.Contains(t, writer.String(), payload)

	fmt.Println(
		writer.String(),
	)
}
