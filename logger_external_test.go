package log_test

// File details how to use logger.

import (
	"bytes"
	"testing"

	"github.com/TudorHulban/log"
	"github.com/TudorHulban/log/timestamp"
)

type T struct {
	l log.Logger
}

func TestExternal(t *testing.T) {
	obj := T{
		l: log.NewLogger(
			&log.ParamsNewLogger{
				LoggerLevel:  log.LevelDEBUG,
				LoggerWriter: new(bytes.Buffer),

				WithTimestamp: timestamp.TimestampNano,
				WithCaller:    true,
				WithColor:     true,
			},
		),
	}

	obj.l.Info("xxx")
	obj.l.Debug("yyy")

	// assert.Contains(t, output.String(), "xxx") - race condition
}
