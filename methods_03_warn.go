package log

func (l *Logger) labelWarn() string {
	if l.withColor {
		return colorWarn(logLevels[LevelWarn])
	}

	return logLevels[LevelWarn]
}

func (l *Logger) Warn(args ...any) {
	if Level(l.logLevel.Load()) > LevelWarn {
		return
	}

	l.logWithLabel(
		l.labelWarn(), args...,
	)
}

func (l *Logger) Warnf(format string, args ...any) {
	if Level(l.logLevel.Load()) > LevelWarn {
		return
	}

	l.logfWithLabel(
		l.labelWarn(), format, args,
	)
}

func (l *Logger) Warnw(msg string, keysAndValues ...any) {
	if Level(l.logLevel.Load()) > LevelWarn {
		return
	}

	l.logwWithLabel(
		l.labelWarn(), msg, keysAndValues...,
	)
}

func (l *Logger) WarnFast(args ...any) {
	if Level(l.logLevel.Load()) > LevelWarn {
		return
	}

	l.logWithLabelFast(
		l.labelWarn(),
		l.estimatedMessageSizeWarn,
		args,
	)
}
