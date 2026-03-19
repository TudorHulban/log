package log

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/log/bytearena"
	"github.com/tudorhulban/log/timestamp"
)

// BenchmarkArenaNilTimestamp-16    	10375689	       110.8 ns/op	      92 B/op	       0 allocs/op
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

// BenchmarkLogger_Print-16    	35267358	        33.05 ns/op	      24 B/op	       0 allocs/op
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

// BenchmarkLogger_PrintWithNoTimestamp-16    	34212000	        34.30 ns/op	      25 B/op	       0 allocs/op
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

// BenchmarkNanoTimestamp-16    	 5805175	       204.6 ns/op	     182 B/op	       1 allocs/op
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

// BenchmarkLogger_NanoTimestamp_JSON-16    	 3617188	       335.4 ns/op	     305 B/op	       3 allocs/op
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

// BenchmarkStandardTimestamp-16    	 6339567	       184.8 ns/op	     169 B/op	       1 allocs/op
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

// BenchmarkYYYYTimestamp-16    	 6268123	       185.1 ns/op	     171 B/op	       1 allocs/op
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
