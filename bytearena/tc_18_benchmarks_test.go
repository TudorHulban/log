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
// BenchmarkArena_ConstantPayload-16    	84908293	        13.89 ns/op	       0 B/op	       0 allocs/op
func BenchmarkArena_ConstantPayload(b *testing.B) {
	writer := helpers.CountWriter{}

	ingestor, errCrIngestor := NewIngestor(1024*1024, &writer)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	payload := []byte(`{"level":"info","msg":"user login","user_id":123}`)

	b.ReportAllocs()
	b.ResetTimer()

	for b.Loop() {
		ingestor.write(
			uint32(len(payload)),

			func(destination []byte) {
				copy(destination, payload)
			},
		)
	}

	_ = writer.TotalBytesWritten.Load()
}

// BenchmarkArena_FormattedPayload-16    	11129125	       109.9 ns/op	      64 B/op	       1 allocs/op
func BenchmarkArena_FormattedPayload(b *testing.B) {
	writer := helpers.CountWriter{}

	ingestor, errCrIngestor := NewIngestor(1024, &writer)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		payload := helpers.SprintfInt(
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

	_ = writer.TotalBytesWritten.Load() // force sink to stay live
}

// go test -run '^$' -bench '^BenchmarkIngestor_Write$' -benchmem
// go test -run '^$' -bench '^BenchmarkIngestor_Write$' -benchmem -race
// BenchmarkIngestor_Write-16     92691792                12.31 ns/op            0 B/op          0 allocs/op
func BenchmarkIngestor_Write(b *testing.B) {
	writer := helpers.CountWriter{}

	ingestor, errCrIngestor := NewIngestor(Size1M, &writer)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

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

	ingestor, errCrIngestor := NewIngestor(Size100K, &writer)
	require.NoError(b, errCrIngestor)
	require.NotNil(b, ingestor)

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
// BenchmarkIngestor_MultipleSizes/size_msg16_arena102400-16         	100000000	        11.76 ns/op	       0 B/op	       0 allocs/op
// BenchmarkIngestor_MultipleSizes/size_msg64_arena102400-16         	93691624	        12.98 ns/op	       0 B/op	       0 allocs/op
// BenchmarkIngestor_MultipleSizes/size_msg256_arena102400-16        	77943147	        15.41 ns/op	       0 B/op	       0 allocs/op
// BenchmarkIngestor_MultipleSizes/size_msg1024_arena102400-16       	67931802	        17.33 ns/op	       0 B/op	       0 allocs/op

// BenchmarkIngestor_MultipleSizes/size_msg16_arena512000-16         	100000000	        11.77 ns/op	       0 B/op	       0 allocs/op
// BenchmarkIngestor_MultipleSizes/size_msg64_arena512000-16         	93074583	        13.06 ns/op	       0 B/op	       0 allocs/op
// BenchmarkIngestor_MultipleSizes/size_msg256_arena512000-16        	78160538	        15.28 ns/op	       0 B/op	       0 allocs/op
// BenchmarkIngestor_MultipleSizes/size_msg1024_arena512000-16       	68533010	        17.50 ns/op	       0 B/op	       0 allocs/op

// BenchmarkIngestor_MultipleSizes/size_msg16_arena1048576-16        	100000000	        11.95 ns/op	       0 B/op	       0 allocs/op
// BenchmarkIngestor_MultipleSizes/size_msg64_arena1048576-16        	93492704	        13.12 ns/op	       0 B/op	       0 allocs/op
// BenchmarkIngestor_MultipleSizes/size_msg256_arena1048576-16       	79064127	        15.30 ns/op	       0 B/op	       0 allocs/op
// BenchmarkIngestor_MultipleSizes/size_msg1024_arena1048576-16      	62103726	        17.42 ns/op	       0 B/op	       0 allocs/op
func BenchmarkIngestor_MultipleSizes(b *testing.B) {
	sizesMessage := []int{16, 64, 256, 1024}
	sizesArena := []int{Size100K, Size500K, Size1M}

	for _, sizeArena := range sizesArena {
		for _, sizeMessage := range sizesMessage {
			b.Run(
				fmt.Sprintf(
					"size_msg%d_arena%d",
					sizeMessage,
					sizeArena,
				),

				func(b *testing.B) {
					ingestor, errCrIngestor := NewIngestor(
						Size1M,
						&helpers.CountWriter{},
					)
					require.NoError(b, errCrIngestor)
					require.NotNil(b, ingestor)

					ctx, cancel := context.WithCancel(context.Background())
					chIngestionEnd := ingestor.StartIngestion(ctx)

					payload := bytes.Repeat([]byte("x"), sizeMessage)

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
}
