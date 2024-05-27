package log

import "os"

func (l Logger) Panic(args ...any) {
	l.Print(args...)

	os.Exit(1)
}

func (l Logger) Panicf(format string, args ...any) {
	l.Printf(format, args...)

	os.Exit(1)
}

func (l Logger) Panicw(msg string, keysAndValues ...any) {
	l.Printw(
		msg,
		keysAndValues...,
	)

	os.Exit(1)
}
