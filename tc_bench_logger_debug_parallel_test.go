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
// BenchmarkLogger_Parallel_DebugFast/1.nil_timestamp-16         	18869478	        65.65 ns/op	      40 B/op	       2 allocs/op
// BenchmarkLogger_Parallel_DebugFast/2.standard_timestamp-16    	18107064	        64.76 ns/op	      40 B/op	       1 allocs/op
// BenchmarkLogger_Parallel_DebugFast/3.yyyy-month_timestamp-16  	17515292	        65.47 ns/op	      40 B/op	       2 allocs/op
// BenchmarkLogger_Parallel_DebugFast/4.nano_timestamp-16        	17986500	        65.53 ns/op	      40 B/op	       1 allocs/op
// BenchmarkLogger_Parallel_DebugFast/5.nano_timestamp_-_json-16 	17570806	        65.62 ns/op	      40 B/op	       2 allocs/op
// BenchmarkLogger_Parallel_DebugFast/6.nano_-_json,_caller-16   	17852004	        65.50 ns/op	      40 B/op	       2 allocs/op
func BenchmarkLogger_Parallel_DebugFast(b *testing.B) {
	runtime.GOMAXPROCS(1)

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

				b.SetParallelism(16)
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
