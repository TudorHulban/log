package log

func (l Logger) Trace(args ...any) {
	l.Print(args...)
}

func (l Logger) Tracef(format string, args ...any) {
	l.Printf(format, args...)
}

func (l Logger) Tracew(msg string, keysAndValues ...any) {
	l.Printw(
		msg,
		keysAndValues...,
	)
}
