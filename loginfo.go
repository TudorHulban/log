package loginfo

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"runtime"
	"strconv"
)

const delim = ": "

var logLevels = []string{"NADA", "INFO", "DEBUG"}

// LogInfo Custom logger with nada, info and debug.
type LogInfo struct {
	logLevel int
	writeTo  io.Writer
	l        *log.Logger
}

// Info Method logging info level without formatting.
func (i LogInfo) Info(args ...interface{}) {
	if i.logLevel > 0 {
		_, file, line, _ := runtime.Caller(1)

		var buf bytes.Buffer
		buf.WriteString(file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + logLevels[1] + delim + fmt.Sprint(args...) + "\n")
		i.writeTo.Write(buf.Bytes())
	}
}

// Infof Method logging info level with formatting.
func (i LogInfo) Infof(format string, args ...interface{}) {
	if i.logLevel > 0 {
		_, file, line, _ := runtime.Caller(1)

		var buf bytes.Buffer
		buf.WriteString(file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + logLevels[1] + delim + fmt.Sprintf(format, args...) + "\n")
		i.writeTo.Write(buf.Bytes())
	}
}

// Debug Method logging info debug level without formatting.
func (i LogInfo) Debug(args ...interface{}) {
	if i.logLevel > 1 {
		_, file, line, _ := runtime.Caller(1)

		var buf bytes.Buffer
		buf.WriteString(file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + logLevels[2] + delim + fmt.Sprint(args...) + "\n")
		i.writeTo.Write(buf.Bytes())
	}
}

// Debugf Method logging info debug level with formatting.
func (i LogInfo) Debugf(format string, args ...interface{}) {
	if i.logLevel > 1 {
		_, file, line, _ := runtime.Caller(1)

		var buf bytes.Buffer
		buf.WriteString(file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + logLevels[2] + delim + fmt.Sprintf(format, args...) + "\n")
		i.writeTo.Write(buf.Bytes())
	}
}

// GetLogLevel Method returns current log level.
func (i LogInfo) GetLogLevel() int {
	return i.logLevel
}

// SetLogLevel Method sets log level.
func (i LogInfo) SetLogLevel(level int) {
	i.logLevel = level
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
		l:        log.New(writeTo, "", log.LstdFlags),
		writeTo:  writeTo,
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
