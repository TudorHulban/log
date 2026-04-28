package log

import (
	"context"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/bytearena/helpers"
	"github.com/tudorhulban/log/timestamp"
)

func BenchmarkLogger_Entry(b *testing.B) {
	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		&helpers.NoopWriter{},
	)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelInfo,

			WithFatalWriter: os.Stdout,
			WithTimestamp:   timestamp.TimestampRFC3339,
			WithJSON:        true,
		},
	)
	require.NoError(b, errCrLogger)
	require.NotNil(b, logger)

	logContext := NewLogContext(logger)
	require.NotNil(b, logContext.cfg.Load().root)

	runtime.GC()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		logger.Print("hi", 123, "world")
	}

	cancel()
	<-chIngestionEnd
}
