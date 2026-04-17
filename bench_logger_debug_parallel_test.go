package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/log/helpers"
	"github.com/tudorhulban/log/timestamp"
)

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkLogger_Parallel_DebugFast/1.nil_timestamp-16         	20935722	        58.98 ns/op	       8 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_DebugFast/2.standard_timestamp-16    	20320483	        59.29 ns/op	       8 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_DebugFast/3.yyyy-month_timestamp-16  	19997485	        59.22 ns/op	       8 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_DebugFast/4.nano_timestamp-16        	20287327	        59.34 ns/op	       8 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_DebugFast/5.nano_timestamp_-_json-16 	20216437	        59.43 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_Parallel_DebugFast/6.nano_-_json,_caller-16   	19846812	        60.27 ns/op	       9 B/op	       1 allocs/op
func BenchmarkLogger_Parallel_DebugFast(b *testing.B) {
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

				b.SetParallelism(1)
				b.ReportAllocs()
				b.ResetTimer()

				b.RunParallel(
					func(pb *testing.PB) {
						i := 0

						for pb.Next() {
							logger.DebugFast(
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
