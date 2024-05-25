package log

import (
	"io"

	safewriter "github.com/TudorHulban/log/safe-writer"
)

type Level int8

type Logger struct {
	w        *safewriter.SafeWriterInfo
	logLevel int8

	withCaller bool // for shorter form in case do not need caller file.
}

func NewLogger(level Level, writeTo io.Writer, withCaller bool) Logger {
	lev := convertLevel(level)

	if writeTo == nil {
		writeTo = io.Discard
	}

	res := Logger{
		logLevel:   lev,
		w:          safewriter.NewSafeWriterInfo(writeTo),
		withCaller: withCaller,
	}

	go res.w.Writer.Listen()

	res.Printf(
		"created logger, level %v.",
		logLevels[lev],
	)

	return res
}

func (l *Logger) write(payload []byte) {
	l.w.ChWrites <- payload
}

func (l *Logger) Close() {
	// TODO: release resources
}
