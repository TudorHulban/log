package log

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/log/query"
	"github.com/tudorhulban/log/timestamp"
)

// test produces
// {"ts":"2026-05-10T13:08:07+03:00","level":"TRACE","msg":"created logger, level TRACE"}
// {"ts":"2026-05-10T13:08:07+03:00","service":"auth","req_id":12345,"cache_hit":true,"root ends":"here","msg":"xxx1"}
// {"ts":"2026-05-10T13:08:07+03:00","service":"auth","req_id":12345,"cache_hit":true,"root ends":"here","area":"some area","msg":"login ok again"}
// {"ts":"2026-05-10T13:08:07+03:00","level":"INFO","service":"auth","req_id":12345,"cache_hit":true,"root ends":"here","area":"some area","zzzz":4.299999999999,"g":8,"msg":"finished"}
// {"ts":"2026-05-10T13:08:07+03:00","level":"ERROR","service":"auth","req_id":12345,"cache_hit":true,"root ends":"here","area":"some area","xxxxxxxxxxxxx":"2","g":10,"msg":"some error"}

func TestContextPrint(t *testing.T) {
	var bufLogs, bufFatal bytes.Buffer

	writer := &bufLogs

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		writer,
		bytearena.WithTickIfDataMilliseconds(10),
	)
	require.NoError(t, errCrIngestor)

	serviceLogging, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelTrace,

			WithFatalWriter: &bufFatal,
			WithTimestamp:   timestamp.TimestampRFC3339Bucharest,
			// WithColor:       true,
			WithJSON: true,
		},
	)
	require.NoError(t, errCrLogger)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	f := NewLogContext(serviceLogging).
		WithRoot("service", "auth").
		SetInt("req_id", 12345).
		SetBool("cache_hit", true).
		SetString("root ends", "here")

	// Execution
	f.Print("xxx1")
	f.SetString("area", "some area")
	f.Print("login ok again")

	go func() {
		f.With("xxxxxxxxxxxxx", "2").
			Error().
			WithGoroutineID().
			Msg("some error")
	}()

	f.With("zzzz", 4.3).
		Info().
		WithGoroutineID().
		Msg("finished")

	// Allow time for the goroutine and the ingestor tick
	time.Sleep(100 * time.Millisecond)

	cancel()
	<-chIngestionEnd

	// --- Processing Lines ---
	require.Zero(t, bufFatal.Len())

	logSet, errParse := query.NewLogset(bufLogs.String())
	require.NoError(t, errParse)

	require.Len(t,
		logSet,
		5,

		logSet,
	)

	require.Len(t, logSet.WithTimestamp(), 5)
	require.NoError(t, logSet.HasKey("msg", 5))
	require.NoError(t, logSet.HasKey("req_id", 4))
	require.NoError(t, logSet.HasKey("level", 3))
	require.NoError(t, logSet.HasKeyWithValue("service", "auth", 4))
	require.NoError(t, logSet.HasKeyWithValue("level", "TRACE", 1))
	require.NoError(t,
		logSet.HasKeysWithValues(1,
			"req_id", 12345,
			"level", "INFO",
			"zzzz", 4.299999999999,
		),
	)
	require.NoError(t,
		logSet.HasKeysWithValues(1,
			"req_id", 12345,
			"level", "ERROR",
		),
	)
}
