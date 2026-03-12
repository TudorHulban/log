package log

import (
	"testing"

	"log"
)

// BenchmarkStandardLogger-16    	 6128383	       196.9 ns/op	       8 B/op	       0 allocs/op
func BenchmarkStandardLogger(b *testing.B) {
	b.ReportAllocs()

	sink := countWriter{}

	log.SetOutput(&sink)
	log.SetFlags(log.LstdFlags)

	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		log.Printf(
			`{"level":"info","msg":"user login","user_id":%d}`,
			i,
		)
	}

	_ = sink.n.Load() // force sink to stay live
}

// and parallel — stdlib log has a global mutex, this exposes it
// BenchmarkStandardLoggerParallel-16    	 8015534	       152.3 ns/op	       0 B/op	       0 allocs/op
func BenchmarkStandardLoggerParallel(b *testing.B) {
	b.ReportAllocs()

	sink := countWriter{}

	log.SetOutput(&sink)
	log.SetFlags(log.LstdFlags)

	b.ResetTimer()

	b.RunParallel(
		func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				log.Printf(
					`{"level":"info","msg":"user login","user_id":%d}`,
					i, // ← variable
				)

				i++
			}
		},
	)
}
