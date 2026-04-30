package log

import "github.com/tudorhulban/log/helpers"

func (*Logger) labelTrace() string {
	return logLevels[LevelTrace]
}

func (l *Logger) Trace(args ...any) {
	if Level(l.logLevel.Load()) > LevelTrace {
		return
	}

	estimatedMessageSizeInfo := helpers.GetEstimatedMessageSize("", args)

	l.logWithLabel(
		l.labelTrace(),
		uint32(estimatedMessageSizeInfo),
		args,
	)
}

func (l *Logger) Tracef(format string, args ...any) {
	if Level(l.logLevel.Load()) > LevelTrace {
		return
	}

	estimatedMessageSizeInfo := helpers.GetEstimatedMessageSize(format, args)

	l.logfWithLabel(
		l.labelTrace(),
		format,
		uint32(estimatedMessageSizeInfo),
		args,
	)
}

func (l *Logger) Tracew(msg string, keysAndValues ...any) {
	if Level(l.logLevel.Load()) > LevelTrace {
		return
	}

	if (len(keysAndValues) & 1) != 0 {
		l.PrintMessage("panicw: odd number of key-value arguments")
	}

	l.logwWithLabel(
		l.labelTrace(), msg, keysAndValues...,
	)
}

func (l *Logger) TraceFast(args ...any) {
	if Level(l.logLevel.Load()) > LevelTrace {
		return
	}

	l.logWithLabel(
		l.labelTrace(),
		l.estimatedMessageSizeTrace,
		args,
	)
}
