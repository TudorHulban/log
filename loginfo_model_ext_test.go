package log_test

/*
File details how to use logger.
*/

import (
	"bytes"
	"testing"

	"github.com/TudorHulban/log"
	"github.com/stretchr/testify/assert"
)

type T struct {
	l *log.Logger
}

func Test01Logger(t *testing.T) {
	output := &bytes.Buffer{}

	obj := T{
		l: log.NewLogger(log.DEBUG, output, true),
	}
	obj.l.Info("xxx")
	assert.Contains(t, output.String(), "xxx")
}

// higher log levels are not sent when lower log level defined
func Test02Logger(t *testing.T) {
	output := &bytes.Buffer{}

	obj := T{
		l: log.NewLogger(log.INFO, output, true),
	}
	obj.l.Debug("xxx")
	assert.NotContains(t, output.String(), "xxx")
}
