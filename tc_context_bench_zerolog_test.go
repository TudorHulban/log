package log

import (
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena/helpers"
)

func TestZerolog_OneField(t *testing.T) {
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()

	// {"level":"info","area":"some area","time":"2026-04-28T17:33:32+03:00","message":"benchmark test"}

	logger.Info().
		Str("area", "some area").
		Msg("benchmark test")
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkZerolog_OneField/G1-16 	 7227042	       166.4 ns/op	       0 B/op	       0 allocs/op
// BenchmarkZerolog_OneField/G2-16 	 7211571	       165.2 ns/op	       0 B/op	       0 allocs/op
func BenchmarkZerolog_OneField(b *testing.B) {
	gomaxprocsValues := []int{1, 2}
	writer := helpers.CountWriterNoBuffer{}

	for _, g := range gomaxprocsValues {
		b.Run(
			fmt.Sprintf("G%d", g),
			func(b *testing.B) {
				prev := runtime.GOMAXPROCS(g)
				defer runtime.GOMAXPROCS(prev)

				logger := zerolog.New(&writer).With().
					Timestamp().
					Logger()

				b.ReportAllocs()
				b.ResetTimer()

				for b.Loop() {
					logger.Info().
						Str("area", "some area").
						Msg("benchmark test")
				}

				require.NotZero(b, writer.TotalBytesWritten.Load())
			},
		)
	}
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
