package log

import (
	"testing"

	"github.com/TudorHulban/log/timestamp"
	"github.com/stretchr/testify/assert"
)

func Test_GetLogLevel(t *testing.T) {
	l := NewLogger(
		&ParamsNewLogger{
			LoggerLevel:   LevelDEBUG,
			WithTimestamp: timestamp.TimestampNil,
		},
	)

	assert.EqualValues(t,
		LevelDEBUG,
		l.GetLogLevel(),
	)
}

func Test_SetLogLevel(t *testing.T) {
	l := NewLogger(
		&ParamsNewLogger{
			WithTimestamp: timestamp.TimestampNil,
		},
	)

	l.SetLogLevel(LevelINFO)

	assert.EqualValues(t,
		LevelINFO,
		l.GetLogLevel(),
	)
}
