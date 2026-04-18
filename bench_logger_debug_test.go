package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/bytearena/helpers"
	"github.com/tudorhulban/log/timestamp"
)

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkLogger_Debugf/1.nil_timestamp-16         	18687692	        69.03 ns/op	      72 B/op	       2 allocs/op
// BenchmarkLogger_Debugf/2.standard_timestamp-16    	10765196	       113.9 ns/op	      72 B/op	       2 allocs/op
// BenchmarkLogger_Debugf/3.yyyy-month_timestamp-16  	10722692	       113.1 ns/op	      72 B/op	       2 allocs/op
// BenchmarkLogger_Debugf/4.nano_timestamp-16        	 9199351	       131.4 ns/op	      72 B/op	       2 allocs/op
// BenchmarkLogger_Debugf/5.nano_timestamp_-_json-16 	 7004761	       173.8 ns/op	     112 B/op	       3 allocs/op
// BenchmarkLogger_Debugf/6.nano_-_json,_caller-16   	 2217819	       544.0 ns/op	     552 B/op	       6 allocs/op
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

				logger, errCrLogger := NewLogger(
					&ParamsNewLogger{
						Ingestor:    ingestor,
						LoggerLevel: LevelDEBUG,

						WithTimestamp: tcase.timestampFunc,
						WithJSON:      tcase.withJSON,
						WithCaller:    tcase.withCaller,
					},
				)
				require.NoError(b, errCrLogger)
				require.NotNil(b, logger)

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
// BenchmarkLogger_DebugFast/1.nil_timestamp-16         	22843332	        53.04 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_DebugFast/2.standard_timestamp-16    	19039914	        64.34 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_DebugFast/3.yyyy-month_timestamp-16  	18960468	        64.87 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_DebugFast/4.nano_timestamp-16        	17616594	        68.71 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_DebugFast/5.nano_timestamp_-_json-16 	14479726	        83.08 ns/op	      14 B/op	       1 allocs/op
// BenchmarkLogger_DebugFast/6.nano_-_json,_caller-16   	 7473280	       162.5 ns/op	      75 B/op	       1 allocs/op
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
						LoggerLevel: LevelDEBUG,

						WithTimestamp: tcase.timestampFunc,
						WithJSON:      tcase.withJSON,
						WithCaller:    tcase.withCaller,
					},
				)
				require.NoError(b, errCrLogger)
				require.NotNil(b, logger)

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
