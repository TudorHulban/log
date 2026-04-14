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
// BenchmarkLogger_Debugf/1._standard_timestamp-16         	 3771642	       324.0 ns/op	     134 B/op	       6 allocs/op
// BenchmarkLogger_Debugf/2._yyyy-month_timestamp-16       	 3736538	       330.1 ns/op	     134 B/op	       6 allocs/op
// BenchmarkLogger_Debugf/3._nano_timestamp-16             	 3522434	       343.5 ns/op	     133 B/op	       6 allocs/op
// BenchmarkLogger_Debugf/4._nano_timestamp_-_json-16      	 4975443	       242.9 ns/op	     136 B/op	       6 allocs/op
// BenchmarkLogger_Debugf/5._nano_timestamp_-_json-16      	 5024364	       242.1 ns/op	     136 B/op	       6 allocs/op
// BenchmarkLogger_Debugf/6._nil_timestamp-16              	 4602891	       281.1 ns/op	     111 B/op	       6 allocs/op
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
			description:   "5. nano timestamp - json",
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
