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
// BenchmarkLogger_Infof/1.nil_timestamp-16         	24787635	        49.41 ns/op	       8 B/op	       0 allocs/op
// BenchmarkLogger_Infof/2.standard_timestamp-16    	18566630	        65.94 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_Infof/3.yyyy-month_timestamp-16  	18228122	        66.08 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_Infof/4.nano_timestamp-16        	16159620	        74.89 ns/op	       7 B/op	       0 allocs/op
// BenchmarkLogger_Infof/5.json_-_rfc3339-16        	13196606	        91.23 ns/op	      10 B/op	       1 allocs/op
// BenchmarkLogger_Infof/6.nano_timestamp_-_json-16 	13767166	        88.79 ns/op	      10 B/op	       1 allocs/op
// BenchmarkLogger_Infof/7.nano_-_json,_caller-16   	 4566430	       268.3 ns/op	     102 B/op	       2 allocs/op
func BenchmarkLogger_Infof(b *testing.B) {
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
			description:   "5.json - rfc3339",
			timestampFunc: timestamp.TimestampRFC3339,
			withJSON:      true,
		},
		{
			description:   "6.nano timestamp - json",
			timestampFunc: timestamp.TimestampNano,
			withJSON:      true,
		},
		{
			description:   "7.nano - json, caller",
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

				b.ReportAllocs()
				b.ResetTimer()

				for i := 0; b.Loop(); i++ {
					logger.Infof(
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
// BenchmarkLogger_InfoFast/1.nil_timestamp-16         	28443475	        41.81 ns/op	       8 B/op	       0 allocs/op
// BenchmarkLogger_InfoFast/2.standard_timestamp-16    	25182224	        46.44 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_InfoFast/3.yyyy-month_timestamp-16  	25570363	        46.68 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_InfoFast/4.nano_timestamp-16        	24437850	        48.42 ns/op	       8 B/op	       0 allocs/op
// BenchmarkLogger_InfoFast/5.json_-_rfc3339-16        	19821165	        56.32 ns/op	      10 B/op	       1 allocs/op
// BenchmarkLogger_InfoFast/6.nano_timestamp_-_json-16 	21348616	        54.92 ns/op	      10 B/op	       1 allocs/op
// BenchmarkLogger_InfoFast/7.nano_-_json,_caller-16   	14252337	        83.78 ns/op	      33 B/op	       1 allocs/op
func BenchmarkLogger_InfoFast(b *testing.B) {
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
			description:   "5.json - rfc3339",
			timestampFunc: timestamp.TimestampRFC3339,
			withJSON:      true,
		},
		{
			description:   "6.nano timestamp - json",
			timestampFunc: timestamp.TimestampNano,
			withJSON:      true,
		},
		{
			description:   "7.nano - json, caller",
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

				logger, errCrLogger := NewLogger(
					&ParamsNewLogger{
						Ingestor:    ingestor,
						LoggerLevel: LevelDebug,

						EstimatedMessageSizeInfo: MessageMediumSize,

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
					logger.InfoFast(
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
