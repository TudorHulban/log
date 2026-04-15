package log

func (l Logger) labelDebug() string {
	if l.withColor {
		return colorDebug(logLevels[LevelDEBUG])
	}

	return logLevels[LevelDEBUG]
}

func (l *Logger) Debug(args ...any) {
	if l.logLevel > LevelDEBUG || l.logLevel == LevelNONE {
		return
	}

	l.logWithLabel(
		l.labelDebug(), args...,
	)
}

func (l *Logger) Debugf(format string, args ...any) {
	if l.logLevel > LevelDEBUG || l.logLevel == LevelNONE {
		return
	}

	l.logfWithLabel(
		l.labelDebug(), format, args...,
	)
}

func (l *Logger) Debugw(msg string, keysAndValues ...any) {
	if l.logLevel > LevelDEBUG || l.logLevel == LevelNONE {
		return
	}

	l.logwWithLabel(
		l.labelDebug(), msg, keysAndValues...,
	)
}

func (l *Logger) DebugFast(args ...any) {
	if l.logLevel > LevelDEBUG || l.logLevel == LevelNONE {
		return
	}

	l.logWithLabelFast(
		l.labelDebug(),
		l.estimatedMessageSizeDebug,
		args...,
	)
}
