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
// BenchmarkLogger_Parallel_Printf/1._standard_timestamp-16         	17041684	        74.35 ns/op	     112 B/op	       3 allocs/op
// BenchmarkLogger_Parallel_Printf/2._yyyy-month_timestamp-16       	16036966	        74.37 ns/op	     112 B/op	       3 allocs/op
// BenchmarkLogger_Parallel_Printf/3._nano_timestamp-16             	15739663	        75.68 ns/op	     112 B/op	       3 allocs/op
// BenchmarkLogger_Parallel_Printf/4._nano_timestamp_-_json-16      	12146430	        99.75 ns/op	     240 B/op	       5 allocs/op
// BenchmarkLogger_Parallel_Printf/5._nil_timestamp-16              	18202852	        66.88 ns/op	      72 B/op	       2 allocs/op
func BenchmarkLogger_Parallel_Printf(b *testing.B) {
	tests := []struct {
		timestampFunc timestamp.Timestamp
		description   string
		withJSON      bool
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
			description: "5. nil timestamp",
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
						logger.Printf(
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

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkLogger_Parallel_PrintfFast/1._standard_timestamp-16         	21208875	        58.62 ns/op	       8 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_PrintfFast/2._yyyy-month_timestamp-16       	20666210	        58.79 ns/op	       8 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_PrintfFast/3._nano_timestamp-16             	20455646	        58.64 ns/op	       8 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_PrintfFast/4._nano_timestamp_-_json-16      	20126348	        58.80 ns/op	       8 B/op	       0 allocs/op
// BenchmarkLogger_Parallel_PrintfFast/5._nil_timestamp-16              	20683252	        59.20 ns/op	       8 B/op	       0 allocs/op
func BenchmarkLogger_Parallel_PrintfFast(b *testing.B) {
	tests := []struct {
		timestampFunc timestamp.Timestamp
		description   string
		withJSON      bool
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
			description: "5. nil timestamp",
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
