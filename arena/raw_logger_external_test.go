package arena_test

import (
	"bytes"
	"context"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/log/arena"
	"github.com/tudorhulban/log/helpers"
)

func TestHowToUse(t *testing.T) {
	var sink bytes.Buffer

	rawLogger := arena.NewRawLogger(arena.Size100K, &sink)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := rawLogger.StartIngestion(ctx)

	payload := "xxx"

	bytesWritten, errWrite := rawLogger.Write([]byte(payload))
	require.NoError(t, errWrite)
	require.Equal(t, len(payload), bytesWritten)

	cancel()
	<-chIngestionEnd

	require.Contains(t, sink.String(), payload)
}

// go test -run '^$' -bench '^BenchmarkRawLogger_Write$' -benchmem
// go test -run '^$' -bench '^BenchmarkRawLogger_Write$' -benchmem -race
// BenchmarkRawLogger_Write-16     92691792                12.31 ns/op            0 B/op          0 allocs/op
func BenchmarkRawLogger_Write(b *testing.B) {
	writer := helpers.CountWriter{}

	rawLogger := arena.NewRawLogger(arena.Size1M, &writer)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := rawLogger.StartIngestion(ctx)

	payload := []byte("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx") // fixed size

	b.ReportAllocs()
	b.ResetTimer()

	var written atomic.Int64

	for b.Loop() {
		if _, errWrite := rawLogger.Write(payload); errWrite == nil {
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

// go test -run '^$' -bench '^BenchmarkRawLogger_WriteParallel$' -benchmem
// go test -run '^$' -bench '^BenchmarkRawLogger_WriteParallel$' -benchmem -race
// BenchmarkRawLogger_WriteParallel-16    	17507642	        67.72 ns/op	       0 B/op	       0 allocs/op
func BenchmarkRawLogger_WriteParallel(b *testing.B) {
	writer := helpers.CountWriter{}

	rawLogger := arena.NewRawLogger(arena.Size100K, &writer)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := rawLogger.StartIngestion(ctx)

	payload := []byte("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")

	b.ReportAllocs()
	b.SetParallelism(16) // tune accordingly
	b.ResetTimer()

	var written atomic.Int64

	b.RunParallel(
		func(pb *testing.PB) {
			for pb.Next() {
				if _, errWrite := rawLogger.Write(payload); errWrite == nil {
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
