package log

import (
	"context"
	"os"
	"testing"

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
	writer := os.Stdout

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		writer,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: Level(LevelDEBUG),

			WithTimestamp: timestamp.TimestampRFC3339Bucharest,
			WithJSON:      true,
		},
	)
	require.NoError(t, errCrLogger)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	var totals uint32 = 30

	data, errCrData := createEmitData(
		totals,
		map[Level]*uint32{},
		new(0),
	)
	require.NoError(t, errCrData)
	require.NotNil(t, data)

	require.Empty(t,
		emitData(
			data,
			l,
		),
	)

	cancel()
	<-chIngestionEnd
}
