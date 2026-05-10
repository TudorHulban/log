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
			WithColor:       true,
			WithJSON:        true,
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

	// produces
	// {"ts":"2026-05-10T13:08:07+03:00","level":"TRACE","msg":"created logger, level TRACE"}
	// {"ts":"2026-05-10T13:08:07+03:00","service":"auth","req_id":12345,"cache_hit":true,"root ends":"here","msg":"xxx1"}
	// {"ts":"2026-05-10T13:08:07+03:00","service":"auth","req_id":12345,"cache_hit":true,"root ends":"here","area":"some area","msg":"login ok again"}
	// {"ts":"2026-05-10T13:08:07+03:00","level":"INFO","service":"auth","req_id":12345,"cache_hit":true,"root ends":"here","area":"some area","zzzz":4.299999999999,"g":8,"msg":"finished"}
	// {"ts":"2026-05-10T13:08:07+03:00","level":"ERROR","service":"auth","req_id":12345,"cache_hit":true,"root ends":"here","area":"some area","xxxxxxxxxxxxx":"2","g":10,"msg":"some error"}

	// Allow time for the goroutine and the ingestor tick
	time.Sleep(100 * time.Millisecond)

	cancel()
	<-chIngestionEnd

	// --- Processing Lines ---
	require.Zero(t, bufFatal.Len())

	fmt.Println(bufLogs.String())

	entries, errParse := newTestEntries(bufLogs.String())
	require.NoError(t, errParse)

	require.Len(t, entries, 5)
	require.True(t, entries.haveTimestamp())
	require.NoError(t, entries.hasKey("msg", 5))
	require.NoError(t, entries.hasKey("req_id", 4))
	require.NoError(t, entries.hasKey("level", 3))
	require.NoError(t, entries.hasKeyWithValue("service", "auth", 4))
	require.NoError(t, entries.hasKeyWithValue("level", "TRACE", 1))
	require.NoError(t,
		entries.hasKeysWithValues(1,
			"req_id", 12345,
			"level", "INFO",
			"zzzz", 4.299999999999,
		),
	)
	require.NoError(t, entries.hasKeysWithValues(1, "req_id", 12345, "level", "ERROR"))
}
