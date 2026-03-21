package log

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/log/bytearena"
	"github.com/tudorhulban/log/timestamp"
)

// BenchmarkLogger_NilTimestamp-16    	35355093	        40.22 ns/op	     106 B/op	       0 allocs/op
func BenchmarkLogger_NilTimestamp(b *testing.B) {
	var sink bytes.Buffer

	ingestor := bytearena.NewIngestor(bytearena.Size100K, &sink)
	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	defer func() {
		cancel()
		<-chIngestionEnd
	}()

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelINFO,
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

// go test -run '^$' -bench '^BenchmarkLogger_Print$' -benchmem
// go test -run '^$' -bench '^BenchmarkLogger_Print$' -benchmem -race
// BenchmarkLogger_Print-16    	66342212	        17.05 ns/op	      52 B/op	       0 allocs/op
func BenchmarkLogger_Print(b *testing.B) {
	var sink bytes.Buffer

	writer := bytearena.NewIngestor(bytearena.Size100K, &sink)
	ctx, cancel := context.WithCancel(context.Background())
	chEnd := writer.StartIngestion(ctx)

	defer func() {
		cancel()
		<-chEnd
	}()

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    writer,
			LoggerLevel: LevelINFO,
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
	var sink bytes.Buffer

	writer := bytearena.NewIngestor(bytearena.Size100K, &sink)
	ctx, cancel := context.WithCancel(context.Background())
	chEnd := writer.StartIngestion(ctx)

	defer func() {
		cancel()
		<-chEnd
	}()

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    writer,
			LoggerLevel: LevelINFO,
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

// BenchmarkLogger_PrintRaw-16    	86922866	        13.38 ns/op	      20 B/op	       0 allocs/op
func BenchmarkLogger_PrintRaw(b *testing.B) {
	var sink bytes.Buffer

	ingestor := bytearena.NewIngestor(bytearena.Size100K, &sink)
	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	defer func() {
		cancel()
		<-chIngestionEnd
	}()

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelINFO,
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

// BenchmarkLogger_NanoTimestamp-16    	29493921	        42.23 ns/op	      69 B/op	       1 allocs/op
func BenchmarkLogger_NanoTimestamp(b *testing.B) {
	var sink bytes.Buffer

	ingestor := bytearena.NewIngestor(bytearena.Size100K, &sink)
	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	defer func() {
		cancel()
		<-chIngestionEnd
	}()

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:      ingestor,
			LoggerLevel:   LevelINFO,
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

// BenchmarkLogger_NanoTimestamp_JSON-16    	22477536	        66.49 ns/op	      92 B/op	       1 allocs/op
func BenchmarkLogger_NanoTimestamp_JSON(b *testing.B) {
	var sink bytes.Buffer

	ingestor := bytearena.NewIngestor(bytearena.Size100K, &sink)
	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	defer func() {
		cancel()
		<-chIngestionEnd
	}()

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:      ingestor,
			LoggerLevel:   LevelINFO,
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

// BenchmarkLogger_StandardTimestamp-16    	29345433	        42.11 ns/op	      69 B/op	       1 allocs/op
func BenchmarkLogger_StandardTimestamp(b *testing.B) {
	var sink bytes.Buffer

	ingestor := bytearena.NewIngestor(bytearena.Size100K, &sink)
	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	defer func() {
		cancel()
		<-chIngestionEnd
	}()

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:      ingestor,
			LoggerLevel:   LevelINFO,
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

// BenchmarkLogger_YYYYTimestamp-16    	30000571	        43.96 ns/op	      68 B/op	       1 allocs/op
func BenchmarkLogger_YYYYTimestamp(b *testing.B) {
	var sink bytes.Buffer

	ingestor := bytearena.NewIngestor(bytearena.Size100K, &sink)
	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	defer func() {
		cancel()
		<-chIngestionEnd
	}()

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:      ingestor,
			LoggerLevel:   LevelINFO,
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
