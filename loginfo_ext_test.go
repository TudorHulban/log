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
	l *log.LogInfo
}

func Test1ELogger(t *testing.T) {
	output := &bytes.Buffer{}

	obj := T{
		l: log.New(3, output),
	}
	obj.l.Print("xxx")
	assert.Contains(t, output.String(), "xxx")
}
