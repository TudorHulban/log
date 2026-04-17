package log

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/bytearena/helpers"
)

func Test_GetLogLevel(t *testing.T) {
	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		&helpers.NoopWriter{},
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelDEBUG,
		},
	)
	require.NoError(t, errCrLogger)
	require.NotNil(t, l)

	require.EqualValues(t,
		LevelDEBUG,
		l.GetLogLevel(),
	)
}

func Test_SetLogLevel(t *testing.T) {
	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		&helpers.NoopWriter{},
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor: ingestor,
		},
	)
	require.NoError(t, errCrLogger)
	require.NotNil(t, l)

	l.SetLogLevel(LevelINFO)

	require.EqualValues(t,
		LevelINFO,
		l.GetLogLevel(),
	)
}
