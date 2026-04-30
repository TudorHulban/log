package log

import "github.com/tudorhulban/log/helpers"

func (l *Logger) labelDebug() string {
	if l.withColor {
		return colorDebug(logLevels[LevelDebug])
	}

	return logLevels[LevelDebug]
}

func (l *Logger) Debug(args ...any) {
	if Level(l.logLevel.Load()) > LevelDebug {
		return
	}

	estimatedMessageSizeInfo := helpers.GetEstimatedMessageSize("", args)

	l.logWithLabel(
		l.labelDebug(),
		uint32(estimatedMessageSizeInfo+_DeltaEstimation),
		args,
	)
}

func (l *Logger) Debugf(format string, args ...any) {
	if Level(l.logLevel.Load()) > LevelDebug {
		return
	}

	estimatedMessageSizeInfo := helpers.GetEstimatedMessageSize(format, args)

	l.logfWithLabel(
		l.labelDebug(),
		format,
		uint32(estimatedMessageSizeInfo),
		args,
	)
}

func (l *Logger) Debugw(msg string, keysAndValues ...any) {
	if Level(l.logLevel.Load()) > LevelDebug {
		return
	}

	if (len(keysAndValues) & 1) != 0 {
		l.PrintMessage("panicw: odd number of key-value arguments")
	}

	l.logwWithLabel(
		l.labelDebug(), msg, keysAndValues...,
	)
}

func (l *Logger) DebugFast(args ...any) {
	if Level(l.logLevel.Load()) > LevelDebug {
		return
	}

	l.logWithLabel(
		l.labelDebug(),
		l.estimatedMessageSizeDebug,
		args,
	)
}
