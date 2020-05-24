package loginfo

import (
	"bytes"
	"fmt"
	"io"
	"log"
)

const delim = ": "

var logLevels = []string{"NADA", "INFO", "DEBUG"}

// LogInfo Custom logger with nada, info and debug.
type LogInfo struct {
	logLevel int
	writeTo  *bytes.Buffer
	l        *log.Logger
}

// Infof Method logging info level with formatting.
func (i LogInfo) Infof(format string, args ...interface{}) {
	if i.logLevel > 0 {
		i.l.Output(3, logLevels[1]+delim+fmt.Sprintf(format, args...))
	}
}

// Info Method logging info level without formatting.
func (i LogInfo) Info(args ...interface{}) {
	if i.logLevel > 0 {
		i.l.Output(3, logLevels[1]+delim+fmt.Sprint(args...))
	}
}

// Debugf Method logging info debug level with formatting.
func (i LogInfo) Debugf(format string, args ...interface{}) {
	if i.logLevel > 1 {
		i.l.Output(2, logLevels[2]+delim+fmt.Sprintf(format, args...))
	}
}

// Debug Method logging info debug level without formatting.
func (i LogInfo) Debug(args ...interface{}) {
	if i.logLevel > 1 {
		i.l.Output(2, logLevels[2]+delim+fmt.Sprint(args...))
	}
}

var thelogger LogInfo // should be type of interface.

// New flyweight constructor for logger.
// 0 - nada, 1 - info, 2 - debug
func New(level int, writeTo io.Writer) (LogInfo, error) {
	if thelogger.l != nil {
		return thelogger, nil
	}
	lev := convertLevel(level)
	thelogger = LogInfo{
		logLevel: lev,
		l:        log.New(writeTo, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile),
	}
	log.Printf("Created logger, level %v.", logLevels[lev])
	return thelogger, nil
}

func convertLevel(level int) int {
	switch {
	case level < 1:
		return 0
	case level == 1:
		return 1
	}
	return 2
}
