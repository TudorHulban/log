package log

import (
	"testing"

	"log"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena/helpers"
)

// BenchmarkStandardLogger-16    	 5019615	       233.7 ns/op	     168 B/op	       0 allocs/op
func BenchmarkStandardLogger(b *testing.B) {
	writer := helpers.CountWriterNoBuffer{}

	log.SetOutput(&writer)
	log.SetFlags(log.LstdFlags)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		log.Printf(
			`{"level":"info","msg":"user login","user_id":%d}`,
			i,
		)
	}

	require.NotZero(b,
		writer.TotalBytesWritten.Load(), // force writer to stay live
	)
}

// and parallel — stdlib log has a global mutex, this exposes it
// BenchmarkStandardLoggerParallel-16    	 6440850	       191.1 ns/op	     258 B/op	       0 allocs/op
func BenchmarkStandardLoggerParallel(b *testing.B) {
	writer := helpers.CountWriterNoBuffer{}

	log.SetOutput(&writer)
	log.SetFlags(log.LstdFlags)

	b.ReportAllocs()
	b.ResetTimer()
	b.SetParallelism(1)

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

	require.NotZero(b,
		writer.TotalBytesWritten.Load(), // force writer to stay live
	)
}
