package log

func (l Logger) labelDebug() string {
	if l.withColor {
		return colorDebug(logLevels[LevelDEBUG])
	}

	return logLevels[LevelDEBUG]
}

func (l *Logger) Debug(args ...any) {
	l.logWithLabel(
		l.labelDebug(), args...,
	)
}

func (l *Logger) Debugf(format string, args ...any) {
	l.logfWithLabel(
		l.labelDebug(), format, args...,
	)
}

func (l *Logger) Debugw(msg string, keysAndValues ...any) {
	l.logwWithLabel(
		l.labelDebug(), msg, keysAndValues...,
	)
}

func (l *Logger) DebugFast(args ...any) {
	l.logWithLabelFast(
		l.labelDebug(), args...,
	)
}
