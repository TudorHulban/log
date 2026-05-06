package log

import (
	"io"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/bytearena/helpers"
	"github.com/tudorhulban/log/timestamp"
)

// Benchmark_Debug-16    	20773839	        58.77 ns/op	       0 B/op	       0 allocs/op
func Benchmark_Debug(b *testing.B) {
	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		io.Discard,
	)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelDebug,

			WithFatalWriter: os.Stdout,
			WithTimestamp:   timestamp.TimestampRFC3339Bucharest,
			WithJSON:        true,
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

// Benchmark_Debugf-16    	20811402	        58.54 ns/op	       0 B/op	       0 allocs/op
func Benchmark_Debugf(b *testing.B) {
	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		io.Discard,
	)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelDebug,

			WithFatalWriter: os.Stdout,
			WithTimestamp:   timestamp.TimestampRFC3339Bucharest,
			WithJSON:        true,
		},
	)
	require.NoError(b, errCrLogger)

	b.SetParallelism(1)
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				l.Debugf("x=%d", 1)
			}
		},
	)
}

// Benchmark_Debugw-16    	30845853	        35.96 ns/op	       0 B/op	       0 allocs/op
func Benchmark_Debugw(b *testing.B) {
	writer := helpers.CountWriterNoBuffer{}

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		&writer,
	)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelDebug,

			WithFatalWriter: &writer,
			WithTimestamp:   timestamp.TimestampRFC3339Bucharest,
			WithJSON:        true,
		},
	)
	require.NoError(b, errCrLogger)

	runtime.GOMAXPROCS(1)

	b.ReportAllocs()
	b.ResetTimer()
	b.SetParallelism(16)

	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				l.Debugw("1some message", "some key", "some value")
			}
		},
	)
}
