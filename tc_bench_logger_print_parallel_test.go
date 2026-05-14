package log

import (
	"context"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/bytearena/helpers"
	"github.com/tudorhulban/log/timestamp"
)

// cpu: AMD Ryzen 5 5600U with Radeon Graphics
// BenchmarkLogger_Parallel_PrintRaw-12    	41441428	        28.67 ns/op	       0 B/op	       0 allocs/op

func BenchmarkLogger_Parallel_PrintRaw(b *testing.B) {
	prev := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(prev)

	writer := helpers.CountWriterNoBuffer{}

	ingestor, errCrIngestor := bytearena.NewIngestor(bytearena.Size100K(), &writer)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	time.Sleep(10 * time.Millisecond) // warmup

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

	b.SetParallelism(16)
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				logger.PrintRaw([]byte(_Payload))
			}
		},
	)

	cancel()
	<-chIngestionEnd

	require.NotZero(b,
		writer.TotalBytesWritten.Load(),
	)
}

// go test -run '^$' -bench '^BenchmarkLogger_Parallel_Printf$' -benchmem -race

// cpu: AMD Ryzen 5 5600U with Radeon Graphics
// BenchmarkLogger_Parallel_Printf/1.nil_timestamp-12         	20619877	        59.20 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Printf/2.standard_timestamp-12    	20188375	        59.13 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Printf/3.yyyy-month_timestamp-12  	20417576	        59.18 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Printf/4.nano_timestamp-12        	20330949	        59.27 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Printf/5.nano_timestamp_-_json-12 	20169340	        59.76 ns/op	       8 B/op	       1 allocs/op

func BenchmarkLogger_Parallel_Printf(b *testing.B) {
	runtime.GOMAXPROCS(1)

	tests := []struct {
		timestampFunc timestamp.Timestamp
		description   string
		withJSON      bool
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

				time.Sleep(10 * time.Millisecond) // warmup

				logger, errCrLogger := NewLogger(
					&ParamsNewLogger{
						Ingestor:    ingestor,
						LoggerLevel: LevelInfo,

						WithFatalWriter: os.Stdout,
						WithTimestamp:   tc.timestampFunc,
						WithJSON:        tc.withJSON,
					},
				)
				require.NoError(b, errCrLogger)
				require.NotNil(b, logger)

				runtime.GC()

				b.SetParallelism(16)
				b.ReportAllocs()
				b.ResetTimer()

				b.RunParallel(
					func(pb *testing.PB) {
						i := 0

						for pb.Next() {
							logger.Printf(
								`{"level":"info","msg":"user login","user_id":%d}`,
								i,
							)

							i++
						}
					},
				)

				cancel()
				<-chIngestionEnd
			},
		)
	}
}

// cpu: AMD Ryzen 5 5600U with Radeon Graphics
// BenchmarkLogger_Parallel_Printw/1._nil_timestamp-12         	20561943	        59.40 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Printw/2._standard_timestamp-12    	20253273	        58.38 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Printw/3._yyyy-month_timestamp-12  	20634520	        58.31 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Printw/4._nano_timestamp-12        	20525589	        58.75 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Printw/5._nano_timestamp_-_json-12 	20244883	        59.15 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_Parallel_Printw/6._nano_-_json,_caller-12   	17569348	        68.56 ns/op	       8 B/op	       1 allocs/op

func BenchmarkLogger_Parallel_Printw(b *testing.B) {
	runtime.GOMAXPROCS(1) // TODO: add multiple values

	tests := []struct {
		timestampFunc timestamp.Timestamp
		description   string

		withJSON   bool
		withCaller bool
	}{
		{
			description: "1. nil timestamp",
		},
		{
			description:   "2. standard timestamp",
			timestampFunc: timestamp.TimestampStandard,
		},
		{
			description:   "3. yyyy-month timestamp",
			timestampFunc: timestamp.TimestampYYYYMonth,
		},
		{
			description:   "4. nano timestamp",
			timestampFunc: timestamp.TimestampNano,
		},
		{
			description:   "5. nano timestamp - json",
			timestampFunc: timestamp.TimestampNano,
			withJSON:      true,
		},
		{
			description:   "6. nano - json, caller",
			timestampFunc: timestamp.TimestampNano,
			withJSON:      true,
			withCaller:    true,
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

				time.Sleep(10 * time.Millisecond) // warmup

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

				b.SetParallelism(16)
				b.ReportAllocs()
				b.ResetTimer()

				b.RunParallel(
					func(pb *testing.PB) {
						i := 0

						for pb.Next() {
							logger.Printf(
								`{"level":"info","msg":"user login","user_id":%d}`,
								i,
							)
							i++
						}
					},
				)

				cancel()
				<-chIngestionEnd
			},
		)
	}
}
