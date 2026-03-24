package log

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/log/bytearena"
	"github.com/tudorhulban/log/timestamp"
)

func TestContext(t *testing.T) {
	writer := os.Stdout

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K,
		writer,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	serviceLogging, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:      ingestor,
			LoggerLevel:   Level(LevelDEBUG),
			WithTimestamp: timestamp.TimestampRFC3339Bucharest,
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

	f.With("xxx", 2).Print()
	f.With("yyy", 3).Print()
	f.With("zzzz", 4.3).Print()

	cancel()
	<-chIngestionEnd
}

func Test_With_JSON_Context(t *testing.T) {
	writer := os.Stdout

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K,
		writer,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	serviceLogging, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:      ingestor,
			LoggerLevel:   Level(LevelDEBUG),
			WithTimestamp: timestamp.TimestampRFC3339Bucharest,
			WithJSON:      true,
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

	f.With("xxx", 2).Print("1")
	f.With("yyy", 3).Print("2")
	f.With("zzzz", 4.3).Print("3")

	cancel()
	<-chIngestionEnd
}
