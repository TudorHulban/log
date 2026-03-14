package arena

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

	m := NewManager(1024, &out)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	done := m.StartConsumer(ctx)

	payload := "hi!"

	require.True(t,
		m.Write(
			int64(len(payload)),
			func(dst []byte) {
				copy(dst, []byte(payload))
			},
		),
	)

	// Wait for consumer shutdown flush.
	<-done

	require.Equal(t,
		payload,
		out.String(),
	)
}

// BenchmarkStandardLogger-16    	 9623833	       125.0 ns/op	      72 B/op	       2 allocs/op
func BenchmarkStandardLogger(b *testing.B) {
	b.ReportAllocs()

	sink := helpers.CountWriter{}

	manager := NewManager(1024, &sink)

	b.ResetTimer()

	for i := 0; b.Loop(); i++ {
		payload := fmt.Sprintf(
			`{"level":"info","msg":"user login","user_id":%d}`,
			i,
		)

		manager.Write(
			int64(len(payload)),
			func(dst []byte) {
				copy(dst, []byte(payload))
			},
		)
	}

	_ = sink.N.Load() // force sink to stay live
}

// BenchmarkArenaWrite-16    	78782551	        14.73 ns/op	       0 B/op	       0 allocs/op
func BenchmarkArenaWrite(b *testing.B) {
	b.ReportAllocs()

	sink := helpers.CountWriter{}
	manager := NewManager(1024*1024, &sink)

	payload := []byte(`{"level":"info","msg":"user login","user_id":123}`)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		manager.Write(
			int64(len(payload)),
			func(dst []byte) {
				copy(dst, payload)
			},
		)
	}

	_ = sink.N.Load()
}
