package log

func (l *Logger) GetLogLevel() int8 {
	return l.logLevel
}

func (l *Logger) SetLogLevel(level int8) {
	l.logLevel = level
}
