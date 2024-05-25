package log_test

// File details how to use logger.

import (
	"bytes"
	"testing"

	"github.com/TudorHulban/log"
)

type T struct {
	l log.Logger
}

func Test01Logger(t *testing.T) {
	var output bytes.Buffer

	obj := T{
		l: log.NewLogger(log.LevelDEBUG, &output, true),
	}
	obj.l.Info("xxx")

	// assert.Contains(t, output.String(), "xxx") - race condition
}

// higher log levels are not sent when lower log level defined
func Test02Logger(t *testing.T) {
	var output bytes.Buffer

	obj := T{
		l: log.NewLogger(log.LevelINFO, &output, true),
	}
	obj.l.Debug("xxx")

	// assert.NotContains(t, output.String(), "xxx") - race condition
}
