package log

func (l *Logger) labelError() string {
	if l.withColor {
		return colorError(logLevels[LevelError])
	}

	return logLevels[LevelError]
}

func (l *Logger) Error(args ...any) {
	if Level(l.logLevel.Load()) > LevelError {
		return
	}

	l.logWithLabel(
		l.labelError(), args...,
	)
}

func (l *Logger) Errorf(format string, args ...any) {
	if Level(l.logLevel.Load()) > LevelError {
		return
	}

	l.logfWithLabel(
		l.labelError(), format, args...,
	)
}

func (l *Logger) Errorw(msg string, keysAndValues ...any) {
	if Level(l.logLevel.Load()) > LevelError {
		return
	}

	l.logwWithLabel(
		l.labelError(), msg, keysAndValues...,
	)
}

func (l *Logger) ErrorFast(args ...any) {
	if Level(l.logLevel.Load()) > LevelError {
		return
	}

	l.logWithLabelFast(
		l.labelError(),
		l.estimatedMessageSizeError,
		args...,
	)
}
