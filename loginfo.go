package loginfo

import (
	"fmt"
	"log"
	"os"
)

// Custom logger with nada, info and debug.
type Info struct {
	l        *log.Logger
	logLevel int // 0 - nada, 1 - info, 2 - debug
}

var logLevels = []string{"NADA", "INFO", "DEBUG"}
var delim = ": "

func (i Info) Infof(format string, args ...interface{}) {
	if i.logLevel > 0 {
		i.l.Output(2, logLevels[1]+delim+fmt.Sprintf(format, args...))
	}
}

func (i Info) Info(args ...interface{}) {
	if i.logLevel > 0 {
		i.l.Output(2, logLevels[1]+delim+fmt.Sprint(args...))
	}
}

func (i Info) Debugf(format string, args ...interface{}) {
	if i.logLevel > 1 {
		i.l.Output(2, logLevels[2]+delim+fmt.Sprintf(format, args...))
	}
}

func (i Info) Debug(args ...interface{}) {
	if i.logLevel > 1 {
		i.l.Output(2, logLevels[2]+delim+fmt.Sprint(args...))
	}
}

var thelogger LogInfo // should be type of interface.

func New(level int) (LogInfo, error) {
	if thelogger != nil {
		return thelogger, nil
	}
	lev := convertLevel(level)
	thelogger = &Info{
		l:        log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile),
		logLevel: lev,
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
