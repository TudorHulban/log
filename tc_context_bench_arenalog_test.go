package log

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/bytearena/helpers"
	"github.com/tudorhulban/log/timestamp"
)

// test produces
// {"ts":"2026-05-11T09:33:49.624Z","level":"INFO","msg":"created logger, level INFO"}
// {"ts":"2026-05-11T09:33:49.624Z","level":"INFO","area":"some area","msg":"benchmark test"}

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
// BenchmarkArenalog_OneField/G1-16 	18325016	        63.49 ns/op	       0 B/op	       0 allocs/op
// BenchmarkArenalog_OneField/G2-16 	11072150	       108.7 ns/op	       6 B/op	       0 allocs/op
// BenchmarkArenalog_OneField/G3-16 	10980145	       108.3 ns/op	       5 B/op	       0 allocs/op
// BenchmarkArenalog_OneField/G4-16 	10916209	       110.0 ns/op	       6 B/op	       0 allocs/op
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
					e, _ := entryPool.Get().(*Entry) //nolint:revive
					entryPool.Put(e)
				}

				var warmupBuffer []byte
				timestamp.TimestampRFC3339(warmupBuffer)

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

// test produces
// {"ts":"2026-05-04T14:29:35Z","level":"INFO","msg":"created logger, level INFO"}
// {"ts":"2026-05-04T14:29:35Z","level":"INFO","service":"auth","req_id":12345,"cache_hit":true,"area":"some area","user":"tudor","attempt":1,"some float":1.113699999999,"success":true,"message":"benchmark test"}

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

	cancel()
	<-chIngestionEnd
}

// go test -run '^$' -bench '^BenchmarkContext_NoJSON_MultipleFields$' -benchmem -memprofile=mem.prof -cpuprofile=cpu.prof
// go tool pprof -alloc_objects mem.out

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkContext_NoJSON_MultipleFields-16    	 9760326	       126.5 ns/op	       4 B/op	       0 allocs/op
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

	for i := 0; i < runtime.GOMAXPROCS(0)*4; i++ {
		e, _ := entryPool.Get().(*Entry) //nolint:revive
		entryPool.Put(e)
	}

	var warmupBuffer []byte
	timestamp.TimestampRFC3339(warmupBuffer)

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
// BenchmarkContext_WithJSON_MultipleFields-16    	 8048110	       151.1 ns/op	       5 B/op	       0 allocs/op
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

	for i := 0; i < runtime.GOMAXPROCS(0)*4; i++ {
		e, _ := entryPool.Get().(*Entry) //nolint:revive
		entryPool.Put(e)
	}

	var warmupBuffer []byte
	timestamp.TimestampRFC3339(warmupBuffer)

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
// BenchmarkArenalog_MultipleFields_Parallel/gomaxprocs=1-16         	16213712	        74.77 ns/op	       0 B/op	       0 allocs/op
// BenchmarkArenalog_MultipleFields_Parallel/gomaxprocs=2-16         	17813103	        64.67 ns/op	       0 B/op	       0 allocs/op
// BenchmarkArenalog_MultipleFields_Parallel/gomaxprocs=3-16         	14062590	        84.90 ns/op	       0 B/op	       0 allocs/op
// BenchmarkArenalog_MultipleFields_Parallel/gomaxprocs=4-16         	13590536	        88.83 ns/op	       0 B/op	       0 allocs/op
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

	// warm up
	for i := 0; i < runtime.GOMAXPROCS(0)*4; i++ {
		e, _ := entryPool.Get().(*Entry) //nolint:revive
		entryPool.Put(e)
	}

	var warmupBuffer []byte
	timestamp.TimestampRFC3339(warmupBuffer)

	time.Sleep(10 * time.Millisecond)

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
