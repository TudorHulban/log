package log

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test1_Logger(t *testing.T) {
	l := NewLogger(3, os.Stderr, true)

	l.Print("0")
	l.Info("1")
	l.Warn("2")
	l.Debug("3")
}

func Test2_GetLevel(t *testing.T) {
	l := NewLogger(3, os.Stderr, true)

	assert.Equal(t, 3, l.GetLogLevel())
}

func Test3_SetLevel(t *testing.T) {
	l := NewLogger(3, os.Stderr, true)

	l.SetLogLevel(1)

	assert.Equal(t, 1, l.GetLogLevel())
}

func Test4_OutputPrint(t *testing.T) {
	output := &bytes.Buffer{}

	l := NewLogger(0, output, true)
	l.Print("xxx")

	assert.Contains(t, output.String(), "xxx")
}

func Test5_OutputInfoLevel(t *testing.T) {
	output := &bytes.Buffer{}

	l := NewLogger(1, output, true)

	l.Info("xxx")
	assert.Contains(t, output.String(), "xxx")

	l.Infof("%d", 1)
	assert.Contains(t, output.String(), "1")

	l.Debug("zzz")
	assert.NotContains(t, output.String(), "zzz")
}

func Test5_NewOutputInfoLevel(t *testing.T) {
	var output bytes.Buffer

	l := NewLogger(1, &output, true)

	l.Info("xxx")
	assert.Contains(t, output.String(), "xxx")

	l.Infof("%d", 1)
	assert.Contains(t, output.String(), "1")

	l.Debug("zzz")
	assert.NotContains(t, output.String(), "zzz")
}

func Test6_Fatal(t *testing.T) {
	output := os.Stdout

	l := NewLogger(0, output, false)

	l.Fatal("xxx")
}
