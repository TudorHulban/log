package log

func (*Logger) labelTrace() string {
	return logLevels[LevelTrace]
}

func (l *Logger) Trace(args ...any) {
	if Level(l.logLevel.Load()) > LevelTrace {
		return
	}

	l.logWithLabel(
		l.labelTrace(), args...,
	)
}

func (l *Logger) Tracef(format string, args ...any) {
	if Level(l.logLevel.Load()) > LevelTrace {
		return
	}

	l.logfWithLabel(
		l.labelTrace(), format, args...,
	)
}

func (l *Logger) Tracew(msg string, keysAndValues ...any) {
	if Level(l.logLevel.Load()) > LevelTrace {
		return
	}

	if (len(keysAndValues) & 1) != 0 {
		l.PrintMessage("panicw: odd number of key-value arguments")
	}

	l.logwWithLabel(
		l.labelTrace(), msg, keysAndValues...,
	)
}

func (l *Logger) TraceFast(args ...any) {
	if Level(l.logLevel.Load()) > LevelTrace {
		return
	}

	l.logWithLabelFast(
		l.labelTrace(),
		l.estimatedMessageSizeTrace,
		args...,
	)
}
