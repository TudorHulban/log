package log

import (
	"io"

	"github.com/TudorHulban/log/timestamp"
)

type Level int8

type Logger struct {
	localWriter   io.Writer
	withTimestamp timestamp.Timestamp

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

		withCaller:    params.WithCaller,
		withTimestamp: params.WithTimestamp,
		withColor:     params.WithColor,
		withJSON:      params.WithJSON,

		localWriter: params.LoggerWriter,
	}

	if params.LoggerWriter == nil {
		result.localWriter = io.Discard
	}

	result.Printf(
		"created logger, level %v.",
		logLevels[params.LoggerLevel],
	)

	return &result
}
