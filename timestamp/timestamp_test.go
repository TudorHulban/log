package timestamp

import (
	"testing"
)

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkTimestamp-16    	27195747	        43.94 ns/op	       0 B/op	       0 allocs/op
func BenchmarkTimestamp(b *testing.B) {
	var scratch [64]byte

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		TimestampYYYYMonth(scratch[:0])
	}
}
