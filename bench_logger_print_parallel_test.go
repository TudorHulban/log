package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/log/helpers"
	"github.com/tudorhulban/log/timestamp"
)

// BenchmarkLogger_Parallel_PrintRaw-16    	22106994	        56.75 ns/op	       0 B/op	       0 allocs/op
func BenchmarkLogger_Parallel_PrintRaw(b *testing.B) {
	var writer helpers.NoopWriter

	ingestor, errCrIngestor := bytearena.NewIngestor(bytearena.Size100K(), &writer)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: Level(LevelINFO),
		},
	)
	require.NoError(b, errCrLogger)
	require.NotNil(b, logger)

	b.SetParallelism(1)
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				logger.PrintRaw([]byte("xxxxxxxxxxxxxxxxxxx"))
			}
		},
	)

	cancel()
	<-chIngestionEnd
}

// go test -run '^$' -bench '^BenchmarkLogger_Parallel_Printf$' -benchmem -race

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkLogger_Parallel_Printf/1.nil_timestamp-16         	19088084	        66.20 ns/op	      72 B/op	       2 allocs/op
// BenchmarkLogger_Parallel_Printf/2.standard_timestamp-16    	14446918	        83.74 ns/op	     200 B/op	       3 allocs/op
// BenchmarkLogger_Parallel_Printf/3.yyyy-month_timestamp-16  	14160460	        84.73 ns/op	     200 B/op	       3 allocs/op
// BenchmarkLogger_Parallel_Printf/4.nano_timestamp-16        	13746822	        87.33 ns/op	     200 B/op	       3 allocs/op
// BenchmarkLogger_Parallel_Printf/5.nano_timestamp_-_json-16 	 9217291	       129.7 ns/op	     360 B/op	       4 allocs/op
func BenchmarkLogger_Parallel_Printf(b *testing.B) {
	tests := []struct {
		timestampFunc timestamp.Timestamp
		description   string
		withJSON      bool
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
						LoggerLevel:   Level(LevelINFO),
						WithTimestamp: tcase.timestampFunc,
						WithJSON:      tcase.withJSON,
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
							logger.Printf(
								`{"level":"info","msg":"user login","user_id":%d}`,
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

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkLogger_Parallel_PrintfFast/1._nil_timestamp-16         	20184763	        61.00 ns/op	      24 B/op	       2 allocs/op
// BenchmarkLogger_Parallel_PrintfFast/2._standard_timestamp-16    	18928455	        61.58 ns/op	      24 B/op	       2 allocs/op
// BenchmarkLogger_Parallel_PrintfFast/3._yyyy-month_timestamp-16  	19693173	        61.62 ns/op	      24 B/op	       2 allocs/op
// BenchmarkLogger_Parallel_PrintfFast/4._nano_timestamp-16        	19073270	        62.06 ns/op	      24 B/op	       2 allocs/op
// BenchmarkLogger_Parallel_PrintfFast/5._nano_timestamp_-_json-16 	18923116	        63.67 ns/op	      26 B/op	       2 allocs/op
func BenchmarkLogger_Parallel_PrintfFast(b *testing.B) {
	tests := []struct {
		timestampFunc timestamp.Timestamp
		description   string
		withJSON      bool
	}{
		{
			description: "1. nil timestamp",
		},
		{
			description:   "2. standard timestamp",
			timestampFunc: timestamp.TimestampStandard,
		},
		{
			description:   "3. yyyy-month timestamp",
			timestampFunc: timestamp.TimestampYYYYMonth,
		},
		{
			description:   "4. nano timestamp",
			timestampFunc: timestamp.TimestampNano,
		},
		{
			description:   "5. nano timestamp - json",
			timestampFunc: timestamp.TimestampNano,
			withJSON:      true,
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
						LoggerLevel:   Level(LevelINFO),
						WithTimestamp: tcase.timestampFunc,
						WithJSON:      tcase.withJSON,
					},
				)
				require.NoError(b, errCrLogger)
				require.NotNil(b, logger)

				b.SetParallelism(1)
				b.ReportAllocs()
				b.ResetTimer()

				b.RunParallel(func(pb *testing.PB) {
					i := 0
					for pb.Next() {
						logger.PrintfFast(
							`{"level":"info","msg":"user login","user_id":%d}`,
							i,
						)
						i++
					}
				})

				cancel()
				<-chIngestionEnd
			},
		)
	}
}
