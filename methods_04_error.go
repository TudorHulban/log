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

	l.logWithLabel(
		l.labelError(),
		helpers.GetEstimatedMessageSize("", args),
		args,
	)
}

func (l *Logger) Errorf(format string, args ...any) {
	if Level(l.logLevel.Load()) > LevelError {
		return
	}

	l.logfWithLabel(
		l.labelError(),
		format,
		helpers.GetEstimatedMessageSize(format, args),
		args,
	)
}

func (l *Logger) Errorw(msg string, keysAndValues ...any) {
	if Level(l.logLevel.Load()) > LevelError {
		return
	}

	l.logwWithLabel(
		l.labelError(),
		msg,
		uint32(len(msg))+helpers.GetEstimatedMessageSize("", keysAndValues),
		keysAndValues,
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
