package log

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/tudorhulban/log/helpers"
)

// BenchmarkZerolog_OneField-12    	 7189495	       167.2 ns/op	       0 B/op	       0 allocs/op
func BenchmarkZerolog_OneField(b *testing.B) {
	var writer helpers.NoopWriter

	logger := zerolog.New(&writer).With().
		Timestamp().
		Logger()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		logger.Info().
			Str("area", "some area").
			Msg("benchmark test")
	}
}

// BenchmarkZerolog_WithFields-12    	 5937602	       198.6 ns/op	       0 B/op	       0 allocs/op
func BenchmarkZerolog_WithFields(b *testing.B) {
	var writer helpers.NoopWriter

	logger := zerolog.New(&writer).With().
		Timestamp().
		Str("service", "auth").
		Int("req_id", 12345).
		Bool("cache_hit", true).
		Logger()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		logger.Info().
			Str("area", "some area").
			Str("user", "tudor").
			Int("attempt", i).
			Msg("benchmark test")
	}
}
