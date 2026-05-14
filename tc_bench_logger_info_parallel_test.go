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
// BenchmarkLogger_Parallel_Info/1.nil_timestamp-16         	21114748	        59.00 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Info/2.standard_timestamp-16    	20245314	        59.36 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Info/3.yyyy-month_timestamp-16  	19998660	        59.42 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Info/4.nano_timestamp-16        	20076030	        59.31 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Info/5.nano_timestamp_-_json-16 	20193655	        59.49 ns/op	       0 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Info/6.nano_-_json,_caller-16   	18852224	        63.66 ns/op	       0 B/op	       0 allocs/op

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
// BenchmarkLogger_Parallel_Infof/1.nil_timestamp-16         	20651301	        60.43 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Infof/2.standard_timestamp-16    	19452838	        62.51 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Infof/3.yyyy-month_timestamp-16  	18907471	        61.86 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Infof/4.nano_timestamp-16        	19386944	        62.66 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Infof/5.nano_timestamp_-_json-16 	19244624	        61.49 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_Infof/6.nano_-_json,_caller-16   	18153415	        67.98 ns/op	       7 B/op	       0 allocs/op

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
// BenchmarkLogger_Parallel_Infow/1.nil_timestamp-16         	19160923	        63.98 ns/op	      32 B/op	       1 allocs/op
// BenchmarkLogger_Parallel_Infow/2.standard_timestamp-16    	15703446	        67.51 ns/op	      32 B/op	       1 allocs/op
// BenchmarkLogger_Parallel_Infow/3.yyyy-month_timestamp-16  	16751235	        75.11 ns/op	      32 B/op	       1 allocs/op
// BenchmarkLogger_Parallel_Infow/4.nano_timestamp-16        	16930497	        76.23 ns/op	      32 B/op	       1 allocs/op
// BenchmarkLogger_Parallel_Infow/5.nano_timestamp_-_json-16 	17022232	        73.78 ns/op	      32 B/op	       1 allocs/op
// BenchmarkLogger_Parallel_Infow/6.nano_-_json,_caller-16   	15994668	        76.77 ns/op	      32 B/op	       1 allocs/op

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
