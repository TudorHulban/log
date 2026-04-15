package log

func (l *Logger) GetLogLevel() Level {
	return l.logLevel
}

func (l *Logger) SetLogLevel(level Level) {
	l.logLevel = level
}
