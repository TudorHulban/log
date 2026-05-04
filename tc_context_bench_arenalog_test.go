package log

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/bytearena/helpers"
	"github.com/tudorhulban/log/timestamp"
)

func TestArenalog_OneField(t *testing.T) {
	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		os.Stdout,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

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
	require.NoError(t, errCrLogger)

	logContext := NewLogContext(logger)

	entry := logContext.WithString("area", "some area")
	entry.Info().Msg("benchmark test")

	cancel()
	<-chIngestionEnd
}

// go test -run '^$' -bench '^BenchmarkArenalog_OneField$' -benchmem -memprofile=mem.prof -cpuprofile=cpu.prof

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkArenalog_OneField/G1-16 	18708489	        63.30 ns/op	       0 B/op	       0 allocs/op
// BenchmarkArenalog_OneField/G2-16 	 9617562	       123.1 ns/op	      13 B/op	       0 allocs/op
// BenchmarkArenalog_OneField/G3-16 	 9792391	       121.7 ns/op	      13 B/op	       0 allocs/op
// BenchmarkArenalog_OneField/G4-16 	 9654943	       123.4 ns/op	      13 B/op	       0 allocs/op
func BenchmarkArenalog_OneField(b *testing.B) {
	gomaxprocsValues := []int{1, 2, 3, 4}
	writer := helpers.CountWriterNoBuffer{}

	for _, g := range gomaxprocsValues {
		b.Run(
			fmt.Sprintf("G%d", g),
			func(b *testing.B) {
				prev := runtime.GOMAXPROCS(g)
				defer runtime.GOMAXPROCS(prev)

				ingestor, errCrIngestor := bytearena.NewIngestor(
					bytearena.Size100K(),
					&writer,
				)
				require.NoError(b, errCrIngestor)
				require.NotNil(b, ingestor)

				ctx, cancel := context.WithCancel(context.Background())
				chIngestionEnd := ingestor.StartIngestion(ctx)

				logger, errCrLogger := NewLogger(
					&ParamsNewLogger{
						Ingestor:    ingestor,
						LoggerLevel: LevelDebug,

						WithFatalWriter: os.Stdout,
						WithTimestamp:   timestamp.TimestampRFC3339,
						WithJSON:        true,
					},
				)
				require.NoError(b, errCrLogger)

				logContext := NewLogContext(logger).
					WithRoot("service", "auth")

				runtime.GC()
				for i := 0; i < runtime.GOMAXPROCS(0)*4; i++ {
					e := entryPool.Get().(*Entry)
					entryPool.Put(e)
				}

				b.ReportAllocs()
				b.ResetTimer()

				for b.Loop() {
					entry := logContext.WithString("area", "some area")
					entry.Info().Msg("benchmark test")
				}

				cancel()
				<-chIngestionEnd

				require.NotZero(b,
					writer.TotalBytesWritten.Load(),
				)
			},
		)
	}
}

func TestArenalog_MultipleFields(t *testing.T) {
	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		os.Stdout,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelInfo,

			WithFatalWriter: os.Stdout,
			WithTimestamp:   timestamp.TimestampRFC3339,
			// WithJSON:        true,
		},
	)
	require.NoError(t, errCrLogger)

	logContext := NewLogContext(logger).
		WithRoot("service", "auth").
		SetInt("req_id", 12345).
		SetBool("cache_hit", true)

	entry := logContext.
		WithString("area", "some area").
		Info().
		WithString("user", "tudor").
		WithInt("attempt", 1).
		WithFloat("some float", 1.1137).
		WithBool("success", true)

	entry.Msg("benchmark test")

	// {"ts":"2026-05-04T14:29:35Z","level":"INFO","msg":"created logger, level INFO"}
	// {"ts":"2026-05-04T14:29:35Z","level":"INFO","service":"auth","req_id":12345,"cache_hit":true,"area":"some area","user":"tudor","attempt":1,"some float":1.113699999999,"success":true,"message":"benchmark test"}

	cancel()
	<-chIngestionEnd
}

