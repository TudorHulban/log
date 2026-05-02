package log

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/log/timestamp"
)

func TestContext(t *testing.T) {
	writer := os.Stdout

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		writer,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	serviceLogging, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelDebug,

			WithFatalWriter: writer,
			WithTimestamp:   timestamp.TimestampRFC3339Bucharest,
		},
	)
	require.NoError(t, errCrLogger)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	f := NewLogContext(serviceLogging).
		WithRoot("service", "auth").
		SetInt("req_id", 12345).
		SetBool("cache_hit", true)

	require.NotNil(t, f.cfg.Load().root)

	f.Print("login ok")

	f.SetString("area", "some area")
	f.Print("login ok again")

	f.With("xxx", 2).
		Error().
		Msg("")
	f.With("yyy", 3).Msg("")
	f.With("zzzz", 4.3).
		Info().
		WithGoroutineID().
		Msg("")

	cancel()
	<-chIngestionEnd
}

func Test_With_JSON_Context(t *testing.T) {
	writer := os.Stdout

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		writer,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	serviceLogging, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelDebug,

			WithFatalWriter: writer,
			WithTimestamp:   timestamp.TimestampRFC3339Bucharest,
			WithJSON:        true,
		},
	)
	require.NoError(t, errCrLogger)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	f := NewLogContext(serviceLogging).
		WithRoot("service", "auth").
		SetInt("req_id", 12345).
		SetBool("cache_hit", true)

	require.NotNil(t, f.cfg.Load().root)

	f.SetString("area", "some area")

	f.With("xxx", 2).Msg("1")
	f.With("yyy", 3).Msg("2")
	f.With("zzzz", 4.3).Msg("3")

	cancel()
	<-chIngestionEnd
}
