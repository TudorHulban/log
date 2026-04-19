package log

func (l *Logger) labelInfo() string {
	if l.withColor {
		return colorInfo(logLevels[LevelINFO])
	}

	return logLevels[LevelINFO]
}

func (l *Logger) Info(args ...any) {
	if l.logLevel > LevelINFO {
		return
	}

	l.logWithLabel(
		l.labelInfo(), args...,
	)
}

func (l *Logger) Infof(format string, args ...any) {
	if l.logLevel > LevelINFO {
		return
	}

	l.logfWithLabel(
		l.labelInfo(), format, args...,
	)
}

func (l *Logger) Infow(msg string, keysAndValues ...any) {
	if l.logLevel > LevelINFO {
		return
	}

	if (len(keysAndValues) & 1) != 0 {
		l.PrintMessage("panicw: odd number of key-value arguments")
	}

	l.logwWithLabel(
		l.labelInfo(), msg, keysAndValues...,
	)
}

func (l *Logger) InfoFast(args ...any) {
	if l.logLevel > LevelINFO {
		return
	}

	l.logWithLabelFast(
		l.labelInfo(),
		l.estimatedMessageSizeInfo,
		args...,
	)
}
