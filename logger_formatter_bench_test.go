package log

import (
	"bytes"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/log/bytearena"
)

// BenchmarkFormatterWithRequest-16    	 5000041	       227.8 ns/op	     692 B/op	       3 allocs/op
func BenchmarkFormatterWithRequest(b *testing.B) {
	var writer bytes.Buffer

	ingestor, errCrIngestor := bytearena.NewIngestor(bytearena.Size100K, &writer)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	defer func() {
		cancel()
		<-chIngestionEnd
	}()

	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: Level(LevelDEBUG),
		},
	)
	require.NoError(b, errCrLogger)

	logContext := NewLogContext(logger).
		WithRoot("service", "auth").
		SetInt("req_id", 12345).
		SetBool("cache_hit", true)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		// 1. Create request with 4 attributes
		entry := logContext.With("area", "some area").
			With("user", "tudor").
			With("attempt", i).
			With("success", true)

		// 2. Print
		entry.Print("benchmark test")
	}
}
