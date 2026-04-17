package log

import (
	"errors"

	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/log/timestamp"
)

type Logger struct {
	ingestor    *bytearena.Ingestor
	fnTimestamp timestamp.Timestamp

	callerLevel int

	estimatedMessageSizeOverall uint32
	estimatedMessageSizeInfo    uint32
	estimatedMessageSizeWarn    uint32
	estimatedMessageSizeDebug   uint32
	estimatedMessageSizeError   uint32

	logLevel Level

	withCaller bool // for shorter form in case do not need caller file.
	withColor  bool
	withJSON   bool
}

type ParamsNewLogger struct {
	Ingestor      *bytearena.Ingestor
	WithTimestamp timestamp.Timestamp

	EstimatedMessageSizeOverall uint32
	EstimatedMessageSizeInfo    uint32
	EstimatedMessageSizeWarn    uint32
	EstimatedMessageSizeDebug   uint32
	EstimatedMessageSizeError   uint32
	LoggerLevel                 Level

	CallerLevel uint8

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
		callerLevel: int(params.CallerLevel),

		fnTimestamp: params.WithTimestamp,
		withColor:   params.WithColor,
		withJSON:    params.WithJSON,

		ingestor: params.Ingestor,

		estimatedMessageSizeOverall: params.EstimatedMessageSizeOverall,
		estimatedMessageSizeDebug:   params.EstimatedMessageSizeDebug,
	}

	if result.estimatedMessageSizeOverall == 0 {
		result.estimatedMessageSizeOverall = MessageSmallSize
	}

	setEstimatedIfZero := func(value *uint32) {
		if *value == 0 {
			*value = result.estimatedMessageSizeOverall
		}
	}

	setEstimatedIfZero(&result.estimatedMessageSizeInfo)
	setEstimatedIfZero(&result.estimatedMessageSizeWarn)
	setEstimatedIfZero(&result.estimatedMessageSizeDebug)
	setEstimatedIfZero(&result.estimatedMessageSizeError)

	if result.callerLevel == 0 {
		result.callerLevel = 1
	}

	result.PrintfFast(
		"created logger, level %v",
		logLevels[params.LoggerLevel],
	)

	return &result,
		nil
}
