package loginfo

// ILogInfo is the type supported when injecting the logger.
type ILogInfo interface {
	Infof(string, ...interface{})
	Info(...interface{})
	Debugf(string, ...interface{})
	Debug(...interface{})
}
