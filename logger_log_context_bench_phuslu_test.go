package log

import (
	"testing"
	"time"

	"github.com/phuslu/log"
	"github.com/tudorhulban/log/helpers"
)

// BenchmarkPhuslu_OneField-12    	 9210207	       129.9 ns/op	       0 B/op	       0 allocs/op
func BenchmarkPhuslu_OneField(b *testing.B) {
	var writer helpers.NoopWriter

	logger := log.Logger{
		Level:      log.InfoLevel,
		TimeFormat: time.RFC3339,
		Writer:     &log.IOWriter{Writer: &writer},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		logger.Info().
			Str("area", "some area").
			Msg("benchmark test")
	}
}

// BenchmarkPhuslu_WithFields-12    	 6902982	       174.4 ns/op	       0 B/op	       0 allocs/op
func BenchmarkPhuslu_WithFields(b *testing.B) {
	var writer helpers.NoopWriter

	logger := log.Logger{
		Level:      log.InfoLevel,
		TimeField:  "ts",
		TimeFormat: time.RFC3339,
		Writer:     &log.IOWriter{Writer: &writer},
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		logger.Info().
			Str("service", "auth").
			Int("req_id", 12345).
			Str("area", "some area").
			Str("user", "tudor").
			Int("attempt", i).
			Bool("success", true).
			Msg("benchmark test")
	}
}
