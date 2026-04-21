package log

import (
	"errors"
	"io"
	"sync/atomic"

	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/log/timestamp"
)

type Logger struct { //nolint:govet
	// Cache line 0 (hot, bytes 0..63):
	// - atomic level and per-level estimated sizes are read on every log call.
	// - booleans and callerLevel are also hot (formatting + caller decision).
	// - ingestor and fnTimestamp are hot-ish and placed here so a single cache load
	//   brings the common hot state into L1.
	logLevel                    atomic.Uint32 // 4
	estimatedMessageSizeOverall uint32        // 4
	estimatedMessageSizeTrace   uint32        // 4
	estimatedMessageSizeDebug   uint32        // 4
	estimatedMessageSizeInfo    uint32        // 4
	estimatedMessageSizeWarn    uint32        // 4
	estimatedMessageSizeError   uint32        // 4

	withCaller  bool  // 1
	withColor   bool  // 1
	withJSON    bool  // 1
	callerLevel uint8 // 1

	ingestor    *bytearena.Ingestor // 8
	fnTimestamp timestamp.Timestamp // 8

	// pad to end of 64‑byte cache line
	_pad0 [16]byte //nolint:unused

	// Cache line 1 (cold/rare, bytes 64..127):
	// - fatalWriter is rare (fatal path). Keep it separate so it doesn't evict hot state.
	fatalWriter io.Writer // 16 (interface: type+data)
}

type ParamsNewLogger struct {
	Ingestor        *bytearena.Ingestor
	WithTimestamp   timestamp.Timestamp
	WithFatalWriter io.Writer

	EstimatedMessageSizeOverall uint32
	EstimatedMessageSizeTrace   uint32
	EstimatedMessageSizeDebug   uint32
	EstimatedMessageSizeInfo    uint32
	EstimatedMessageSizeWarn    uint32
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

	if params.WithFatalWriter == nil {
		return nil,
			errors.New("nil fatal writer")
	}

	result := Logger{
		withCaller:  params.WithCaller,
		callerLevel: params.CallerLevel,

		fnTimestamp: params.WithTimestamp,
		withColor:   params.WithColor,
		withJSON:    params.WithJSON,

		ingestor: params.Ingestor,

		estimatedMessageSizeOverall: params.EstimatedMessageSizeOverall,
		estimatedMessageSizeTrace:   params.EstimatedMessageSizeTrace,
		estimatedMessageSizeDebug:   params.EstimatedMessageSizeDebug,
		estimatedMessageSizeInfo:    params.EstimatedMessageSizeInfo,
		estimatedMessageSizeWarn:    params.EstimatedMessageSizeWarn,
		estimatedMessageSizeError:   params.EstimatedMessageSizeError,
	}

	result.SetLogLevel(params.LoggerLevel)

	if result.estimatedMessageSizeOverall == 0 {
		result.estimatedMessageSizeOverall = MessageSmallSize
	}

	setEstimatedIfZero := func(value *uint32) {
		if *value == 0 {
			*value = result.estimatedMessageSizeOverall
		}
	}

	setEstimatedIfZero(&result.estimatedMessageSizeTrace)
	setEstimatedIfZero(&result.estimatedMessageSizeDebug)
	setEstimatedIfZero(&result.estimatedMessageSizeInfo)
	setEstimatedIfZero(&result.estimatedMessageSizeWarn)
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
