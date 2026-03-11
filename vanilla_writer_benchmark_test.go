package log

import (
	"bytes"
	"io"
	"strconv"
	"sync"
	"testing"
	"time"
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

// BenchmarkLogger-16    	12526819	        98.19 ns/op	      40 B/op	       2 allocs/op
func BenchmarkLoggerZeroAlloc(b *testing.B) {
	b.ReportAllocs()

	sink := io.Discard
	pool := sync.Pool{New: func() any {
		buf := &bytes.Buffer{}
		buf.Grow(256)
		return buf
	}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf := pool.Get().(*bytes.Buffer)
		buf.Reset()

		buf.WriteString(`{"level":"info","ts":"`)
		buf.Write(time.Now().AppendFormat(buf.AvailableBuffer(), time.RFC3339)) // no alloc
		buf.WriteString(`","msg":"user login","user_id":`)
		buf.Write(strconv.AppendInt(buf.AvailableBuffer(), int64(i), 10)) // no alloc
		buf.WriteByte('}')
		buf.WriteByte('\n')

		_, _ = sink.Write(buf.Bytes())
		pool.Put(buf)
	}
}
