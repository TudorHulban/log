package log

import (
	"errors"
	"fmt"

	"github.com/tudorhulban/log/bytearena"
	"github.com/tudorhulban/log/timestamp"
)

type Level int8

type Logger struct {
	ingestor    *bytearena.Ingestor
	fnTimestamp timestamp.Timestamp

	logLevel int8

	withCaller bool // for shorter form in case do not need caller file.
	withColor  bool
	withJSON   bool
}

type ParamsNewLogger struct {
	Ingestor      *bytearena.Ingestor
	WithTimestamp timestamp.Timestamp

	LoggerLevel Level

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

		ingestor: params.Ingestor,
	}

	result.Printf(
		"created logger, level %v",
		logLevels[params.LoggerLevel],
	)

	return &result,
		nil
}

func (*Logger) appendJSON(buf, ts []byte, level, format string, args ...any) []byte {
	buf = append(buf, `{"timestamp":"`...)
	buf = append(buf, ts...)
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
