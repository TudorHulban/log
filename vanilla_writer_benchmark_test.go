package log

import (
	"io"
	"testing"
)

// Benchmark_Vanilla_Logger-16    	100000000	        10.07 ns/op	       1 B/op	       1 allocs/op
func Benchmark_Vanilla_Logger_Parallel(b *testing.B) {
	w := io.Discard

	b.ResetTimer()

	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				w.Write(
					[]byte(
						"1",
					),
				)
			}
		},
	)
}

// Benchmark_Vanilla_Logger-16    	24263023	        53.31 ns/op	       1 B/op	       1 allocs/op
func Benchmark_Vanilla_Logger(b *testing.B) {
	w := io.Discard

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		w.Write(
			[]byte(
				"1",
			),
		)
	}
}
