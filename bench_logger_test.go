package log

import (
	"testing"

	"github.com/tudorhulban/log/timestamp"
)

// BenchmarkLogger-16    	 8522778	       145.4 ns/op	       8 B/op	       0 allocs/op
func BenchmarkLogger(b *testing.B) {
	b.ReportAllocs()

	sink := countWriter{}

	logger := NewLogger(
		&ParamsNewLogger{
			LoggerWriter:  &sink,
			LoggerLevel:   LevelINFO,
			WithTimestamp: timestamp.TimestampNano,
		},
	)

	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		logger.Printf(
			`{"level":"info","msg":"user login","user_id":%d}`,
			i,
		)
	}

	_ = sink.n.Load() // force sink to stay live
}
