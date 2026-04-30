package log

import "github.com/tudorhulban/log/helpers"

func (l *Logger) labelInfo() string {
	if l.withColor {
		return colorInfo(logLevels[LevelInfo])
	}

	return logLevels[LevelInfo]
}

func (l *Logger) Info(args ...any) {
	if Level(l.logLevel.Load()) > LevelInfo {
		return
	}

	estimatedMessageSizeInfo := helpers.GetEstimatedMessageSize("", args)

	l.logWithLabel(
		l.labelInfo(),
		uint32(estimatedMessageSizeInfo),
		args,
	)
}

func (l *Logger) Infof(format string, args ...any) {
	if Level(l.logLevel.Load()) > LevelInfo {
		return
	}

	estimatedMessageSizeInfo := helpers.GetEstimatedMessageSize(format, args)

	l.logfWithLabel(
		l.labelInfo(),
		format,
		uint32(estimatedMessageSizeInfo),
		args,
	)
}

func (l *Logger) Infow(msg string, keysAndValues ...any) {
	if Level(l.logLevel.Load()) > LevelInfo {
		return
	}

	if (len(keysAndValues) & 1) != 0 {
		l.PrintMessage("panicw: odd number of key-value arguments")
	}

	l.logwWithLabel(
		l.labelInfo(), msg, keysAndValues...,
	)
}

func (l *Logger) InfoFast(args ...any) {
	if Level(l.logLevel.Load()) > LevelInfo {
		return
	}

	l.logWithLabel(
		l.labelInfo(),
		l.estimatedMessageSizeInfo,
		args,
	)
}
