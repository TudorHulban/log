package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"time"
)

// Not using defer for release to pool for performance reasons.

const delim = ": "

// helpers for constructor
const (
	// NONE no logging
	NONE = 0
	// INFO level
	INFO = 1
	// WARN level
	WARN = 2
	// DEBUG level
	DEBUG = 3
)

var logLevels = []string{"NONE", "INFO", "WARN", "DEBUG"}

// Logger Custom logger with nada, info and debug.
type Logger struct {
	logLevel int
	writeTo  io.Writer
	// for shorter form in case do not need caller file.
	withCaller bool
}

// New Constructor with levels 0 - nada, 1 - info, 2 - warn, 3 - debug.
func NewLogger(level int, writeTo io.Writer, withCaller bool) *Logger {
	lev := convertLevel(level)

	if writeTo == nil {
		writeTo = os.Stdout
	}

	result := Logger{
		logLevel:   lev,
		writeTo:    writeTo,
		withCaller: withCaller,
	}

	result.Printf("Created logger, level %v.", logLevels[lev])
	return &result
}

// Print Method prints without formatting.
func (i *Logger) Print(args ...interface{}) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()

	buf.WriteString(timestamp() + " " + fmt.Sprint(args...) + "\n")
	i.writeTo.Write(buf.Bytes())

	bufPool.Put(buf)
}

// Printf Method prints without formatting.
func (i *Logger) Printf(format string, args ...interface{}) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()

	buf.WriteString(timestamp() + " " + fmt.Sprintf(format, args...) + "\n")
	i.writeTo.Write(buf.Bytes())

	bufPool.Put(buf)
}

// Info Method logging info level without formatting.
func (i *Logger) Info(args ...interface{}) {
	if i.logLevel > 0 {
		buf := bufPool.Get().(*bytes.Buffer)
		buf.Reset()

		if i.withCaller {
			_, file, line, _ := runtime.Caller(1)
			buf.WriteString(timestamp() + " " + file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + logLevels[1] + delim + fmt.Sprint(args...) + "\n")
		} else {
			buf.WriteString(timestamp() + " " + logLevels[1] + delim + fmt.Sprint(args...) + "\n")
		}

		i.writeTo.Write(buf.Bytes())
		bufPool.Put(buf)
	}
}

// Infof Method logging info level with formatting.
func (i *Logger) Infof(format string, args ...interface{}) {
	if i.logLevel > 0 {
		buf := bufPool.Get().(*bytes.Buffer)
		buf.Reset()

		if i.withCaller {
			_, file, line, _ := runtime.Caller(1)
			buf.WriteString(timestamp() + " " + file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + logLevels[1] + delim + fmt.Sprintf(format, args...) + "\n")
		} else {
			buf.WriteString(timestamp() + " " + logLevels[1] + delim + fmt.Sprintf(format, args...) + "\n")
		}

		i.writeTo.Write(buf.Bytes())
		bufPool.Put(buf)
	}
}

// Warn Method logging warn level without formatting.
func (i *Logger) Warn(args ...interface{}) {
	if i.logLevel > 1 {
		buf := bufPool.Get().(*bytes.Buffer)
		buf.Reset()

		if i.withCaller {
			_, file, line, _ := runtime.Caller(1)
			buf.WriteString(timestamp() + " " + file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + cWarn()(logLevels[2]) + delim + fmt.Sprint(args...) + "\n")
		} else {
			buf.WriteString(timestamp() + " " + cWarn()(logLevels[2]) + delim + fmt.Sprint(args...) + "\n")
		}

		i.writeTo.Write(buf.Bytes())
		bufPool.Put(buf)
	}
}

// Warnf Method logging warn level with formatting.
func (i *Logger) Warnf(format string, args ...interface{}) {
	if i.logLevel > 1 {
		buf := bufPool.Get().(*bytes.Buffer)
		buf.Reset()

		if i.withCaller {
			_, file, line, _ := runtime.Caller(1)
			buf.WriteString(timestamp() + " " + file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + cWarn()(logLevels[2]) + delim + fmt.Sprintf(format, args...) + "\n")
		} else {
			buf.WriteString(timestamp() + " " + cWarn()(logLevels[2]) + delim + fmt.Sprintf(format, args...) + "\n")
		}

		i.writeTo.Write(buf.Bytes())
		bufPool.Put(buf)
	}
}

// Debug Method logging info debug level without formatting.
func (i *Logger) Debug(args ...interface{}) {
	if i.logLevel > 2 {
		buf := bufPool.Get().(*bytes.Buffer)
		buf.Reset()

		if i.withCaller {
			_, file, line, _ := runtime.Caller(1)
			buf.WriteString(timestamp() + " " + file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + cDebug()(logLevels[3]) + delim + fmt.Sprint(args...) + "\n")
		} else {
			buf.WriteString(timestamp() + " " + cDebug()(logLevels[3]) + delim + fmt.Sprint(args...) + "\n")
		}

		i.writeTo.Write(buf.Bytes())
		bufPool.Put(buf)
	}
}

// Debugf Method logging info debug level with formatting.
func (i *Logger) Debugf(format string, args ...interface{}) {
	if i.logLevel > 2 {
		buf := bufPool.Get().(*bytes.Buffer)
		buf.Reset()

		if i.withCaller {
			_, file, line, _ := runtime.Caller(1)
			buf.WriteString(timestamp() + " " + file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + cDebug()(logLevels[3]) + delim + fmt.Sprintf(format, args...) + "\n")
		} else {
			buf.WriteString(timestamp() + " " + cDebug()(logLevels[3]) + delim + fmt.Sprintf(format, args...) + "\n")

		}

		i.writeTo.Write(buf.Bytes())
		bufPool.Put(buf)
	}
}

// GetLogLevel Method returns current log level.
func (i *Logger) GetLogLevel() int {
	return i.logLevel
}

// SetLogLevel Method sets log level.
func (i *Logger) SetLogLevel(level int) {
	i.logLevel = level
}

func convertLevel(level int) int {
	switch {
	case level < 1:
		return 0
	case level == 1:
		return 1
	case level == 2:
		return 2
	}

	return 3
}

// timestamp Helper provides time in format YYYYMonth HH24:Minutes:Seconds.Miliseconds
func timestamp() string {
	now := time.Now()
	theMonth := "0" + strconv.FormatInt(int64(now.Month()), 10)
	theMonth = theMonth[len(theMonth)-2:]

	theHour := "0" + strconv.FormatInt(int64(now.Hour()), 10)
	theHour = theHour[len(theHour)-2:]

	theMin := "0" + strconv.FormatInt(int64(now.Minute()), 10)
	theMin = theMin[len(theMin)-2:]

	theSec := "0" + strconv.FormatInt(int64(now.Second()), 10)
	theSec = theSec[len(theSec)-2:]

	theMilisec := "00" + strconv.FormatInt(int64(now.Nanosecond()/1000000), 10)
	theMilisec = theMilisec[len(theMilisec)-3:]

	return strconv.FormatInt(int64(now.Year()), 10) + theMonth + " " + theHour + ":" + theMin + ":" + theSec + "." + theMilisec
}
