package log

import (
	"fmt"
	"io"

	"github.com/tudorhulban/log/timestamp"
)

type Level int8

type Logger struct {
	localWriter io.Writer
	fnTimestamp timestamp.Timestamp

	logLevel int8

	withCaller bool // for shorter form in case do not need caller file.
	withColor  bool
	withJSON   bool
}

type ParamsNewLogger struct {
	LoggerWriter  io.Writer
	WithTimestamp timestamp.Timestamp

	LoggerLevel Level

	WithCaller bool
	WithColor  bool
	WithJSON   bool
}

func NewLogger(params *ParamsNewLogger) *Logger {
	result := Logger{
		logLevel: convertLevel(params.LoggerLevel),

		withCaller:  params.WithCaller,
		fnTimestamp: params.WithTimestamp,
		withColor:   params.WithColor,
		withJSON:    params.WithJSON,

		localWriter: params.LoggerWriter,
	}

	if params.LoggerWriter == nil {
		result.localWriter = io.Discard
	}

	if params.WithTimestamp == nil {
		result.fnTimestamp = timestamp.TimestampNil
	}

	result.Printf(
		"created logger, level %v",
		logLevels[params.LoggerLevel],
	)

	return &result
}

func (l *Logger) appendJSON(buf, timestamp []byte, level, format string, args ...any) []byte {
	buf = append(buf, `{"timestamp":"`...)
	buf = append(buf, timestamp...)
	buf = append(buf, `","level":"`...)
	buf = append(buf, level...)
	buf = append(buf, `","message":"`...)
	buf = fmt.Appendf(buf, format, args...)
	buf = append(buf, "\"}\n"...)

	return buf
}

func (l Logger) labelInfo() string {
	return ternary(
		l.withColor,

		colorInfo(logLevels[LevelINFO]),
		logLevels[LevelINFO],
	)
}
