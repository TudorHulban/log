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
// BenchmarkLogger_Infof/1.nil_timestamp-16         	17719626	        67.83 ns/op	      72 B/op	       1 allocs/op
// BenchmarkLogger_Infof/2.standard_timestamp-16    	10622220	       114.7 ns/op	      72 B/op	       2 allocs/op
// BenchmarkLogger_Infof/3.yyyy-month_timestamp-16  	10527896	       113.8 ns/op	      72 B/op	       2 allocs/op
// BenchmarkLogger_Infof/4.nano_timestamp-16        	 9061526	       133.8 ns/op	      72 B/op	       1 allocs/op
// BenchmarkLogger_Infof/5.json_-_rfc3339-16        	 6598698	       183.0 ns/op	     112 B/op	       2 allocs/op
// BenchmarkLogger_Infof/6.nano_timestamp_-_json-16 	 6775639	       179.0 ns/op	     112 B/op	       2 allocs/op
// BenchmarkLogger_Infof/7.nano_-_json,_caller-16   	 2183959	       551.3 ns/op	     552 B/op	       5 allocs/op
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
// BenchmarkLogger_InfoFast/1.nil_timestamp-16         	12252991	        97.87 ns/op	      46 B/op	       2 allocs/op
// BenchmarkLogger_InfoFast/2.standard_timestamp-16    	11375336	       104.6 ns/op	      46 B/op	       2 allocs/op
// BenchmarkLogger_InfoFast/3.yyyy-month_timestamp-16  	11259042	       105.0 ns/op	      46 B/op	       2 allocs/op
// BenchmarkLogger_InfoFast/4.nano_timestamp-16        	10979743	       108.9 ns/op	      46 B/op	       2 allocs/op
// BenchmarkLogger_InfoFast/5.json_-_rfc3339-16        	 9975960	       119.6 ns/op	      48 B/op	       2 allocs/op
// BenchmarkLogger_InfoFast/6.nano_timestamp_-_json-16 	10051143	       118.9 ns/op	      48 B/op	       2 allocs/op
// BenchmarkLogger_InfoFast/7.nano_-_json,_caller-16   	 6654192	       179.3 ns/op	      88 B/op	       2 allocs/op
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

						EstimatedMessageSizeInfo: MessageMediumSize,

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
