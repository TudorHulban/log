package log

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/tudorhulban/log/bytearena"
	"github.com/tudorhulban/log/helpers"
	"github.com/tudorhulban/log/timestamp"
)

// BenchmarkNilTimestamp-16    	14189306	        87.69 ns/op	       8 B/op	       0 allocs/op
func BenchmarkNilTimestamp(b *testing.B) {
	b.ReportAllocs()

	sink := helpers.CountWriter{}

	logger := NewLogger(
		&ParamsNewLogger{
			LoggerWriter:  &sink,
			LoggerLevel:   LevelINFO,
			WithTimestamp: timestamp.TimestampNil,
		},
	)

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

	logger := NewLogger(
		&ParamsNewLogger{
			LoggerWriter:  writer,
			LoggerLevel:   LevelINFO,
			WithTimestamp: timestamp.TimestampNil,
		},
	)

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

func BenchmarkLogger_Print(b *testing.B) {
	var sink bytes.Buffer

	writer := bytearena.NewIngestor(bytearena.Size100K, &sink)
	ctx, cancel := context.WithCancel(context.Background())
	chEnd := writer.StartIngestion(ctx)

	defer func() {
		cancel()
		<-chEnd
	}()

	logger := NewLogger(
		&ParamsNewLogger{
			LoggerWriter:  writer,
			LoggerLevel:   LevelINFO,
			WithTimestamp: timestamp.TimestampNil,
		},
	)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		logger.Print("hi", 123, "world")
	}
}

// BenchmarkLogger-16    	 8522778	       145.4 ns/op	       8 B/op	       0 allocs/op
func BenchmarkNanoTimestamp(b *testing.B) {
	b.ReportAllocs()

	sink := helpers.CountWriter{}

	logger := NewLogger(
		&ParamsNewLogger{
			LoggerWriter:  &sink,
			LoggerLevel:   LevelINFO,
			WithTimestamp: timestamp.TimestampNano,
		},
	)

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

	logger := NewLogger(
		&ParamsNewLogger{
			LoggerWriter:  &sink,
			LoggerLevel:   LevelINFO,
			WithTimestamp: timestamp.TimestampStandard,
		},
	)

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

	logger := NewLogger(
		&ParamsNewLogger{
			LoggerWriter:  &sink,
			LoggerLevel:   LevelINFO,
			WithTimestamp: timestamp.TimestampYYYYMonth,
		},
	)

	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		logger.Printf(
			`{"level":"info","msg":"user login","user_id":%d}`,
			i,
		)
	}

	_ = sink.NumberWrites.Load() // force sink to stay live
}
