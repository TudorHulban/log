package log

func (l Logger) labelWarn() string {
	if l.withColor {
		return colorWarn(logLevels[LevelWARN])
	}

	return logLevels[LevelWARN]
}

func (l *Logger) Warn(args ...any) {
	if Level(l.logLevel.Load()) > LevelWARN {
		return
	}

	l.logWithLabel(
		l.labelWarn(), args...,
	)
}

func (l *Logger) Warnf(format string, args ...any) {
	if Level(l.logLevel.Load()) > LevelWARN {
		return
	}

	l.logfWithLabel(
		l.labelWarn(), format, args...,
	)
}

func (l *Logger) Warnw(msg string, keysAndValues ...any) {
	if Level(l.logLevel.Load()) > LevelWARN {
		return
	}

	l.logwWithLabel(
		l.labelWarn(), msg, keysAndValues...,
	)
}

func (l *Logger) WarnFast(args ...any) {
	if Level(l.logLevel.Load()) > LevelWARN {
		return
	}

	l.logWithLabelFast(
		l.labelWarn(),
		l.estimatedMessageSizeWarn,
		args...,
	)
}
