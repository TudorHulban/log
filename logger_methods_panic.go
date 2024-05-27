package log

import (
	"fmt"
	"os"
)

func (l Logger) Panic(args ...any) {
	panic(
		fmt.Sprint(args...),
	)
}

func (l Logger) Panicf(format string, args ...any) {
	panic(
		fmt.Sprintf(format, args...),
	)
}

func (l Logger) Panicw(msg string, keysAndValues ...any) {
	l.Printw(
		msg,
		keysAndValues...,
	)

	os.Exit(1)
}
