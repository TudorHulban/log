package log

import (
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/log/bytearena"
	"github.com/tudorhulban/log/helpers"
	"github.com/tudorhulban/log/timestamp"

	"github.com/phuslu/log"
	phuslog "github.com/phuslu/log"
)

// BenchmarkContext_NoJSON_OneField-12    	24259626	        48.80 ns/op	       4 B/op	       0 allocs/op
func BenchmarkContext_NoJSON_OneField(b *testing.B) {
	var writer helpers.NoopWriter

	ingestor, errCrIngestor := bytearena.NewIngestor(bytearena.Size100K, &writer)
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
			Ingestor:      ingestor,
			LoggerLevel:   Level(LevelDEBUG),
			WithTimestamp: timestamp.TimestampRFC3339,
		},
	)
	require.NoError(b, errCrLogger)

	logContext := NewLogContext(logger).
		WithRoot("service", "auth")

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		// 1. Create request with 4 attributes
		entry := logContext.With("area", "some area")

		// 2. Print
		entry.Print("benchmark test")
	}
}

// go test -bench=BenchmarkContext_NoJSON_MultipleFields -benchmem -memprofile=mem.out
// go tool pprof -alloc_space ./your.test mem.out

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkContext_NoJSON_MultipleFields-12    	10104504	       120.2 ns/op	      15 B/op	       1 allocs/op
func BenchmarkContext_NoJSON_MultipleFields(b *testing.B) {
	var writer helpers.NoopWriter

	ingestor, errCrIngestor := bytearena.NewIngestor(bytearena.Size100K, &writer)
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
			Ingestor:      ingestor,
			LoggerLevel:   Level(LevelDEBUG),
			WithTimestamp: timestamp.TimestampRFC3339,
		},
	)
	require.NoError(b, errCrLogger)

	logContext := NewLogContext(logger).
		WithRoot("service", "auth").
		SetInt("req_id", 12345).
		SetBool("cache_hit", true)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		// 1. Create request with 4 attributes
		entry := logContext.With("area", "some area").
			With("user", "tudor").
			With("attempt", i).
			With("success", true)

		// 2. Print
		entry.Print("benchmark test")
	}
}

// BenchmarkContext_WithJSON_OneField-12    	17870000	        67.06 ns/op	       7 B/op	       0 allocs/op
func BenchmarkContext_WithJSON_OneField(b *testing.B) {
	var writer helpers.NoopWriter

	ingestor, errCrIngestor := bytearena.NewIngestor(bytearena.Size100K, &writer)
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
			Ingestor:      ingestor,
			LoggerLevel:   Level(LevelDEBUG),
			WithTimestamp: timestamp.TimestampRFC3339,
			WithJSON:      true,
		},
	)
	require.NoError(b, errCrLogger)

	logContext := NewLogContext(logger).
		WithRoot("service", "auth")

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
// BenchmarkContext_WithJSON_MultipleFields-12    	 9880800	       122.4 ns/op	      12 B/op	       0 allocs/op
func BenchmarkContext_WithJSON_MultipleFields(b *testing.B) {
	var writer helpers.NoopWriter

	ingestor, errCrIngestor := bytearena.NewIngestor(bytearena.Size100K, &writer)
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
			Ingestor:      ingestor,
			LoggerLevel:   Level(LevelDEBUG),
			WithTimestamp: timestamp.TimestampRFC3339,
			WithJSON:      true,
		},
	)
	require.NoError(b, errCrLogger)

	logContext := NewLogContext(logger).
		WithRoot("service", "auth").
		SetInt("req_id", 12345).
		SetBool("cache_hit", true)

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

// BenchmarkZerolog_OneField-12    	 7189495	       167.2 ns/op	       0 B/op	       0 allocs/op
func BenchmarkZerolog_OneField(b *testing.B) {
	var writer helpers.NoopWriter

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
}

// BenchmarkZerolog_WithFields-12    	 5937602	       198.6 ns/op	       0 B/op	       0 allocs/op
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

// BenchmarkPhuslu_OneField-12    	 9210207	       129.9 ns/op	       0 B/op	       0 allocs/op
func BenchmarkPhuslu_OneField(b *testing.B) {
	var writer helpers.NoopWriter

	logger := phuslog.Logger{
		Level:      phuslog.InfoLevel,
		TimeFormat: time.RFC3339,
		Writer:     &phuslog.IOWriter{Writer: &writer},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		logger.Info().
			Str("area", "some area").
			Msg("benchmark test")
	}
}

// BenchmarkPhuslu_WithFields-12    	 6902982	       174.4 ns/op	       0 B/op	       0 allocs/op
func BenchmarkPhuslu_WithFields(b *testing.B) {
	var writer helpers.NoopWriter

	logger := log.Logger{
		Level:      log.InfoLevel,
		TimeField:  "ts",
		TimeFormat: time.RFC3339,
		Writer:     &phuslog.IOWriter{Writer: &writer},
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
