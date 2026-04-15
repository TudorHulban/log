package log

func (l Logger) labelWarn() string {
	if l.withColor {
		return colorDebug(logLevels[LevelWARN])
	}

	return logLevels[LevelWARN]
}

func (l *Logger) Warn(args ...any) {
	if l.logLevel > LevelWARN || l.logLevel == LevelNONE {
		return
	}

	l.logWithLabel(
		l.labelWarn(), args...,
	)
}

func (l *Logger) Warnf(format string, args ...any) {
	if l.logLevel > LevelWARN || l.logLevel == LevelNONE {
		return
	}

	l.logfWithLabel(
		l.labelWarn(), format, args...,
	)
}

func (l *Logger) Warnw(msg string, keysAndValues ...any) {
	if l.logLevel > LevelWARN || l.logLevel == LevelNONE {
		return
	}

	l.logwWithLabel(
		l.labelWarn(), msg, keysAndValues...,
	)
}

func (l *Logger) WarnFast(args ...any) {
	if l.logLevel > LevelWARN || l.logLevel == LevelNONE {
		return
	}

	l.logWithLabelFast(
		l.labelWarn(),
		l.estimatedMessageSizeWarn,
		args...,
	)
}
