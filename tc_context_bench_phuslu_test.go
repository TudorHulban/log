package log

import (
	"runtime"
	"testing"
	"time"

	"github.com/phuslu/log"
	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena/helpers"
)

// BenchmarkPhuslu_OneField-12    	 9210207	       129.9 ns/op	       0 B/op	       0 allocs/op
func BenchmarkPhuslu_OneField(b *testing.B) {
	writer := helpers.CountWriterNoBuffer{}

	logger := log.Logger{
		Level:      log.InfoLevel,
		TimeFormat: time.RFC3339,
		Writer:     &log.IOWriter{Writer: &writer},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		logger.Info().
			Str("area", "some area").
			Msg("benchmark test")
	}

	require.NotZero(b,
		writer.TotalBytesWritten.Load(),
	)
}

// BenchmarkPhuslu_Parallel_OneField-16    	 9290592	       129.4 ns/op	       0 B/op	       0 allocs/op
func BenchmarkPhuslu_Parallel_OneField(b *testing.B) {
	writer := helpers.CountWriterNoBuffer{}

	logger := log.Logger{
		Level:      log.InfoLevel,
		TimeFormat: time.RFC3339,
		Writer:     &log.IOWriter{Writer: &writer},
	}

	runtime.GOMAXPROCS(1)

	b.ReportAllocs()
	b.ResetTimer()
	b.SetParallelism(16)

	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().
					Str("area", "some area").
					Msg("benchmark test")
			}
		},
	)

	require.NotZero(b,
		writer.TotalBytesWritten.Load(),
	)
}

// BenchmarkPhuslu_WithFields-16    	 6391726	       187.9 ns/op	       0 B/op	       0 allocs/op
func BenchmarkPhuslu_WithFields(b *testing.B) {
	var writer helpers.NoopWriter

	logger := log.Logger{
		Level:      log.InfoLevel,
		TimeField:  "ts",
		TimeFormat: time.RFC3339,
		Writer:     &log.IOWriter{Writer: &writer},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		logger.Info().
			Str("service", "auth").
			Int("req_id", 12345).
			Str("area", "some area").
			Str("user", "tudor").
			Int("attempt", i).
			Bool("success", true).
			Msg("benchmark test")
	}
}
