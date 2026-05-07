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

// go test -run '^$' -bench '^BenchmarkLogger_Print$' -benchmem

// BenchmarkLogger_Print-16    	32919981	        36.53 ns/op	       0 B/op	       0 allocs/op
func BenchmarkLogger_Print(b *testing.B) {
	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		&helpers.NoopWriter{},
	)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelInfo,

			WithFatalWriter: os.Stdout,
		},
	)
	require.NoError(b, errCrLogger)
	require.NotNil(b, logger)

	runtime.GC()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		logger.Print("hi", 123, "world")
	}

	cancel()
	<-chIngestionEnd
}

// BenchmarkLogger_PrintWithNoTimestamp-16    	39153303	        30.57 ns/op	       0 B/op	       0 allocs/op
func BenchmarkLogger_PrintWithNoTimestamp(b *testing.B) {
	var writer helpers.NoopWriter

	ingestor, errCrIngestor := bytearena.NewIngestor(bytearena.Size100K(), &writer)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelInfo,

			WithFatalWriter: os.Stdout,
		},
	)
	require.NoError(b, errCrLogger)
	require.NotNil(b, logger)

	runtime.GC()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		logger.PrintWithNoTimestamp("hi", 123, "world")
	}

	cancel()
	<-chIngestionEnd
}

// BenchmarkLogger_PrintRaw-16    	69086016	        17.51 ns/op	       0 B/op	       0 allocs/op
func BenchmarkLogger_PrintRaw(b *testing.B) {
	var writer helpers.NoopWriter

	ingestor, errCrIngestor := bytearena.NewIngestor(bytearena.Size100K(), &writer)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelInfo,

			WithFatalWriter: os.Stdout,
		},
	)
	require.NoError(b, errCrLogger)
	require.NotNil(b, logger)

	runtime.GC()

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		logger.PrintRaw(
			[]byte("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"), // 32 bytes
		)
	}

	cancel()
	<-chIngestionEnd
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkLogger_Printf/1.nil_timestamp-16         	12283927	        99.22 ns/op	       8 B/op	       0 allocs/op
// BenchmarkLogger_Printf/2.standard_timestamp-16    	10164178	       117.4 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_Printf/3.yyyy-month_timestamp-16  	10367665	       116.6 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_Printf/4.nano_timestamp-16        	 9150955	       130.7 ns/op	       8 B/op	       0 allocs/op
// BenchmarkLogger_Printf/5.nano_timestamp_-_json-16 	 6316914	       185.7 ns/op	      33 B/op	       1 allocs/op
func BenchmarkLogger_Printf(b *testing.B) {
	tests := []struct {
		timestampFunc timestamp.Timestamp
		description   string
		withJSON      bool
		withCaller    bool
	}{
		{
			description: "1.nil timestamp",
		},
		{
			description:   "2.standard timestamp",
			timestampFunc: timestamp.TimestampStandard,
		},
		{
			description:   "3.yyyy-month timestamp",
			timestampFunc: timestamp.TimestampYYYYMonth,
		},
		{
			description:   "4.nano timestamp",
			timestampFunc: timestamp.TimestampNano,
		},
		{
			description:   "5.nano timestamp - json",
			timestampFunc: timestamp.TimestampNano,
			withJSON:      true,
		},
	}

	for _, tc := range tests {
		b.Run(
			tc.description,
			func(b *testing.B) {
				ingestor, errCrIngestor := bytearena.NewIngestor(
					bytearena.Size100K(),
					&helpers.NoopWriter{},
				)
				require.NoError(b, errCrIngestor)
				require.NotNil(b, ingestor)

				ctx, cancel := context.WithCancel(context.Background())
				chIngestionEnd := ingestor.StartIngestion(ctx)

				logger, errCrLogger := NewLogger(
					&ParamsNewLogger{
						Ingestor:    ingestor,
						LoggerLevel: LevelInfo,

						WithFatalWriter: os.Stdout,
						WithTimestamp:   tc.timestampFunc,
						WithJSON:        tc.withJSON,
						WithCaller:      tc.withCaller,
					},
				)
				require.NoError(b, errCrLogger)
				require.NotNil(b, logger)

				runtime.GC()

				b.ReportAllocs()
				b.ResetTimer()

				for i := 0; b.Loop(); i++ {
					logger.Printf(
						`{"level":"info","msg":"user login","user_id":%d}`,
						i,
					)
				}

				cancel()
				<-chIngestionEnd
			},
		)
	}
}
