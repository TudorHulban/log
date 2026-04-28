package log

import (
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/phuslu/log"
	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena/helpers"
)

func TestPhuslu_OneField(t *testing.T) {
	logger := log.Logger{
		Level:      log.InfoLevel,
		TimeFormat: time.RFC3339,
		Writer: &log.IOWriter{
			Writer: os.Stdout,
		},
	}

	// {"time":"2026-04-28T12:09:35+03:00","level":"info","area":"some area","message":"benchmark test"}

	logger.Info().
		Str("area", "some area").
		Msg("benchmark test")
}

// go test -run '^$' -bench '^BenchmarkPhuslu_OneField$' -benchmem

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkPhuslu_OneField/G1-16 	 8996448	       134.4 ns/op	       0 B/op	       0 allocs/op
// BenchmarkPhuslu_OneField/G2-16 	 9046921	       132.7 ns/op	       0 B/op	       0 allocs/op

func BenchmarkPhuslu_OneField(b *testing.B) {
	gomaxprocsValues := []int{1, 2}
	writer := helpers.CountWriterNoBuffer{}

	for _, g := range gomaxprocsValues {
		b.Run(
			fmt.Sprintf("G%d", g),
			func(b *testing.B) {
				prev := runtime.GOMAXPROCS(g)
				defer runtime.GOMAXPROCS(prev)

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

				require.NotZero(b, writer.TotalBytesWritten.Load())
			},
		)
	}
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkPhuslu_Parallel_OneField/gomaxprocs=1-16         	 8888434	       131.1 ns/op	       0 B/op	       0 allocs/op
// BenchmarkPhuslu_Parallel_OneField/gomaxprocs=2-16         	16020363	        74.65 ns/op	       0 B/op	       0 allocs/op
// BenchmarkPhuslu_Parallel_OneField/gomaxprocs=3-16         	23736805	        50.80 ns/op	       0 B/op	       0 allocs/op
// BenchmarkPhuslu_Parallel_OneField/gomaxprocs=4-16         	23618685	        50.08 ns/op	       0 B/op	       0 allocs/op
func BenchmarkPhuslu_Parallel_OneField(b *testing.B) {
	gomaxprocsValues := []int{1, 2, 3, 4}

	writer := helpers.CountWriterNoBuffer{}

	logger := log.Logger{
		Level:      log.InfoLevel,
		TimeFormat: time.RFC3339,
		Writer:     &log.IOWriter{Writer: &writer},
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.SetParallelism(16)

	for _, g := range gomaxprocsValues {
		inner := g

		b.Run(
			fmt.Sprintf("gomaxprocs=%d", g),
			func(b *testing.B) {
				prev := runtime.GOMAXPROCS(inner)
				defer runtime.GOMAXPROCS(prev)

				b.RunParallel(
					func(pb *testing.PB) {
						for pb.Next() {
							logger.Info().
								Str("area", "some area").
								Msg("benchmark test")
						}
					},
				)
			},
		)
	}

	require.NotZero(
		b,
		writer.TotalBytesWritten.Load(),
		"1. writer must record bytes",
	)
}

// BenchmarkPhuslu_WithFields-16    	 6631579	       179.0 ns/op	       0 B/op	       0 allocs/op
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
