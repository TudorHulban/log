package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
)

type Logger struct {
	logLevel int
	writeTo  io.Writer

	// for shorter form in case do not need caller file.
	withCaller bool
}

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

	go res.Printf(
		"created logger, level %v.",
		logLevels[lev],
	)

	return &res
}

func (l *Logger) Print(args ...interface{}) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()

	buf.WriteString(timestamp() + " " + fmt.Sprint(args...) + "\n")

	l.writeTo.Write(buf.Bytes())

	bufPool.Put(buf)
}

func (l *Logger) Printf(format string, args ...interface{}) {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()

	buf.WriteString(timestamp() + " " + fmt.Sprintf(format, args...) + "\n")

	l.writeTo.Write(buf.Bytes())

	bufPool.Put(buf)
}

func (l *Logger) Info(args ...interface{}) {
	if l.logLevel == 0 {
		return
	}

	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)

	buf.Reset()

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		buf.WriteString(timestamp() + " " + file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + logLevels[1] + delim + fmt.Sprint(args...) + "\n")

		l.writeTo.Write(buf.Bytes())

		return
	}

	buf.WriteString(timestamp() + " " + logLevels[1] + delim + fmt.Sprint(args...) + "\n")

	l.writeTo.Write(buf.Bytes())
}

func (l *Logger) Infof(format string, args ...interface{}) {
	if l.logLevel == 0 {
		return
	}

	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)

	buf.Reset()

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		buf.WriteString(timestamp() + " " + file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + logLevels[1] + delim + fmt.Sprintf(format, args...) + "\n")

		l.writeTo.Write(buf.Bytes())

		return
	}

	buf.WriteString(timestamp() + " " + logLevels[1] + delim + fmt.Sprintf(format, args...) + "\n")
	l.writeTo.Write(buf.Bytes())
}

func (l *Logger) Warn(args ...interface{}) {
	if l.logLevel < 2 {
		return
	}

	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)

	buf.Reset()

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		buf.WriteString(timestamp() + " " + file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + colorWarn()(logLevels[2]) + delim + fmt.Sprint(args...) + "\n")

		l.writeTo.Write(buf.Bytes())

		return
	}

	buf.WriteString(timestamp() + " " + colorWarn()(logLevels[2]) + delim + fmt.Sprint(args...) + "\n")

	l.writeTo.Write(buf.Bytes())
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.logLevel < 2 {
		return
	}

	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)

	buf.Reset()

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		buf.WriteString(timestamp() + " " + file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + colorWarn()(logLevels[2]) + delim + fmt.Sprintf(format, args...) + "\n")

		l.writeTo.Write(buf.Bytes())

		return
	}

	buf.WriteString(timestamp() + " " + colorWarn()(logLevels[2]) + delim + fmt.Sprintf(format, args...) + "\n")

	l.writeTo.Write(buf.Bytes())
}

func (l *Logger) Debug(args ...interface{}) {
	if l.logLevel < 3 {
		return
	}

	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)

	buf.Reset()

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		buf.WriteString(timestamp() + " " + file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + colorDebug()(logLevels[3]) + delim + fmt.Sprint(args...) + "\n")

		l.writeTo.Write(buf.Bytes())

		return
	}

	buf.WriteString(timestamp() + " " + colorDebug()(logLevels[3]) + delim + fmt.Sprint(args...) + "\n")

	l.writeTo.Write(buf.Bytes())
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.logLevel < 3 {
		return
	}

	buf := bufPool.Get().(*bytes.Buffer)
	defer bufPool.Put(buf)

	buf.Reset()

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		buf.WriteString(timestamp() + " " + file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + colorDebug()(logLevels[3]) + delim + fmt.Sprintf(format, args...) + "\n")

		l.writeTo.Write(buf.Bytes())

		return
	}

	buf.WriteString(timestamp() + " " + colorDebug()(logLevels[3]) + delim + fmt.Sprintf(format, args...) + "\n")

	l.writeTo.Write(buf.Bytes())
}

func (l *Logger) Fatal(args ...interface{}) {
	l.Print(args...)

	os.Exit(1)
}

func (l *Logger) GetLogLevel() int {
	return l.logLevel
}

func (l *Logger) SetLogLevel(level int) {
	l.logLevel = level
}
