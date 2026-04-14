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
// BenchmarkLogger_Debugf/1._standard_timestamp-12         	 8762857	       137.6 ns/op	      80 B/op	       3 allocs/op
// BenchmarkLogger_Debugf/2._yyyy-month_timestamp-12       	 8349907	       140.3 ns/op	      80 B/op	       3 allocs/op
// BenchmarkLogger_Debugf/3._nano_timestamp-12             	 7462010	       160.1 ns/op	      80 B/op	       3 allocs/op
// BenchmarkLogger_Debugf/4._nano_timestamp_-_json-12      	 5490399	       220.6 ns/op	     120 B/op	       5 allocs/op
// BenchmarkLogger_Debugf/5._nano_timestamp_-_json,_caller-12         	 2058727	       584.8 ns/op	     496 B/op	       8 allocs/op
// BenchmarkLogger_Debugf/6._nil_timestamp-12                         	15195619	        80.61 ns/op	      32 B/op	       3 allocs/op
func BenchmarkLogger_Debugf(b *testing.B) {
	tests := []struct {
		timestampFunc timestamp.Timestamp
		description   string
		withJSON      bool
		withCaller    bool
	}{
		{
			description:   "1. standard timestamp",
			timestampFunc: timestamp.TimestampStandard,
		},
		{
			description:   "2. yyyy-month timestamp",
			timestampFunc: timestamp.TimestampYYYYMonth,
		},
		{
			description:   "3. nano timestamp",
			timestampFunc: timestamp.TimestampNano,
		},
		{
			description:   "4. nano timestamp - json",
			timestampFunc: timestamp.TimestampNano,
			withJSON:      true,
		},
		{
			description:   "5. nano timestamp - json, caller",
			timestampFunc: timestamp.TimestampNano,
			withJSON:      true,
			withCaller:    true,
		},
		{
			description: "6. nil timestamp",
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
						Ingestor:      ingestor,
						LoggerLevel:   Level(LevelDEBUG),
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
