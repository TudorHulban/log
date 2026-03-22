package log

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/log/bytearena"
	"github.com/tudorhulban/log/timestamp"
)

// BenchmarkLogger_NilTimestamp-16    	35434408	        33.57 ns/op	     106 B/op	       0 allocs/op
func BenchmarkLogger_NilTimestamp(b *testing.B) {
	var writer bytes.Buffer

	ingestor := bytearena.NewIngestor(bytearena.Size100K, &writer)
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

// BenchmarkLogger_NilTimestamp_Safe-16    	 9477679	       123.9 ns/op	     164 B/op	       1 allocs/op
func BenchmarkLogger_NilTimestamp_Safe(b *testing.B) {
	var writer bytes.Buffer

	ingestor := bytearena.NewIngestor(bytearena.Size100K, &writer)
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
// BenchmarkLogger_Print-16    	66342212	        17.05 ns/op	      52 B/op	       0 allocs/op
func BenchmarkLogger_Print(b *testing.B) {
	var writer bytes.Buffer

	ingestor := bytearena.NewIngestor(bytearena.Size100K, &writer)
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

// BenchmarkLogger_PrintWithNoTimestamp-16    	69018649	        16.96 ns/op	      50 B/op	       0 allocs/op
func BenchmarkLogger_PrintWithNoTimestamp(b *testing.B) {
	var writer bytes.Buffer

	ingestor := bytearena.NewIngestor(bytearena.Size100K, &writer)
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

// BenchmarkLogger_PrintRaw-16    	85781427	        13.16 ns/op	      20 B/op	       0 allocs/op
func BenchmarkLogger_PrintRaw(b *testing.B) {
	var writer bytes.Buffer

	ingestor := bytearena.NewIngestor(bytearena.Size100K, &writer)
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

// BenchmarkLogger_NanoTimestamp-16    	29470323	        39.59 ns/op	      67 B/op	       0 allocs/op
func BenchmarkLogger_NanoTimestamp(b *testing.B) {
	var writer bytes.Buffer

	ingestor := bytearena.NewIngestor(bytearena.Size100K, &writer)
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

// BenchmarkLogger_NanoTimestamp_JSON-16    	26897366	        44.00 ns/op	      72 B/op	       0 allocs/op
func BenchmarkLogger_NanoTimestamp_JSON(b *testing.B) {
	var writer bytes.Buffer

	ingestor := bytearena.NewIngestor(bytearena.Size100K, &writer)
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

// BenchmarkLogger_StandardTimestamp-16    	27820491	        48.12 ns/op	      77 B/op	       1 allocs/op
func BenchmarkLogger_StandardTimestamp(b *testing.B) {
	var writer bytes.Buffer

	ingestor := bytearena.NewIngestor(bytearena.Size100K, &writer)
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

// BenchmarkLogger_YYYYTimestamp-16    	26702970	        39.50 ns/op	      73 B/op	       1 allocs/op
func BenchmarkLogger_YYYYTimestamp(b *testing.B) {
	var writer bytes.Buffer

	ingestor := bytearena.NewIngestor(bytearena.Size100K, &writer)
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

// BenchmarkLogger_Printf_TimestampNano-16    	29792701	        42.18 ns/op	      66 B/op	       0 allocs/op
func BenchmarkLogger_Printf_TimestampNano(b *testing.B) {
	var writer bytes.Buffer

	ingestor := bytearena.NewIngestor(bytearena.Size100K, &writer)
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

// BenchmarkLogger_Print_TimestampNano-16    	58096893	        20.44 ns/op	      60 B/op	       0 allocs/op
func BenchmarkLogger_Print_TimestampNano(b *testing.B) {
	var writer bytes.Buffer

	ingestor := bytearena.NewIngestor(bytearena.Size100K, &writer)
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
