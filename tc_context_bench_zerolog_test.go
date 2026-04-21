package log

import (
	"runtime"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena/helpers"
)

// BenchmarkZerolog_Serial_OneField-16    	 7467723	       158.3 ns/op	       0 B/op	       0 allocs/op
func BenchmarkZerolog_Serial_OneField(b *testing.B) {
	writer := helpers.CountWriterNoBuffer{}

	logger := zerolog.New(&writer).With().
		Timestamp().
		Logger()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		logger.Info().
			Str("area", "some area").
			Msg("benchmark test")
	}

	_ = writer.TotalBytesWritten.Load() // force writer to stay live
}

// BenchmarkZerolog_Parallel_OneField-16    	 7567926	       156.4 ns/op	       0 B/op	       0 allocs/op
func BenchmarkZerolog_Parallel_OneField(b *testing.B) {
	runtime.GOMAXPROCS(1)

	writer := helpers.CountWriterNoBuffer{}

	logger := zerolog.New(&writer).With().
		Timestamp().
		Logger()

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
		writer.TotalBytesWritten.Load(), // force writer to stay live
	)
}

// BenchmarkZerolog_WithFields-16    	 5978826	       200.9 ns/op	       0 B/op	       0 allocs/op
func BenchmarkZerolog_WithFields(b *testing.B) {
	var writer helpers.NoopWriter

	logger := zerolog.New(&writer).With().
		Timestamp().
		Str("service", "auth").
		Int("req_id", 12345).
		Bool("cache_hit", true).
		Logger()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		logger.Info().
			Str("area", "some area").
			Str("user", "tudor").
			Int("attempt", i).
			Msg("benchmark test")
	}
}
