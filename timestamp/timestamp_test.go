package timestamp

import (
	"testing"
)

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkTimestamp-16    	 6041724	       207.3 ns/op	      31 B/op	       2 allocs/op
func BenchmarkTimestamp(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		TimestampYYYYMonth(nil)
	}
}
