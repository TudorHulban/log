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

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkLogger_Parallel_Info/1.nil_timestamp-16         	21219388	        58.37 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Info/2.standard_timestamp-16    	20901224	        58.61 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Info/3.yyyy-month_timestamp-16  	20817206	        58.57 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Info/4.nano_timestamp-16        	20793866	        58.50 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Info/5.nano_timestamp_-_json-16 	20825125	        58.69 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Info/6.nano_-_json,_caller-16   	10581820	       108.2 ns/op	     248 B/op	       2 allocs/op

func BenchmarkLogger_Parallel_Info(b *testing.B) {
	runtime.GOMAXPROCS(1)

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
		{
			description:   "6.nano - json, caller",
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
					&helpers.CountWriterNoBuffer{},
				)
				require.NoError(b, errCrIngestor)
				require.NotNil(b, ingestor)

				ctx, cancel := context.WithCancel(context.Background())
				chIngestionEnd := ingestor.StartIngestion(ctx)

				time.Sleep(10 * time.Millisecond) // warmup

				logger, errCrLogger := NewLogger(
					&ParamsNewLogger{
						Ingestor:    ingestor,
						LoggerLevel: LevelDebug,

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
						for pb.Next() {
							logger.Info(_Payload)
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
// BenchmarkLogger_Parallel_Infof/1.nil_timestamp-16         	21049404	        59.19 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Infof/2.standard_timestamp-16    	20375635	        60.84 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Infof/3.yyyy-month_timestamp-16  	20215396	        60.64 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Infof/4.nano_timestamp-16        	20296920	        59.71 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Infof/5.nano_timestamp_-_json-16 	20112969	        59.68 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Infof/6.nano_-_json,_caller-16   	 9600229	       122.5 ns/op	     256 B/op	       3 allocs/op

func BenchmarkLogger_Parallel_Infof(b *testing.B) {
	runtime.GOMAXPROCS(1)

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
		{
			description:   "6.nano - json, caller",
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
					&helpers.CountWriterNoBuffer{},
				)
				require.NoError(b, errCrIngestor)
				require.NotNil(b, ingestor)

				ctx, cancel := context.WithCancel(context.Background())
				chIngestionEnd := ingestor.StartIngestion(ctx)

				time.Sleep(10 * time.Millisecond) // warmup

				logger, errCrLogger := NewLogger(
					&ParamsNewLogger{
						Ingestor:    ingestor,
						LoggerLevel: LevelDebug,

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
							logger.Infof(
								"%d",
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
// BenchmarkLogger_Parallel_Infow/1.nil_timestamp-16         	19439230	        62.13 ns/op	      32 B/op	       1 allocs/op
// BenchmarkLogger_Parallel_Infow/2.standard_timestamp-16    	19434534	        62.66 ns/op	      32 B/op	       1 allocs/op
// BenchmarkLogger_Parallel_Infow/3.yyyy-month_timestamp-16  	19168959	        62.48 ns/op	      32 B/op	       1 allocs/op
// BenchmarkLogger_Parallel_Infow/4.nano_timestamp-16        	19246233	        62.74 ns/op	      32 B/op	       1 allocs/op
// BenchmarkLogger_Parallel_Infow/5.nano_timestamp_-_json-16 	18490947	        63.57 ns/op	      32 B/op	       1 allocs/op
// BenchmarkLogger_Parallel_Infow/6.nano_-_json,_caller-16   	 8349188	       137.0 ns/op	     280 B/op	       3 allocs/op

func BenchmarkLogger_Parallel_Infow(b *testing.B) {
	runtime.GOMAXPROCS(1)

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
		{
			description:   "6.nano - json, caller",
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
					&helpers.CountWriterNoBuffer{},
				)
				require.NoError(b, errCrIngestor)
				require.NotNil(b, ingestor)

				ctx, cancel := context.WithCancel(context.Background())
				chIngestionEnd := ingestor.StartIngestion(ctx)

				time.Sleep(10 * time.Millisecond) // warmup

				logger, errCrLogger := NewLogger(
					&ParamsNewLogger{
						Ingestor:    ingestor,
						LoggerLevel: LevelDebug,

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
						for pb.Next() {
							logger.Infow("key", _Payload)
						}
					},
				)

				cancel()
				<-chIngestionEnd
			},
		)
	}
}
