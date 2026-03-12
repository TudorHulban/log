package log

import (
	"testing"

	"github.com/tudorhulban/log/timestamp"
)

// BenchmarkLogger-16    	 5513952	       218.3 ns/op	     288 B/op	       3 allocs/op
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
