package loginfo

// LogInfo is the type supported when injecting the logger.
type LogInfo interface {
	Infof(string, ...interface{})
	Info(...interface{})
	Debugf(string, ...interface{})
	Debug(...interface{})
}
