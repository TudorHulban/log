package log

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/bytearena/helpers"
	"github.com/tudorhulban/log/timestamp"
)

// go test -run '^$' -bench '^BenchmarkLogger_Print$' -benchmem

// BenchmarkLogger_Print-16    	19233896	        63.44 ns/op	      24 B/op	       2 allocs/op
func BenchmarkLogger_Print(b *testing.B) {
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
			LoggerLevel: LevelINFO,
		},
	)
	require.NoError(b, errCrLogger)
	require.NotNil(b, logger)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		logger.Print("hi", 123, "world")
	}

	cancel()
	<-chIngestionEnd
}

// BenchmarkLogger_PrintWithNoTimestamp-16    	30770088	        39.77 ns/op	       0 B/op	       0 allocs/op
func BenchmarkLogger_PrintWithNoTimestamp(b *testing.B) {
	var writer helpers.NoopWriter

	ingestor, errCrIngestor := bytearena.NewIngestor(bytearena.Size100K(), &writer)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelINFO,
		},
	)
	require.NoError(b, errCrLogger)
	require.NotNil(b, logger)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		logger.PrintWithNoTimestampFast("hi", 123, "world")
	}

	cancel()
	<-chIngestionEnd
}

// BenchmarkLogger_PrintRaw-16    	70315348	        16.87 ns/op	       0 B/op	       0 allocs/op
func BenchmarkLogger_PrintRaw(b *testing.B) {
	var writer helpers.NoopWriter

	ingestor, errCrIngestor := bytearena.NewIngestor(bytearena.Size100K(), &writer)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelINFO,
		},
	)
	require.NoError(b, errCrLogger)
	require.NotNil(b, logger)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		logger.PrintRaw(
			[]byte("xxxxxxxxxxxxxxxxxxx"),
		)
	}

	cancel()
	<-chIngestionEnd
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkLogger_Printf/1.nil_timestamp-16         	 9420816	       130.3 ns/op	      72 B/op	       2 allocs/op
// BenchmarkLogger_Printf/2.standard_timestamp-16    	 5771464	       210.5 ns/op	     200 B/op	       3 allocs/op
// BenchmarkLogger_Printf/3.yyyy-month_timestamp-16  	 5723036	       210.7 ns/op	     200 B/op	       3 allocs/op
// BenchmarkLogger_Printf/4.nano_timestamp-16        	 5187085	       232.3 ns/op	     200 B/op	       3 allocs/op
// BenchmarkLogger_Printf/5.nano_timestamp_-_json-16 	 3347700	       360.8 ns/op	     360 B/op	       4 allocs/op
func BenchmarkLogger_Printf(b *testing.B) {
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
						LoggerLevel: LevelINFO,

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
					logger.Printf(
						`{"level":"info","msg":"user login","user_id":%d}`,
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
// BenchmarkLogger_PrintFast/1._standard_timestamp-16         	19161618	        61.47 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_PrintFast/2._yyyy-month_timestamp-16       	18268606	        63.50 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_PrintFast/3._nano_timestamp-16             	17402479	        66.70 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_PrintFast/4._nano_timestamp_-_json-16      	17916865	        67.46 ns/op	       8 B/op	       1 allocs/op
// BenchmarkLogger_PrintFast/5._nil_timestamp-16              	23356777	        50.03 ns/op	       8 B/op	       1 allocs/op
func BenchmarkLogger_PrintFast(b *testing.B) {
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
						LoggerLevel:   LevelINFO,
						WithTimestamp: tcase.timestampFunc,
						WithJSON:      tcase.withJSON,
					},
				)
				require.NoError(b, errCrLogger)
				require.NotNil(b, logger)

				b.ReportAllocs()
				b.ResetTimer()

				for i := 0; b.Loop(); i++ {
					logger.PrintFast(
						`{"level":"info","msg":"user login","user_id":%d}`,
						i,
					)
				}

				cancel()
				<-chIngestionEnd
			},
		)
	}
}
