package loginfo

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createLogger(level int, t *testing.T) *LogInfo {
	l, err := New(level, os.Stderr)
	if assert.Nil(t, err) {
		return l
	}
	return &LogInfo{}
}

func Test1Logger(t *testing.T) {
	logger := createLogger(3, t)
	logger.Print("0")
	logger.Info("1")
	logger.Warn("2")
	logger.Debug("3")
}

// Test2Logger Tests if logger is created with correct log level.
func Test2GetLevel(t *testing.T) {
	l := createLogger(3, t)
	assert.Equal(t, 3, l.GetLogLevel())
}

// Test3Logger Tests if logger log level is updated.
func Test3SetLevel(t *testing.T) {
	l := createLogger(3, t)
	l.SetLogLevel(1)
	assert.Equal(t, 1, l.GetLogLevel())
}

// Test4Output Tests logger output.
func Test4Output(t *testing.T) {
	output := &bytes.Buffer{}
	l, err := New(0, output)
	l.Print("xxx")

	if assert.Nil(t, err) {
		assert.Contains(t, output.String(), "xxx")
	}
}

// Test5Output Tests logger Info level output.
func Test5Output(t *testing.T) {
	output := &bytes.Buffer{}
	l, err := New(1, output)

	if assert.Nil(t, err) {
		l.Info("xxx")
		assert.Contains(t, output.String(), "xxx")
		l.Debug("zzz")
		assert.NotContains(t, output.String(), "zzz")
	}
}
