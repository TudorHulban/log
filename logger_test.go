package log

import "testing"

func TestLogger(t *testing.T) {
	l := NewLogger(
		LevelDEBUG,
		nil,
		false,
	)

	go l.PrintMessage("xxxx")

	l.PrintMessage("xxxx")
}
