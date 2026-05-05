package timestamp

import (
	"bytes"
	"testing"
	"time"
)

func TestTimestampRFC3339(t *testing.T) {
	t.Run(
		"1. Test basic appending to an empty slice",
		func(t *testing.T) {
			got := TimestampRFC3339(nil)

			// Validate format by attempting to parse it back
			_, errParse := time.Parse(time.RFC3339Nano, string(got))
			if errParse != nil {
				t.Errorf("Result is not a valid RFC3339 timestamp: %v", errParse)
			}
		},
	)

	t.Run(
		"2. Test appending to an existing slice (functional integrity)",
		func(t *testing.T) {
			prefix := []byte("log_prefix: ")
			got := TimestampRFC3339(prefix)

			if !bytes.HasPrefix(got, prefix) {
				t.Fatalf("Expected prefix not found. Got: %s", string(got))
			}

			timestampPart := got[len(prefix):]
			_, errParse := time.Parse(time.RFC3339Nano, string(timestampPart))
			if errParse != nil {
				t.Errorf("Appended part is not a valid timestamp: %v", errParse)
			}
		},
	)

	t.Run(
		"3. Consistency check (values should be close to time.Now)",
		func(t *testing.T) {
			now := time.Now().UTC()
			got := TimestampRFC3339(nil)
			parsed, _ := time.Parse(time.RFC3339Nano, string(got))

			// Check if the timestamp is within a reasonable drift (e.g., 1 second)
			// because the cache update might be slightly behind time.Now()
			diff := parsed.Sub(now)
			if diff < 0 {
				diff = -diff
			}

			if diff > time.Second {
				t.Errorf("Timestamp drift too high: %v", diff)
			}
		},
	)
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkTimestamp-16    	27195747	        43.94 ns/op	       0 B/op	       0 allocs/op
func BenchmarkTimestampYYYYMonth(b *testing.B) {
	var scratch [64]byte

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		TimestampYYYYMonth(scratch[:0])
	}
}

// BenchmarkTimestampRFC3339-16    	26308006	        45.12 ns/op	       0 B/op	       0 allocs/op
func BenchmarkTimestampRFC3339(b *testing.B) {
	var scratch [64]byte

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		TimestampRFC3339(scratch[:0])
	}
}