// go test -run '^$' -bench '^BenchmarkContext_NoJSON_MultipleFields$' -benchmem -memprofile=mem.prof -cpuprofile=cpu.prof
// go tool pprof -alloc_objects mem.out

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkContext_NoJSON_MultipleFields-16    	 9660744	       125.6 ns/op	       0 B/op	       0 allocs/op
func BenchmarkContext_NoJSON_MultipleFields(b *testing.B) {
	var writer helpers.NoopWriter

	ingestor, errCrIngestor := bytearena.NewIngestor(bytearena.Size100K(), &writer)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	defer func() {
		cancel()
		<-chIngestionEnd
	}()

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelDebug,

			WithFatalWriter: os.Stdout,
			WithTimestamp:   timestamp.TimestampRFC3339,
		},
	)
	require.NoError(b, errCrLogger)

	logContext := NewLogContext(logger).
		WithRoot("service", "auth").
		SetInt("req_id", 12345).
		SetBool("cache_hit", true)

	runtime.GC()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		// 1. Create request with several attributes
		entry := logContext.
			WithString("area", "some area").
			Info().
			WithString("user", "tudor").
			WithInt("attempt", int64(i)).
			WithFloat("some float", 1.1137).
			WithBool("success", true)

		// 2. Print
		entry.Msg("benchmark test")
	}
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkContext_WithJSON_MultipleFields-16    	 7178278	       165.8 ns/op	      14 B/op	       0 allocs/op
func BenchmarkContext_WithJSON_MultipleFields(b *testing.B) {
	var writer helpers.NoopWriter

	ingestor, errCrIngestor := bytearena.NewIngestor(bytearena.Size100K(), &writer)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	defer func() {
		cancel()
		<-chIngestionEnd
	}()

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelDebug,

			WithFatalWriter: os.Stdout,
			WithTimestamp:   timestamp.TimestampRFC3339,
			WithJSON:        true,
		},
	)
	require.NoError(b, errCrLogger)

	logContext := NewLogContext(logger).
		WithRoot("service", "auth").
		SetInt("req_id", 12345).
		SetBool("cache_hit", true)

	runtime.GC()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		// 1. Create request with several attributes
		entry := logContext.
			WithString("area", "some area").
			Info().
			WithString("user", "tudor").
			WithInt("attempt", int64(i)).
			WithFloat("some float", 1.1137).
			WithBool("success", true)

		// 2. Print
		entry.Msg("benchmark test")
	}
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkArenalog_MultipleFields_Parallel/gomaxprocs=1-16         	15683703	        77.47 ns/op	       0 B/op	       0 allocs/op
// BenchmarkArenalog_MultipleFields_Parallel/gomaxprocs=2-16         	14636638	       108.6 ns/op	       0 B/op	       0 allocs/op
// BenchmarkArenalog_MultipleFields_Parallel/gomaxprocs=3-16         	14384590	        93.58 ns/op	       0 B/op	       0 allocs/op
// BenchmarkArenalog_MultipleFields_Parallel/gomaxprocs=4-16         	13956750	        85.85 ns/op	       0 B/op	       0 allocs/op
func BenchmarkArenalog_MultipleFields_Parallel(b *testing.B) {
	gomaxprocsValues := []int{1, 2, 3, 4}

	writer := helpers.CountWriterNoBuffer{}

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		&writer,
	)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	// keep ingestion running for the whole benchmark
	defer func() {
		cancel()
		<-chIngestionEnd
	}()

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:        ingestor,
			LoggerLevel:     LevelDebug,
			WithFatalWriter: os.Stdout,
			WithTimestamp:   timestamp.TimestampRFC3339,
			WithJSON:        true,
		},
	)
	require.NoError(b, errCrLogger)

	logContext := NewLogContext(logger).
		WithRoot("service", "auth").
		SetInt("req_id", 12345).
		SetBool("cache_hit", true)

	runtime.GC()

	b.SetParallelism(16)

	for _, g := range gomaxprocsValues {
		inner := g

		b.Run(
			fmt.Sprintf("gomaxprocs=%d", inner),
			func(b *testing.B) {
				prev := runtime.GOMAXPROCS(inner)
				defer runtime.GOMAXPROCS(prev)

				b.ReportAllocs()
				b.ResetTimer()

				b.RunParallel(
					func(pb *testing.PB) {
						i := int64(0)
						for pb.Next() {
							entry := logContext.
								WithString("area", "some area").
								Info().
								WithString("user", "tudor").
								WithInt("attempt", i).
								WithFloat("some float", 1.1137).
								WithBool("success", true)

							entry.Msg("benchmark test")

							i++
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
