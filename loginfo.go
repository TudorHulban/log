package loginfo

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"time"
)

const delim = ": "

var logLevels = []string{"NADA", "INFO", "WARN", "DEBUG"}

// LogInfo Custom logger with nada, info and debug.
type LogInfo struct {
	logLevel int
	writeTo  io.Writer
}

// Print Method prints without formatting.
func (i *LogInfo) Print(args ...interface{}) {
	var buf bytes.Buffer
	buf.WriteString(timestamp() + " " + fmt.Sprint(args...) + "\n")
	i.writeTo.Write(buf.Bytes())
}

// Printf Method prints without formatting.
func (i *LogInfo) Printf(format string, args ...interface{}) {
	var buf bytes.Buffer
	buf.WriteString(timestamp() + " " + fmt.Sprintf(format, args...) + "\n")
	i.writeTo.Write(buf.Bytes())
}

// Info Method logging info level without formatting.
func (i *LogInfo) Info(args ...interface{}) {
	if i.logLevel > 0 {
		_, file, line, _ := runtime.Caller(1)

		var buf bytes.Buffer
		buf.WriteString(timestamp() + " " + file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + logLevels[1] + delim + fmt.Sprint(args...) + "\n")
		i.writeTo.Write(buf.Bytes())
	}
}

// Infof Method logging info level with formatting.
func (i *LogInfo) Infof(format string, args ...interface{}) {
	if i.logLevel > 0 {
		_, file, line, _ := runtime.Caller(1)

		var buf bytes.Buffer
		buf.WriteString(timestamp() + " " + file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + logLevels[1] + delim + fmt.Sprintf(format, args...) + "\n")
		i.writeTo.Write(buf.Bytes())
	}
}

// Warn Method logging warn level without formatting.
func (i *LogInfo) Warn(args ...interface{}) {
	if i.logLevel > 1 {
		_, file, line, _ := runtime.Caller(1)

		var buf bytes.Buffer
		buf.WriteString(timestamp() + " " + file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + cWarn()(logLevels[2]) + delim + fmt.Sprint(args...) + "\n")
		i.writeTo.Write(buf.Bytes())
	}
}

// Warnf Method logging warn level with formatting.
func (i *LogInfo) Warnf(format string, args ...interface{}) {
	if i.logLevel > 1 {
		_, file, line, _ := runtime.Caller(1)

		var buf bytes.Buffer
		buf.WriteString(timestamp() + " " + file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + cWarn()(logLevels[2]) + delim + fmt.Sprintf(format, args...) + "\n")
		i.writeTo.Write(buf.Bytes())
	}
}

// Debug Method logging info debug level without formatting.
func (i *LogInfo) Debug(args ...interface{}) {
	if i.logLevel > 2 {
		_, file, line, _ := runtime.Caller(1)

		var buf bytes.Buffer
		buf.WriteString(timestamp() + " " + file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + cDebug()(logLevels[3]) + delim + fmt.Sprint(args...) + "\n")
		i.writeTo.Write(buf.Bytes())
	}
}

// Debugf Method logging info debug level with formatting.
func (i *LogInfo) Debugf(format string, args ...interface{}) {
	if i.logLevel > 2 {
		_, file, line, _ := runtime.Caller(1)

		var buf bytes.Buffer
		buf.WriteString(timestamp() + " " + file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + cDebug()(logLevels[3]) + delim + fmt.Sprintf(format, args...) + "\n")
		i.writeTo.Write(buf.Bytes())
	}
}

// GetLogLevel Method returns current log level.
func (i *LogInfo) GetLogLevel() int {
	return i.logLevel
}

// SetLogLevel Method sets log level.
func (i *LogInfo) SetLogLevel(level int) {
	i.logLevel = level
}

// New Constructor with levels 0 - nada, 1 - info, 2 - warn, 3 - debug.
func New(level int, writeTo io.Writer) (*LogInfo, error) {
	lev := convertLevel(level)
	result := LogInfo{
		logLevel: lev,
		writeTo:  writeTo,
	}
	result.Printf("Created logger, level %v.", logLevels[lev])
	return &result, nil
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
