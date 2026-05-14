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
// BenchmarkLogger_Infof/1.nil_timestamp-12         	24220348	        50.29 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_Infof/2.standard_timestamp-12    	23435428	        51.70 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_Infof/3.yyyy-month_timestamp-12  	23365196	        51.58 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_Infof/4.nano_timestamp-12        	15806343	        74.48 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_Infof/5.json_-_rfc3339-12        	16623351	        73.25 ns/op	      10 B/op	       1 allocs/op
// BenchmarkLogger_Infof/6.nano_timestamp_-_json-12 	12663598	        93.26 ns/op	      11 B/op	       1 allocs/op
// BenchmarkLogger_Infof/7.nano_-_json_caller-12    	 4257105	       283.7 ns/op	      13 B/op	       1 allocs/op

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
			description:   "7.nano - json caller",
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

// cpu: AMD Ryzen 5 5600U with Radeon Graphics
// BenchmarkLogger_InfoFast/1.nil_timestamp-12         	24612154	        50.32 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_InfoFast/2.standard_timestamp-12    	22802058	        51.95 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_InfoFast/3.yyyy-month_timestamp-12  	22875586	        52.04 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_InfoFast/4.nano_timestamp-12        	16112794	        74.60 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_InfoFast/5.json_-_rfc3339-12        	16694774	        73.44 ns/op	      11 B/op	       1 allocs/op
// BenchmarkLogger_InfoFast/6.nano_timestamp_-_json-12 	12238436	        94.69 ns/op	      11 B/op	       1 allocs/op
// BenchmarkLogger_InfoFast/7.nano_-_json,_caller-12   	 4207267	       282.9 ns/op	      13 B/op	       1 allocs/op

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
