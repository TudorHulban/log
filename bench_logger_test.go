package log

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/log/bytearena"
	"github.com/tudorhulban/log/helpers"
	"github.com/tudorhulban/log/timestamp"
)

// BenchmarkNilTimestamp-16    	14189306	        87.69 ns/op	       8 B/op	       0 allocs/op
func BenchmarkNilTimestamp(b *testing.B) {
	b.ReportAllocs()

	sink := helpers.CountWriter{}

	ingestor := bytearena.NewIngestor(
		bytearena.Size100K,
		&sink,
	)

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelINFO,
		},
	)
	require.NoError(b, errCrLogger)

	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		logger.Printf(
			`{"level":"info","msg":"user login","user_id":%d}`,
			i,
		)
	}

	_ = sink.NumberWrites.Load() // force sink to stay live
}

func BenchmarkArenaNilTimestamp(b *testing.B) {
	b.ReportAllocs()

	writer := bytearena.NewIngestor(bytearena.Size100K, io.Discard)

	ctx, cancel := context.WithCancel(context.Background())

	chIngestionEnd := writer.StartIngestion(ctx)

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    writer,
			LoggerLevel: LevelINFO,
		},
	)
	require.NoError(b, errCrLogger)

	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		logger.Printf(
			`{"level":"info","msg":"user login","user_id":%d}`,
			i,
		)
	}

	cancel()

	// Wait for consumer shutdown flush.
	<-chIngestionEnd
}

// BenchmarkLogger_Print-16    	16389417	        68.32 ns/op	     282 B/op	       1 allocs/op
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

// BenchmarkLogger_PrintWithNoTimestamp-16    	33209036	        35.24 ns/op	      26 B/op	       0 allocs/op
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

// BenchmarkLogger_PrintFast-16    	86922866	        13.38 ns/op	      20 B/op	       0 allocs/op
func BenchmarkLogger_PrintFast(b *testing.B) {
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
		logger.PrintFast(
			[]byte("xxxxxxxxxxxxxxxxxxx"),
		)
	}
}

// BenchmarkLogger-16    	 8522778	       145.4 ns/op	       8 B/op	       0 allocs/op
func BenchmarkNanoTimestamp(b *testing.B) {
	b.ReportAllocs()

	sink := helpers.CountWriter{}

	ingestor := bytearena.NewIngestor(
		bytearena.Size100K,
		&sink,
	)

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:      ingestor,
			LoggerLevel:   LevelINFO,
			WithTimestamp: timestamp.TimestampNano,
		},
	)
	require.NoError(b, errCrLogger)

	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		logger.Printf(
			`{"level":"info","msg":"user login","user_id":%d}`,
			i,
		)
	}

	_ = sink.NumberWrites.Load() // force sink to stay live
}

// BenchmarkStandardTimestamp-16    	 9234607	       130.2 ns/op	       8 B/op	       0 allocs/op
func BenchmarkStandardTimestamp(b *testing.B) {
	b.ReportAllocs()

	sink := helpers.CountWriter{}

	ingestor := bytearena.NewIngestor(
		bytearena.Size100K,
		&sink,
	)

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:      ingestor,
			LoggerLevel:   LevelINFO,
			WithTimestamp: timestamp.TimestampStandard,
		},
	)
	require.NoError(b, errCrLogger)

	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		logger.Printf(
			`{"level":"info","msg":"user login","user_id":%d}`,
			i,
		)
	}

	_ = sink.NumberWrites.Load() // force sink to stay live
}

// BenchmarkYYYYTimestamp-16    	 9197252	       131.7 ns/op	       8 B/op	       0 allocs/op
func BenchmarkYYYYTimestamp(b *testing.B) {
	b.ReportAllocs()

	sink := helpers.CountWriter{}

	ingestor := bytearena.NewIngestor(
		bytearena.Size100K,
		&sink,
	)

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:      ingestor,
			LoggerLevel:   LevelINFO,
			WithTimestamp: timestamp.TimestampYYYYMonth,
		},
	)
	require.NoError(b, errCrLogger)

	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		logger.Printf(
			`{"level":"info","msg":"user login","user_id":%d}`,
			i,
		)
	}

	_ = sink.NumberWrites.Load() // force sink to stay live
}
