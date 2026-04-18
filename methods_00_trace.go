package log

func (Logger) labelTrace() string {
	return logLevels[LevelTrace]
}

func (l Logger) Trace(args ...any) {
	if l.logLevel > LevelTrace {
		return
	}

	l.logWithLabel(
		l.labelTrace(), args...,
	)
}

func (l Logger) Tracef(format string, args ...any) {
	if l.logLevel > LevelTrace {
		return
	}

	l.logfWithLabel(
		l.labelTrace(), format, args...,
	)
}

func (l Logger) Tracew(msg string, keysAndValues ...any) {
	if l.logLevel > LevelTrace {
		return
	}

	l.logwWithLabel(
		l.labelTrace(), msg, keysAndValues...,
	)
}
