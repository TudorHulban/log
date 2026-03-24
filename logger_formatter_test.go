package log

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/log/bytearena"
	"github.com/tudorhulban/log/timestamp"
)

func TestFormater(t *testing.T) {
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

	f := NewFormatter(serviceLogging).
		WithRoot("service", "auth").
		WithInt("req_id", 12345).
		WithBool("cache_hit", true)

	require.NotNil(t, f.cfg.Load().root)

	f.Print("login ok")

	f.WithString("area", "some area")
	f.Print("login ok again")

	cancel()
	<-chIngestionEnd
}
