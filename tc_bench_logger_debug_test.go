package log

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/bytearena/helpers"
	"github.com/tudorhulban/log/timestamp"
)

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkLogger_Debugf/1.nil_timestamp-16         	19284004	        65.89 ns/op	      72 B/op	       2 allocs/op
// BenchmarkLogger_Debugf/2.standard_timestamp-16    	10270768	       118.5 ns/op	      72 B/op	       2 allocs/op
// BenchmarkLogger_Debugf/3.yyyy-month_timestamp-16  	10580866	       115.2 ns/op	      72 B/op	       2 allocs/op
// BenchmarkLogger_Debugf/4.nano_timestamp-16        	 9146551	       133.0 ns/op	      72 B/op	       2 allocs/op
// BenchmarkLogger_Debugf/5.nano_timestamp_-_json-16 	 6894847	       179.1 ns/op	     112 B/op	       3 allocs/op
// BenchmarkLogger_Debugf/6.nano_-_json,_caller-16   	 2191314	       551.9 ns/op	     552 B/op	       6 allocs/op
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
						LoggerLevel: LevelDebug,

						WithFatalWriter: os.Stdout,
						WithTimestamp:   tcase.timestampFunc,
						WithJSON:        tcase.withJSON,
						WithCaller:      tcase.withCaller,
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
// BenchmarkLogger_DebugFast/1.nil_timestamp-16         	 8259055	       142.5 ns/op	      52 B/op	       2 allocs/op
// BenchmarkLogger_DebugFast/2.standard_timestamp-16    	 7686872	       157.4 ns/op	      52 B/op	       2 allocs/op
// BenchmarkLogger_DebugFast/3.yyyy-month_timestamp-16  	 7665985	       154.4 ns/op	      52 B/op	       2 allocs/op
// BenchmarkLogger_DebugFast/4.nano_timestamp-16        	 7306746	       163.1 ns/op	      52 B/op	       2 allocs/op
// BenchmarkLogger_DebugFast/5.nano_timestamp_-_json-16 	 6636604	       183.7 ns/op	      57 B/op	       2 allocs/op
// BenchmarkLogger_DebugFast/6.nano_-_json,_caller-16   	 4004169	       298.9 ns/op	     131 B/op	       3 allocs/op
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
