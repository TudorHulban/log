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
// BenchmarkLogger_Infof/1.nil_timestamp-16         	23420386	        53.70 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_Infof/2.standard_timestamp-16    	22038463	        54.88 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_Infof/3.yyyy-month_timestamp-16  	22626078	        54.07 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_Infof/4.nano_timestamp-16        	 1621762	       739.2 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_Infof/5.json_-_rfc3339-16        	17267485	        70.38 ns/op	      10 B/op	       1 allocs/op
// BenchmarkLogger_Infof/6.nano_timestamp_-_json-16 	 1813675	       657.5 ns/op	      10 B/op	       1 allocs/op
// BenchmarkLogger_Infof/7.nano_-_json_caller-16    	  562224	      1925 ns/op	     264 B/op	       3 allocs/op

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

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkLogger_InfoFast/1.nil_timestamp-16         	24311269	        50.75 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_InfoFast/2.standard_timestamp-16    	22581175	        53.27 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_InfoFast/3.yyyy-month_timestamp-16  	22399359	        53.05 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_InfoFast/4.nano_timestamp-16        	 1636036	       731.6 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_InfoFast/5.json_-_rfc3339-16        	17275332	        70.38 ns/op	      10 B/op	       1 allocs/op
// BenchmarkLogger_InfoFast/6.nano_timestamp_-_json-16 	 1811938	       663.8 ns/op	      11 B/op	       1 allocs/op
// BenchmarkLogger_InfoFast/7.nano_-_json,_caller-16   	  505116	      2274 ns/op	     265 B/op	       3 allocs/op

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
