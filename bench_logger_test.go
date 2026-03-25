package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/log/helpers"
	"github.com/tudorhulban/log/timestamp"
)

// BenchmarkLogger_NilTimestamp-16    	36902835	        32.48 ns/op	       8 B/op	       1 allocs/op
func BenchmarkLogger_NilTimestamp(b *testing.B) {
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
			Ingestor:    ingestor,
			LoggerLevel: Level(LevelINFO),
		},
	)
	require.NoError(b, errCrLogger)
	require.NotNil(b, logger)

	for i := 0; b.Loop(); i++ {
		logger.Printf(
			`{"level":"info","msg":"user login","user_id":%d}`,
			i,
		)
	}
}

// BenchmarkLogger_NilTimestamp_Safe-16    	 9587048	       125.6 ns/op	      72 B/op	       2 allocs/op
func BenchmarkLogger_NilTimestamp_Safe(b *testing.B) {
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
			Ingestor:    ingestor,
			LoggerLevel: Level(LevelINFO),
		},
	)
	require.NoError(b, errCrLogger)
	require.NotNil(b, logger)

	for i := 0; b.Loop(); i++ {
		logger.PrintfSafe(
			`{"level":"info","msg":"user login","user_id":%d}`,
			i,
		)
	}
}

// go test -run '^$' -bench '^BenchmarkLogger_Print$' -benchmem
// go test -run '^$' -bench '^BenchmarkLogger_Print$' -benchmem -race
// BenchmarkLogger_Print-16    	73801082	        16.25 ns/op	       0 B/op	       0 allocs/op
func BenchmarkLogger_Print(b *testing.B) {
	var writer helpers.NoopWriter

	ingestor, errCrIngestor := bytearena.NewIngestor(bytearena.Size100K, &writer)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chEnd := ingestor.StartIngestion(ctx)

	defer func() {
		cancel()
		<-chEnd
	}()

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: Level(LevelINFO),
		},
	)
	require.NoError(b, errCrLogger)
	require.NotNil(b, logger)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		logger.Print("hi", 123, "world")
	}
}

// BenchmarkLogger_PrintWithNoTimestamp-16    	74154141	        16.26 ns/op	       0 B/op	       0 allocs/op
func BenchmarkLogger_PrintWithNoTimestamp(b *testing.B) {
	var writer helpers.NoopWriter

	ingestor, errCrIngestor := bytearena.NewIngestor(bytearena.Size100K, &writer)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chEnd := ingestor.StartIngestion(ctx)

	defer func() {
		cancel()
		<-chEnd
	}()

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: Level(LevelINFO),
		},
	)
	require.NoError(b, errCrLogger)
	require.NotNil(b, logger)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		logger.PrintWithNoTimestamp("hi", 123, "world")
	}
}

// BenchmarkLogger_PrintRaw-16    	100000000	        10.87 ns/op	       0 B/op	       0 allocs/op
func BenchmarkLogger_PrintRaw(b *testing.B) {
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
			Ingestor:    ingestor,
			LoggerLevel: Level(LevelINFO),
		},
	)
	require.NoError(b, errCrLogger)
	require.NotNil(b, logger)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		logger.PrintRaw(
			[]byte("xxxxxxxxxxxxxxxxxxx"),
		)
	}
}

// BenchmarkLogger_NanoTimestamp-16    	28327180	        43.53 ns/op	       8 B/op	       1 allocs/op
func BenchmarkLogger_NanoTimestamp(b *testing.B) {
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
			LoggerLevel:   Level(LevelINFO),
			WithTimestamp: timestamp.TimestampNano,
		},
	)
	require.NoError(b, errCrLogger)
	require.NotNil(b, logger)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		logger.Printf(
			`{"level":"info","msg":"user login","user_id":%d}`,
			i,
		)
	}
}

// BenchmarkLogger_NanoTimestamp_JSON-16    	24533622	        47.86 ns/op	       8 B/op	       0 allocs/op
func BenchmarkLogger_NanoTimestamp_JSON(b *testing.B) {
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
			LoggerLevel:   Level(LevelINFO),
			WithTimestamp: timestamp.TimestampNano,
			WithJSON:      true,
		},
	)
	require.NoError(b, errCrLogger)
	require.NotNil(b, logger)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		logger.Printf(
			`{"level":"info","msg":"user login","user_id":%d}`,
			i,
		)
	}
}

// BenchmarkLogger_StandardTimestamp-16    	30710284	        39.61 ns/op	       8 B/op	       1 allocs/op
func BenchmarkLogger_StandardTimestamp(b *testing.B) {
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
			LoggerLevel:   Level(LevelINFO),
			WithTimestamp: timestamp.TimestampStandard,
		},
	)
	require.NoError(b, errCrLogger)
	require.NotNil(b, logger)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		logger.Printf(
			`{"level":"info","msg":"user login","user_id":%d}`,
			i,
		)
	}
}

// BenchmarkLogger_YYYYTimestamp-16    	29826070	        39.19 ns/op	       8 B/op	       1 allocs/op
func BenchmarkLogger_YYYYTimestamp(b *testing.B) {
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
			LoggerLevel:   Level(LevelINFO),
			WithTimestamp: timestamp.TimestampYYYYMonth,
		},
	)
	require.NoError(b, errCrLogger)
	require.NotNil(b, logger)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		logger.Printf(
			`{"level":"info","msg":"user login","user_id":%d}`,
			i,
		)
	}
}

// BenchmarkLogger_Printf_TimestampNano-16    	28026976	        42.45 ns/op	       8 B/op	       1 allocs/op
func BenchmarkLogger_Printf_TimestampNano(b *testing.B) {
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
			LoggerLevel:   Level(LevelINFO),
			WithTimestamp: timestamp.TimestampNano,
		},
	)
	require.NoError(b, errCrLogger)
	require.NotNil(b, logger)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		logger.Printf(
			`{"level":"info","msg":"user login","user_id":%d}`,
			i,
		)
	}
}

// BenchmarkLogger_Print_TimestampNano-16    	56562405	        21.79 ns/op	       0 B/op	       0 allocs/op
func BenchmarkLogger_Print_TimestampNano(b *testing.B) {
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
			LoggerLevel:   Level(LevelINFO),
			WithTimestamp: timestamp.TimestampNano,
		},
	)
	require.NoError(b, errCrLogger)
	require.NotNil(b, logger)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		logger.Print("hi", 123, "world")
	}
}
