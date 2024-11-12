package log

import (
	"io"
	"testing"
)

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// Benchmark_Vanilla_Logger_Parallel-16    	908120602	         1.316 ns/op	       1 B/op	       1 allocs/op
func Benchmark_Vanilla_Logger_Parallel(b *testing.B) {
	w := io.Discard

	b.ResetTimer()

	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				_, _ = w.Write(
					[]byte(
						"1",
					),
				)
			}
		},
	)
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// Benchmark_Vanilla_Logger-16    	135784692	         8.801 ns/op	       1 B/op	       1 allocs/op
func Benchmark_Vanilla_Logger(b *testing.B) {
	w := io.Discard

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = w.Write(
			[]byte(
				"1",
			),
		)
	}
}
