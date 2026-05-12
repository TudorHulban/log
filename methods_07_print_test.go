package log

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/log/query"
	"github.com/tudorhulban/log/timestamp"
)

// test produces
// {"ts":"2026-05-12T15:51:49+03:00","level":"TRACE","msg":"created logger, level TRACE"}
// {"ts":"2026-05-12T15:51:49+03:00","level":"TRACE","msg":"xxx"}
// {"ts":"2026-05-12T15:51:49+03:00","level":"TRACE","msg":"LOG_ERR(odd_args): xxx","key":"(MISSING)"}

func TestErrorsPrint(t *testing.T) {
	var buf bytes.Buffer

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		&buf,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelTrace,

			WithFatalWriter: &buf,
			WithTimestamp:   timestamp.TimestampRFC3339Bucharest,
			WithJSON:        true,
		},
	)
	require.NoError(t, errCrLogger)
	require.NotNil(t, l)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	// 1. Valid call - No extra keys
	l.Printw("xxx")

	// 2. Invalid call - Odd number of keys (1 key, 0 values)
	// This should trigger the internal "panicw" log line
	l.Printw("xxx", "key")

	cancel()
	<-chIngestionEnd

	// Parse the results from the buffer
	logSet, errCr := query.NewLogset(buf.String())
	require.NoError(t, errCr)

	require.Len(t, logSet, 3)

	require.NoError(t, logSet.HasKeyWithValue("msg", "xxx", 1))
	require.NoError(t, logSet.HasKey("msg", 3))
	require.NoError(t, logSet.HasKeyWithValue("key", _Missing, 1))
	require.NoError(t, logSet.HasKey("level", 3))
}
