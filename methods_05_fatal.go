package log

import "os"

func (l Logger) Fatal(args ...any) {
	l.Print(args...)

	os.Exit(1)
}

func (l Logger) Fatalf(format string, args ...any) {
	l.Printf(format, args...)

	os.Exit(1)
}

func (l Logger) Fatalw(msg string, keysAndValues ...any) {
	l.Printw(
		msg,
		keysAndValues...,
	)

	os.Exit(1)
}
