package log

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/log/timestamp"
)

func TestDebug(t *testing.T) {
	var writer bytes.Buffer

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
			WithCaller:    true,
			WithColor:     false, // JSON + color = messy output
			WithJSON:      true,
		},
	)
	require.NoError(t, errCrLogger)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	l.Debug("0")
	l.Debugf("%d", 1)

	cancel()
	<-chIngestionEnd

	out := writer.String()
	require.NotEmpty(t, out)

	// JSON mode → each log entry is a JSON object on its own line
	lines := strings.Split(strings.TrimSpace(out), "\n")
	require.Len(t, lines, 2)

	// First log: Debug("0")
	require.Contains(t, lines[0], `"level":"debug"`)
	require.Contains(t, lines[0], `"msg":"0"`)

	// Second log: Debugf("%d", 1)
	require.Contains(t, lines[1], `"level":"debug"`)
	require.Contains(t, lines[1], `"msg":"1"`)

	// Timestamp must exist in both
	require.Contains(t, lines[0], `"ts":`)
	require.Contains(t, lines[1], `"ts":`)

	// Caller info must exist if enabled
	require.Contains(t, lines[0], `"caller":`)
	require.Contains(t, lines[1], `"caller":`)
}

// Benchmark_Debug-16    	  741261	      1905 ns/op	    5256 B/op	      19 allocs/op
func Benchmark_Debug(b *testing.B) {
	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		io.Discard,
	)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:      ingestor,
			LoggerLevel:   Level(LevelDEBUG),
			WithTimestamp: timestamp.TimestampRFC3339Bucharest,
			WithJSON:      true,
		},
	)
	require.NoError(b, errCrLogger)

	b.SetParallelism(1)
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				l.Debug("1")
			}
		},
	)
}

// Benchmark_Debug_Fast-16    	21640418	        57.29 ns/op	       0 B/op	       0 allocs/op
func Benchmark_Debug_Fast(b *testing.B) {
	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		io.Discard,
	)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:      ingestor,
			LoggerLevel:   Level(LevelDEBUG),
			WithTimestamp: timestamp.TimestampRFC3339Bucharest,
			WithJSON:      true,

			EstimatedMessageSize: MessageLargeSize,
		},
	)
	require.NoError(b, errCrLogger)

	b.SetParallelism(1)
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				l.DebugFast("1")
			}
		},
	)
}
