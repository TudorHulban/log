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

	l.logWithLabel(
		l.labelInfo(),
		uint32(
			helpers.GetEstimatedMessageSize("", args),
		),
		args,
	)
}

func (l *Logger) Infof(format string, args ...any) {
	if Level(l.logLevel.Load()) > LevelInfo {
		return
	}

	l.logfWithLabel(
		l.labelInfo(),
		format,
		uint32(
			helpers.GetEstimatedMessageSize(format, args),
		),
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
		l.labelInfo(),
		msg,
		uint32(
			len(msg)+helpers.GetEstimatedMessageSize("", keysAndValues),
		),
		keysAndValues,
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
