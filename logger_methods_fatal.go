package log

import "os"

func (l *Logger) Fatal(args ...any) {
	l.Print(args...)

	os.Exit(1)
}
