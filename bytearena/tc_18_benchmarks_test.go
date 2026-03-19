package bytearena

import (
	"bytes"
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/log/helpers"
)

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkArena_ConstantPayload-16    	87987344	        13.46 ns/op	       0 B/op	       0 allocs/op
func BenchmarkArena_ConstantPayload(b *testing.B) {
	b.ReportAllocs()

	sink := helpers.CountWriter{}
	ingestor := NewIngestor(1024*1024, &sink)

	payload := []byte(`{"level":"info","msg":"user login","user_id":123}`)

	b.ResetTimer()

	for b.Loop() {
		ingestor.write(
			uint32(len(payload)),

			func(destination []byte) {
				copy(destination, payload)
			},
		)
	}

	_ = sink.TotalBytesWritten.Load()
}

// BenchmarkStandardLogger-16    	 9623833	       125.0 ns/op	      72 B/op	       2 allocs/op
func BenchmarkArena_FormattedPayload(b *testing.B) {
	b.ReportAllocs()

	sink := helpers.CountWriter{}
	ingestor := NewIngestor(1024, &sink)

	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		payload := fmt.Sprintf(
			`{"level":"info","msg":"user login","user_id":%d}`,
			i,
		)

		ingestor.write(
			uint32(len(payload)),
			func(destination []byte) {
				copy(destination, []byte(payload))
			},
		)
	}

	_ = sink.TotalBytesWritten.Load() // force sink to stay live
}

// go test -run '^$' -bench '^BenchmarkIngestor_Write$' -benchmem
// go test -run '^$' -bench '^BenchmarkIngestor_Write$' -benchmem -race
// BenchmarkIngestor_Write-16     92691792                12.31 ns/op            0 B/op          0 allocs/op
func BenchmarkIngestor_Write(b *testing.B) {
	writer := helpers.CountWriter{}

	ingestor := NewIngestor(Size1M, &writer)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	payload := []byte("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")

	b.ReportAllocs()
	b.ResetTimer()

	var written atomic.Int64

	for b.Loop() {
		if _, errWrite := ingestor.Write(payload); errWrite == nil {
			written.Add(1)
		}
	}

	cancel()
	<-chIngestionEnd

	b.Log(written.Load())

	require.EqualValues(b,
		writer.TotalBytesWritten.Load(),
		written.Load()*int64(len(payload)),
	)
}

// go test -run '^$' -bench '^BenchmarkIngestor_WriteParallel$' -benchmem
// go test -run '^$' -bench '^BenchmarkIngestor_WriteParallel$' -benchmem -race
// BenchmarkIngestor_WriteParallel-16    	21773329	        56.61 ns/op	       0 B/op	       0 allocs/op
func BenchmarkIngestor_WriteParallel(b *testing.B) {
	writer := helpers.CountWriter{}

	ingestor := NewIngestor(Size100K, &writer)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	payload := []byte("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")

	b.ReportAllocs()
	b.SetParallelism(16) // tune accordingly
	b.ResetTimer()

	var written atomic.Int64

	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				if _, errWrite := ingestor.Write(payload); errWrite == nil {
					written.Add(1)
				}
			}
		},
	)

	cancel()
	<-chIngestionEnd

	b.Log(written.Load())

	require.EqualValues(b,
		writer.TotalBytesWritten.Load(),
		written.Load()*int64(len(payload)),
	)
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkIngestor_MultipleSizes/size_16-16         	95403202	        11.65 ns/op	       0 B/op	       0 allocs/op
// BenchmarkIngestor_MultipleSizes/size_64-16         	91622713	        13.21 ns/op	       0 B/op	       0 allocs/op
// BenchmarkIngestor_MultipleSizes/size_256-16        	67990924	        17.30 ns/op	       0 B/op	       0 allocs/op
// BenchmarkIngestor_MultipleSizes/size_1024-16       	55336317	        21.57 ns/op	       0 B/op	       0 allocs/op
func BenchmarkIngestor_MultipleSizes(b *testing.B) {
	sizes := []int{16, 64, 256, 1024}

	for _, size := range sizes {
		b.Run(
			fmt.Sprintf("size_%d", size),

			func(b *testing.B) {
				ingestor := NewIngestor(
					Size1M,
					&helpers.CountWriter{},
				)

				ctx, cancel := context.WithCancel(context.Background())
				chIngestionEnd := ingestor.StartIngestion(ctx)

				payload := bytes.Repeat([]byte("x"), size)

				b.ReportAllocs()
				b.ResetTimer()

				for b.Loop() {
					_, _ = ingestor.Write(payload)
				}

				cancel()
				<-chIngestionEnd
			},
		)
	}
}
