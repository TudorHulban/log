package log

func (l Logger) labelError() string {
	if l.withColor {
		return colorError(logLevels[LevelERROR])
	}

	return logLevels[LevelERROR]
}

func (l *Logger) Error(args ...any) {
	l.logWithLabel(
		l.labelError(), args...,
	)
}

func (l *Logger) Errorf(format string, args ...any) {
	l.logfWithLabel(
		l.labelError(), format, args...,
	)
}

func (l *Logger) Errorw(msg string, keysAndValues ...any) {
	l.logwWithLabel(
		l.labelError(), msg, keysAndValues...,
	)
}

func (l *Logger) ErrorFast(args ...any) {
	l.logWithLabelFast(
		l.labelError(), args...,
	)
}
