package log

import "github.com/tudorhulban/log/helpers"

func (l *Logger) labelError() string {
	if l.withColor {
		return colorError(logLevels[LevelError])
	}

	return logLevels[LevelError]
}

func (l *Logger) Error(args ...any) {
	if Level(l.logLevel.Load()) > LevelError {
		return
	}

	estimatedMessageSizeInfo := helpers.GetEstimatedMessageSize("", args)

	l.logWithLabel(
		l.labelError(),
		uint32(estimatedMessageSizeInfo),
		args,
	)
}

func (l *Logger) Errorf(format string, args ...any) {
	if Level(l.logLevel.Load()) > LevelError {
		return
	}

	estimatedMessageSizeInfo := helpers.GetEstimatedMessageSize(format, args)

	l.logfWithLabel(
		l.labelError(),
		format,
		uint32(estimatedMessageSizeInfo),
		args,
	)
}

func (l *Logger) Errorw(msg string, keysAndValues ...any) {
	if Level(l.logLevel.Load()) > LevelError {
		return
	}

	l.logwWithLabel(
		l.labelError(), msg, keysAndValues...,
	)
}

func (l *Logger) ErrorFast(args ...any) {
	if Level(l.logLevel.Load()) > LevelError {
		return
	}

	l.logWithLabel(
		l.labelError(),
		l.estimatedMessageSizeError,
		args,
	)
}
