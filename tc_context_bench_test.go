package log

import (
	"context"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/bytearena/helpers"
	"github.com/tudorhulban/log/timestamp"
)

// TODO; alloc?
// BenchmarkContext_NoJSON_OneField-16    	 9055264	       132.8 ns/op	      27 B/op	       1 allocs/op
func BenchmarkContext_NoJSON_OneField(b *testing.B) {
	writer := helpers.CountWriterNoBuffer{}

	ingestor, errCrIngestor := bytearena.NewIngestor(bytearena.Size100K(), &writer)
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
		},
	)
	require.NoError(b, errCrLogger)

	logContext := NewLogContext(logger).
		WithRoot("service", "auth")

	runtime.GC()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		// 1. Create request with 4 attributes
		entry := logContext.With("area", "some area")

		// 2. Print
		entry.Print("benchmark test")
	}

	cancel()
	<-chIngestionEnd

	require.NotZero(b,
		writer.TotalBytesWritten.Load(),
	)
}

// go test -run=^$ -bench=^BenchmarkContext_NoJSON_MultipleFields$ -benchmem -memprofile=mem.out
// go tool pprof -alloc_objects mem.out

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkContext_NoJSON_MultipleFields-16    	 6536090	       182.7 ns/op	      29 B/op	       1 allocs/op
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
		// 1. Create request with 4 attributes
		entry := logContext.With("area", "some area").
			WithString("user", "tudor").
			WithInt("attempt", i).
			WithBool("success", true)

		// 2. Print
		entry.Print("benchmark test")
	}
}

// BenchmarkContext_WithJSON_OneField-16    	 8952406	       134.6 ns/op	      31 B/op	       1 allocs/op
func BenchmarkContext_WithJSON_OneField(b *testing.B) {
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
		WithRoot("service", "auth")

	runtime.GC()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		// 1. Create request with 4 attributes
		entry := logContext.With("area", "some area")

		// 2. Print
		entry.Print("benchmark test")
	}
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkContext_WithJSON_MultipleFields-16    	 6506314	       185.1 ns/op	      35 B/op	       2 allocs/op
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
		// 1. Create request with 4 attributes
		entry := logContext.With("area", "some area").
			WithString("user", "tudor").
			WithInt("attempt", i).
			WithBool("success", true)

		// 2. Print
		entry.Print("benchmark test")
	}
}
