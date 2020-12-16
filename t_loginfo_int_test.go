package log

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test1Logger(t *testing.T) {
	logger := New(3, os.Stderr, true)
	logger.Print("0")
	logger.Info("1")
	logger.Warn("2")
	logger.Debug("3")
}

// Test2Logger Tests if logger is created with correct log level.
func Test2GetLevel(t *testing.T) {
	l := New(3, os.Stderr, true)
	assert.Equal(t, 3, l.GetLogLevel())
}

// Test3Logger Tests if logger log level is updated.
func Test3SetLevel(t *testing.T) {
	l := New(3, os.Stderr, true)
	l.SetLogLevel(1)
	assert.Equal(t, 1, l.GetLogLevel())
}

// Test4Output Tests logger output.
func Test4Output(t *testing.T) {
	output := &bytes.Buffer{}
	l := New(0, output, true)
	l.Print("xxx")

	assert.Contains(t, output.String(), "xxx")
}

// Test5Output Tests logger Info level output.
func Test5Output(t *testing.T) {
	output := &bytes.Buffer{}
	l := New(1, output, true)

	l.Info("xxx")
	assert.Contains(t, output.String(), "xxx")
	l.Debug("zzz")
	assert.NotContains(t, output.String(), "zzz")
}
