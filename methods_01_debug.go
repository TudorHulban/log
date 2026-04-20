package log

func (l *Logger) labelDebug() string {
	if l.withColor {
		return colorDebug(logLevels[LevelDEBUG])
	}

	return logLevels[LevelDEBUG]
}

func (l *Logger) Debug(args ...any) {
	if Level(l.logLevel.Load()) > LevelDEBUG {
		return
	}

	l.logWithLabel(
		l.labelDebug(), args...,
	)
}

func (l *Logger) Debugf(format string, args ...any) {
	if Level(l.logLevel.Load()) > LevelDEBUG {
		return
	}

	l.logfWithLabel(
		l.labelDebug(), format, args...,
	)
}

func (l *Logger) Debugw(msg string, keysAndValues ...any) {
	if Level(l.logLevel.Load()) > LevelDEBUG {
		return
	}

	if (len(keysAndValues) & 1) != 0 {
		l.PrintMessage("panicw: odd number of key-value arguments")
	}

	l.logwWithLabel(
		l.labelDebug(), msg, keysAndValues...,
	)
}

func (l *Logger) DebugFast(args ...any) {
	if Level(l.logLevel.Load()) > LevelDEBUG {
		return
	}

	l.logWithLabelFast(
		l.labelDebug(),
		l.estimatedMessageSizeDebug,
		args...,
	)
}
