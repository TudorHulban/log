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
// BenchmarkLogger_Debugf/1.nil_timestamp-16         	19058520	        65.01 ns/op	      72 B/op	       1 allocs/op
// BenchmarkLogger_Debugf/2.standard_timestamp-16    	10527141	       115.8 ns/op	      72 B/op	       2 allocs/op
// BenchmarkLogger_Debugf/3.yyyy-month_timestamp-16  	10778720	       112.3 ns/op	      72 B/op	       2 allocs/op
// BenchmarkLogger_Debugf/4.nano_timestamp-16        	 9408266	       128.9 ns/op	      72 B/op	       1 allocs/op
// BenchmarkLogger_Debugf/5.nano_timestamp_-_json-16 	 7244959	       167.2 ns/op	     112 B/op	       2 allocs/op
// BenchmarkLogger_Debugf/6.nano_-_json,_caller-16   	 2203927	       545.8 ns/op	     552 B/op	       5 allocs/op
func BenchmarkLogger_Debugf(b *testing.B) {
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

	for _, tcase := range tests {
		b.Run(
			tcase.description,
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
						LoggerLevel: LevelDebug,

						WithFatalWriter: os.Stdout,
						WithTimestamp:   tcase.timestampFunc,
						WithJSON:        tcase.withJSON,
						WithCaller:      tcase.withCaller,
					},
				)
				require.NoError(b, errCrLogger)
				require.NotNil(b, logger)

				runtime.GC()

				b.ReportAllocs()
				b.ResetTimer()

				for i := 0; b.Loop(); i++ {
					logger.Debugf(
						"%d",
						i,
					)
				}

				cancel()
				<-chIngestionEnd
			},
		)
	}
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkLogger_DebugFast/1.nil_timestamp-16         	10002454	       118.3 ns/op	      50 B/op	       2 allocs/op
// BenchmarkLogger_DebugFast/2.standard_timestamp-16    	 8879348	       132.0 ns/op	      50 B/op	       2 allocs/op
// BenchmarkLogger_DebugFast/3.yyyy-month_timestamp-16  	 8820680	       131.8 ns/op	      50 B/op	       2 allocs/op
// BenchmarkLogger_DebugFast/4.nano_timestamp-16        	 8543960	       138.1 ns/op	      50 B/op	       2 allocs/op
// BenchmarkLogger_DebugFast/5.nano_timestamp_-_json-16 	 7802482	       153.5 ns/op	      55 B/op	       2 allocs/op
// BenchmarkLogger_DebugFast/6.nano_-_json,_caller-16   	 4751898	       252.7 ns/op	     121 B/op	       3 allocs/op
func BenchmarkLogger_DebugFast(b *testing.B) {
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

	for _, tcase := range tests {
		b.Run(
			tcase.description,
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
						LoggerLevel: LevelDebug,

						WithFatalWriter: os.Stdout,
						WithTimestamp:   tcase.timestampFunc,
						WithJSON:        tcase.withJSON,
						WithCaller:      tcase.withCaller,
					},
				)
				require.NoError(b, errCrLogger)
				require.NotNil(b, logger)

				runtime.GC()

				b.ReportAllocs()
				b.ResetTimer()

				for i := 0; b.Loop(); i++ {
					logger.DebugFast(
						"%d",
						i,
					)
				}

				cancel()
				<-chIngestionEnd
			},
		)
	}
}
