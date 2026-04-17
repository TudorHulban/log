package log

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/log/timestamp"
)

func TestCreateEmitData(t *testing.T) {
	var totals uint32 = 10

	state, errCrState := createEmitData(
		totals,
		map[Level]*uint32{},
		new(1),
	)
	require.NoError(t, errCrState)
	require.NotNil(t, state)

	require.Len(t, state.dictionary, int(totals))
}

func TestEmitData(t *testing.T) {
	var totals uint32 = 300

	data, errCrData := createEmitData(
		totals,
		map[Level]*uint32{},
		new(0),
	)
	require.NoError(t, errCrData)
	require.NotNil(t, data)

	writer := newTrackingWriter(
		os.Stdout,
		// io.Discard,
		data,
	)

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		writer,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelDEBUG,

			WithTimestamp: timestamp.TimestampRFC3339Bucharest,
			WithJSON:      true,
		},
	)
	require.NoError(t, errCrLogger)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	time.Sleep(10 * time.Millisecond) // warm up

	require.Empty(t,
		emitData(
			data,
			l,
		),
	)

	cancel()
	<-chIngestionEnd

	require.Zero(t,
		writer.UnreceivedCount(),
	)
}
