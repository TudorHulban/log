package log

import "github.com/tudorhulban/log/helpers"

func (l *Logger) labelWarn() string {
	if l.withColor {
		return colorWarn(logLevels[LevelWarn])
	}

	return logLevels[LevelWarn]
}

func (l *Logger) Warn(args ...any) {
	if Level(l.logLevel.Load()) > LevelWarn {
		return
	}

	l.logWithLabel(
		l.labelWarn(),
		helpers.GetEstimatedMessageSize("", args),
		args,
	)
}

func (l *Logger) Warnf(format string, args ...any) {
	if Level(l.logLevel.Load()) > LevelWarn {
		return
	}

	l.logfWithLabel(
		l.labelWarn(),
		format,
		helpers.GetEstimatedMessageSize(format, args),
		args,
	)
}

func (l *Logger) Warnw(msg string, keysAndValues ...any) {
	if Level(l.logLevel.Load()) > LevelWarn {
		return
	}

	l.logwWithLabel(
		l.labelWarn(),
		msg,
		uint32(len(msg))+helpers.GetEstimatedMessageSize("", keysAndValues),
		keysAndValues...,
	)
}
