package log

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test1Logger(t *testing.T) {
	l := NewLogger(3, os.Stderr, true)

	l.Print("0")
	l.Info("1")
	l.Warn("2")
	l.Debug("3")
}

// Test2Logger Tests if logger is created with correct log level.
func Test2GetLevel(t *testing.T) {
	l := NewLogger(3, os.Stderr, true)
	assert.Equal(t, 3, l.GetLogLevel())
}

// Test3Logger Tests if logger log level is updated.
func Test3SetLevel(t *testing.T) {
	l := NewLogger(3, os.Stderr, true)
	l.SetLogLevel(1)
	assert.Equal(t, 1, l.GetLogLevel())
}

// Test4Output Tests logger output.
func Test4Output(t *testing.T) {
	output := &bytes.Buffer{}
	l := NewLogger(0, output, true)
	l.Print("xxx")

	assert.Contains(t, output.String(), "xxx")
}

// Test5Output Tests logger Info level output.
func Test5Output(t *testing.T) {
	output := &bytes.Buffer{}
	l := NewLogger(1, output, true)

	l.Info("xxx")
	assert.Contains(t, output.String(), "xxx")
	l.Debug("zzz")
	assert.NotContains(t, output.String(), "zzz")
}
