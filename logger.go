package log

import (
	"io"

	"github.com/TudorHulban/log/timestamp"
)

type Level int8

type Logger struct {
	localWriter io.Writer
	logLevel    int8

	withTimestamp timestamp.Timestamp
	withCaller    bool // for shorter form in case do not need caller file.
}

type ParamsNewLogger struct {
	LoggerLevel  Level
	LoggerWriter io.Writer

	WithTimestamp timestamp.Timestamp
	WithCaller    bool
}

func NewLogger(params *ParamsNewLogger) Logger {
	result := Logger{
		logLevel: convertLevel(params.LoggerLevel),

		withCaller:    params.WithCaller,
		withTimestamp: params.WithTimestamp,

		localWriter: params.LoggerWriter,
	}

	if params.LoggerWriter == nil {
		result.localWriter = io.Discard
	}

	result.Printf(
		"created logger, level %v.",
		logLevels[params.LoggerLevel],
	)

	return result
}
