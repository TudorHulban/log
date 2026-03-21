package log

func (l *Logger) GetLogLevel() uint8 {
	return l.logLevel
}

func (l *Logger) SetLogLevel(level uint8) {
	l.logLevel = level
}
