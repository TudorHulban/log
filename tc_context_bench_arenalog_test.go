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

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkArenalog_OneField/G1-16 	24744356	        44.83 ns/op	       0 B/op	       0 allocs/op
// BenchmarkArenalog_OneField/G2-16 	 9803336	       126.5 ns/op	      16 B/op	       0 allocs/op
// BenchmarkArenalog_OneField/G3-16 	 9544459	       124.4 ns/op	      16 B/op	       0 allocs/op
// BenchmarkArenalog_OneField/G4-16 	 9686191	       125.5 ns/op	      16 B/op	       0 allocs/op
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
			WithJSON:        true,
		},
	)
	require.NoError(t, errCrLogger)

	logContext := NewLogContext(logger)

	entry := logContext.WithString("area", "some area")
	entry.Info().
		WithFloat("some float", 1.1137).
		Msg("benchmark test")

	cancel()
	<-chIngestionEnd
}

// go test -run=^$ -bench=^BenchmarkContext_NoJSON_MultipleFields$ -benchmem -memprofile=mem.out
// go tool pprof -alloc_objects mem.out

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkContext_NoJSON_MultipleFields-16    	 7380024	       162.3 ns/op	     432 B/op	       2 allocs/op
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
		entry := logContext.With("area", "some area").
			WithString("user", "tudor").
			WithInt("attempt", int64(i)).
			WithFloat("some float", 1.1137).
			WithBool("success", true)

		// 2. Print
		entry.Msg("benchmark test")
	}
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkContext_WithJSON_MultipleFields-16    	 6539997	       185.3 ns/op	     496 B/op	       2 allocs/op
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
		entry := logContext.With("area", "some area").
			WithString("user", "tudor").
			WithInt("attempt", int64(i)).
			WithFloat("some float", 1.1137).
			WithBool("success", true)

		// 2. Print
		entry.Msg("benchmark test")
	}
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkArenalog_SeveralFields_Parallel/gomaxprocs=1-16         	 7205157	       167.1 ns/op	     496 B/op	       2 allocs/op
// BenchmarkArenalog_SeveralFields_Parallel/gomaxprocs=2-16         	10468606	       116.1 ns/op	     496 B/op	       2 allocs/op
// BenchmarkArenalog_SeveralFields_Parallel/gomaxprocs=3-16         	10905969	       110.2 ns/op	     496 B/op	       2 allocs/op
// BenchmarkArenalog_SeveralFields_Parallel/gomaxprocs=4-16         	11598610	       103.1 ns/op	     496 B/op	       2 allocs/op
func BenchmarkArenalog_SeveralFields_Parallel(b *testing.B) {
	gomaxprocsValues := []int{1, 2, 3, 4}

	writer := helpers.CountWriterNoBuffer{}

	ingestor, errCrIngestor := bytearena.NewIngestor(bytearena.Size100K(), &writer)
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
							entry := logContext.With("area", "some area").
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
