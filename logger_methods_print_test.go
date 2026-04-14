package log

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/log/timestamp"
)

func TestLogger_Print(t *testing.T) {
	writer := bytes.Buffer{}

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		&writer,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:      ingestor,
			LoggerLevel:   Level(LevelDEBUG),
			WithTimestamp: timestamp.TimestampRFC3339Bucharest,
		},
	)
	require.NoError(t, errCrLogger)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	payload := "xxx"

	l.PrintFast(payload)

	cancel()
	<-chIngestionEnd

	require.Contains(t, writer.String(), payload)

	fmt.Println(
		writer.String(),
	)
}

func TestLogger_NoTimestampPrint(t *testing.T) {
	writer := bytes.Buffer{}

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		&writer,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: Level(LevelDEBUG),
		},
	)
	require.NoError(t, errCrLogger)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	l.PrintWithNoTimestampFast("hi", 123, "world")

	cancel()
	<-chIngestionEnd

	require.Contains(t, writer.String(), "hi 123 world")

	fmt.Println(
		writer.String(),
	)
}

func TestLogger_NanoPrint(t *testing.T) {
	writer := bytes.Buffer{}

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		&writer,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:      ingestor,
			LoggerLevel:   Level(LevelDEBUG),
			WithTimestamp: timestamp.TimestampNano,
		},
	)
	require.NoError(t, errCrLogger)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	go l.PrintMessageFast("xxx1")
	go l.PrintMessageFast("xxx2")
	go l.PrintMessageFast("xxx3")

	l.PrintwFast(
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

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		&writer,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:      ingestor,
			LoggerLevel:   Level(LevelDEBUG),
			WithTimestamp: timestamp.TimestampYYYYMonth,
		},
	)
	require.NoError(t, errCrLogger)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	go l.PrintMessageFast("xxx1")
	go l.PrintMessageFast("xxx2")
	go l.PrintMessageFast("xxx3")

	l.PrintwFast(
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

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		&writer,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:      ingestor,
			LoggerLevel:   Level(LevelDEBUG),
			WithTimestamp: timestamp.TimestampNano,

			WithJSON:  true,
			WithColor: true,
		},
	)
	require.NoError(t, errCrLogger)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	payload := "xxx"

	l.PrintfFast("%s", payload)

	cancel()
	<-chIngestionEnd

	require.Contains(t, writer.String(), payload)

	fmt.Println(
		writer.String(),
	)
}

func TestLogger_JSON_Print_No_Timestamp(t *testing.T) {
	writer := bytes.Buffer{}

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		&writer,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: Level(LevelDEBUG),

			WithJSON:  true,
			WithColor: true,
		},
	)
	require.NoError(t, errCrLogger)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	payload := "xxx"

	l.PrintfFast("%s", payload)

	cancel()
	<-chIngestionEnd

	require.Contains(t, writer.String(), payload)

	fmt.Println(
		writer.String(),
	)
}
