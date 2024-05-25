package log

import (
	"os"
	"testing"
)

func TestMix(t *testing.T) {
	l := NewLogger(
		LevelDEBUG,
		os.Stdout,
		true,
	)

	go l.Print("0")
	go l.Info("1")
	go l.Warn("2")
	l.Debug("3")

	// <-l.w.ChStop
}
