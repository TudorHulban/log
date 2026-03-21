package log

import (
	"errors"

	"github.com/tudorhulban/log/bytearena"
	"github.com/tudorhulban/log/timestamp"
)

type Level uint8

type Logger struct {
	ingestor    *bytearena.Ingestor
	fnTimestamp timestamp.Timestamp

	estimatedMessageSize uint32
	logLevel             uint8

	withCaller bool // for shorter form in case do not need caller file.
	withColor  bool
	withJSON   bool
}

type ParamsNewLogger struct {
	Ingestor      *bytearena.Ingestor
	WithTimestamp timestamp.Timestamp

	EstimatedMessageSize uint32
	LoggerLevel          Level

	WithCaller bool
	WithColor  bool
	WithJSON   bool
}

func NewLogger(params *ParamsNewLogger) (*Logger, error) {
	if params.Ingestor == nil {
		return nil,
			errors.New("nil ingestor")
	}

	result := Logger{
		logLevel: convertLevel(params.LoggerLevel),

		withCaller:  params.WithCaller,
		fnTimestamp: params.WithTimestamp,
		withColor:   params.WithColor,
		withJSON:    params.WithJSON,

		ingestor:             params.Ingestor,
		estimatedMessageSize: params.EstimatedMessageSize,
	}

	if result.estimatedMessageSize == 0 {
		result.estimatedMessageSize = MessageMediumSize
	}

	result.Printf(
		"created logger, level %v",
		logLevels[params.LoggerLevel],
	)

	return &result,
		nil
}

func (l *Logger) labelInfo() string {
	if l.withColor {
		return colorInfo(logLevels[LevelINFO])
	} else { //nolint:revive
		return logLevels[LevelINFO]
	}
}
