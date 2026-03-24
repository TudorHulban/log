package helpers

import "testing"

func BenchmarkNoopWriter(b *testing.B) {
	w := NoopWriter{}
	buf := make([]byte, 128)

	b.ReportAllocs()

	for i := 0; b.Loop(); i++ {
		w.Write(buf)
	}
}
