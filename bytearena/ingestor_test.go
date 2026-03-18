package bytearena

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/log/helpers"
)

func TestManagerSingleWrite(t *testing.T) {
	var out bytes.Buffer

	rawLogger := NewIngestor(1024, &out)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	chIngestionEnd := rawLogger.StartIngestion(ctx)

	payload := "hi!"

	require.NoError(t,
		rawLogger.write(
			uint32(len(payload)),

			func(destination []byte) {
				copy(destination, []byte(payload))
			},
		),
	)

	// Wait for consumer shutdown flush.
	<-chIngestionEnd

	require.Equal(t,
		payload,
		out.String(),
	)
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

// BenchmarkArenaWrite-16    	67048297	        18.09 ns/op	       0 B/op	       0 allocs/op
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
