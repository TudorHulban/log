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

// BenchmarkLogger_Parallel_PrintRaw-16    	43104222	        27.22 ns/op	       0 B/op	       0 allocs/op
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
				logger.PrintRaw([]byte("xxxxxxxxxxxxxxxxxxx"))
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

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkLogger_Parallel_Printf/1.nil_timestamp-16         	20686093	        60.59 ns/op	       8 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Printf/2.standard_timestamp-16    	20122062	        60.64 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Printf/3.yyyy-month_timestamp-16  	20306156	        60.52 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Printf/4.nano_timestamp-16        	20453836	        59.92 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Printf/5.nano_timestamp_-_json-16 	20047094	        59.86 ns/op	       8 B/op	       0 allocs/op
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

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkLogger_Parallel_PrintfFast/1._nil_timestamp-16         	21116536	        59.38 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_PrintfFast/2._standard_timestamp-16    	20543920	        59.36 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_PrintfFast/3._yyyy-month_timestamp-16  	20435289	        59.44 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_PrintfFast/4._nano_timestamp-16        	20611592	        60.00 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_PrintfFast/5._nano_timestamp_-_json-16 	20369793	        60.52 ns/op	       8 B/op	       0 allocs/op
func BenchmarkLogger_Parallel_PrintfFast(b *testing.B) {
	runtime.GOMAXPROCS(1)

	tests := []struct {
		timestampFunc timestamp.Timestamp
		description   string
		withJSON      bool
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
							logger.PrintfFast(
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
