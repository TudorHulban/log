package timestamp

import (
	"testing"
)

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkTimestamp-16    	  755229	      1443 ns/op	      31 B/op	       2 allocs/op

func BenchmarkTimestamp(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		TimestampYYYYMonth()
	}
}
