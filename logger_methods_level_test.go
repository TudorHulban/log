package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetLogLevel(t *testing.T) {
	l := NewLogger(LevelDEBUG, nil, true)

	assert.EqualValues(t,
		LevelDEBUG,
		l.GetLogLevel(),
	)
}

func Test_SetLogLevel(t *testing.T) {
	l := NewLogger(LevelDEBUG, nil, true)

	l.SetLogLevel(LevelINFO)

	assert.EqualValues(t,
		LevelINFO,
		l.GetLogLevel(),
	)
}
