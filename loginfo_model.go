package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
)

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

type Logger struct {
	logLevel int
	writeTo  io.Writer

	// for shorter form in case do not need caller file.
	withCaller bool
}

// New Constructor with levels 0 - no logging, 1 - info, 2 - warn, 3 - debug.
func NewLogger(level int, writeTo io.Writer, withCaller bool) *Logger {
	lev := convertLevel(level)

	if writeTo == nil {
		writeTo = os.Stdout
	}

	res := Logger{
		logLevel:   lev,
		writeTo:    writeTo,
		withCaller: withCaller,
	}

	res.Printf("Created logger, level %v.", logLevels[lev])
	return &res
}

func (i *Logger) Print(args ...interface{}) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()

	buf.WriteString(timestamp() + " " + fmt.Sprint(args...) + "\n")
	i.writeTo.Write(buf.Bytes())

	bufPool.Put(buf)
}

func (i *Logger) Printf(format string, args ...interface{}) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()

	buf.WriteString(timestamp() + " " + fmt.Sprintf(format, args...) + "\n")
	i.writeTo.Write(buf.Bytes())

	bufPool.Put(buf)
}

func (i *Logger) Info(args ...interface{}) {
	if i.logLevel == 0 {
		return
	}

	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)

	buf.Reset()

	if i.withCaller {
		_, file, line, _ := runtime.Caller(1)
		buf.WriteString(timestamp() + " " + file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + logLevels[1] + delim + fmt.Sprint(args...) + "\n")
		i.writeTo.Write(buf.Bytes())
		return
	}

	buf.WriteString(timestamp() + " " + logLevels[1] + delim + fmt.Sprint(args...) + "\n")
	i.writeTo.Write(buf.Bytes())
}

func (i *Logger) Infof(format string, args ...interface{}) {
	if i.logLevel == 0 {
		return
	}

	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)

	buf.Reset()

	if i.withCaller {
		_, file, line, _ := runtime.Caller(1)
		buf.WriteString(timestamp() + " " + file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + logLevels[1] + delim + fmt.Sprintf(format, args...) + "\n")
		i.writeTo.Write(buf.Bytes())
		return
	}

	buf.WriteString(timestamp() + " " + logLevels[1] + delim + fmt.Sprintf(format, args...) + "\n")
	i.writeTo.Write(buf.Bytes())
}

func (i *Logger) Warn(args ...interface{}) {
	if i.logLevel < 2 {
		return
	}

	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)

	buf.Reset()

	if i.withCaller {
		_, file, line, _ := runtime.Caller(1)
		buf.WriteString(timestamp() + " " + file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + colorWarn()(logLevels[2]) + delim + fmt.Sprint(args...) + "\n")
		i.writeTo.Write(buf.Bytes())
		return
	}

	buf.WriteString(timestamp() + " " + colorWarn()(logLevels[2]) + delim + fmt.Sprint(args...) + "\n")
	i.writeTo.Write(buf.Bytes())
}

func (i *Logger) Warnf(format string, args ...interface{}) {
	if i.logLevel < 2 {
		return
	}

	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)

	buf.Reset()

	if i.withCaller {
		_, file, line, _ := runtime.Caller(1)
		buf.WriteString(timestamp() + " " + file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + colorWarn()(logLevels[2]) + delim + fmt.Sprintf(format, args...) + "\n")
		i.writeTo.Write(buf.Bytes())
		return
	}

	buf.WriteString(timestamp() + " " + colorWarn()(logLevels[2]) + delim + fmt.Sprintf(format, args...) + "\n")
	i.writeTo.Write(buf.Bytes())
}

func (i *Logger) Debug(args ...interface{}) {
	if i.logLevel < 3 {
		return
	}

	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)

	buf.Reset()

	if i.withCaller {
		_, file, line, _ := runtime.Caller(1)
		buf.WriteString(timestamp() + " " + file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + colorDebug()(logLevels[3]) + delim + fmt.Sprint(args...) + "\n")
		i.writeTo.Write(buf.Bytes())
		return
	}

	buf.WriteString(timestamp() + " " + colorDebug()(logLevels[3]) + delim + fmt.Sprint(args...) + "\n")
	i.writeTo.Write(buf.Bytes())
}

func (i *Logger) Debugf(format string, args ...interface{}) {
	if i.logLevel < 3 {
		return
	}

	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)

	buf.Reset()

	if i.withCaller {
		_, file, line, _ := runtime.Caller(1)
		buf.WriteString(timestamp() + " " + file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + colorDebug()(logLevels[3]) + delim + fmt.Sprintf(format, args...) + "\n")
		i.writeTo.Write(buf.Bytes())
		return
	}

	buf.WriteString(timestamp() + " " + colorDebug()(logLevels[3]) + delim + fmt.Sprintf(format, args...) + "\n")
	i.writeTo.Write(buf.Bytes())
}

func (i *Logger) Fatal(args ...interface{}) {
	i.Print(args...)

	os.Exit(1)
}

func (i *Logger) GetLogLevel() int {
	return i.logLevel
}

func (i *Logger) SetLogLevel(level int) {
	i.logLevel = level
}
