package log

func (l Logger) labelError() string {
	if l.withColor {
		return colorError(logLevels[LevelERROR])
	}

	return logLevels[LevelERROR]
}

func (l *Logger) Error(args ...any) {
	if l.logLevel < LevelDEBUG {
		return
	}

	l.logWithLabel(
		l.labelError(), args...,
	)
}

func (l *Logger) Errorf(format string, args ...any) {
	if l.logLevel < LevelDEBUG {
		return
	}

	l.logfWithLabel(
		l.labelError(), format, args...,
	)
}

func (l *Logger) Errorw(msg string, keysAndValues ...any) {
	if l.logLevel < LevelDEBUG {
		return
	}

	l.logwWithLabel(
		l.labelError(), msg, keysAndValues...,
	)
}

func (l *Logger) ErrorFast(args ...any) {
	if l.logLevel < LevelDEBUG {
		return
	}

	l.logWithLabelFast(
		l.labelError(),
		l.estimatedMessageSizeError,
		args...,
	)
}
