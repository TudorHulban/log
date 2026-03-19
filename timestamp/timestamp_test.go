package timestamp

import (
	"testing"
)

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkTimestamp-16    	19970774	        59.96 ns/op	      24 B/op	       1 allocs/op
func BenchmarkTimestamp(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		TimestampYYYYMonth()
	}
}
