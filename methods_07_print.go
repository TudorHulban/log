package log

import (
	"fmt"

	"github.com/tudorhulban/log/helpers"
)

// - Print / PrintMessage / PrintWithNoTimestamp / Printw / Printf
//   Build into their own buffer first, then reserve exactly what is needed.
//   Guaranteed no truncation.
//   One allocation per call.
//   Use when message size is unbounded or correctness is required.

// PrintRaw is always safe.
// The caller owns the buffer and the reservation is sized to match.

func (l *Logger) labelPrint() string {
	if l.withColor {
		return colorDebug(logLevels[l.GetLogLevel()])
	}

	return logLevels[l.GetLogLevel()]
}

func (l *Logger) PrintMessage(msg string) {
	if l.withJSON {
		buf := make([]byte, 0, _PreallocationJSON)

		buf = l.appendJSON(
			buf,
			l.labelInfo(),
			"",
			0,
			[]byte(msg),
		)

		region, errWrite := l.ingestor.TryWrite(uint32(len(buf))) //nolint:gosec
		if errWrite == nil {
			copy(region.Buf(), buf)
			l.ingestor.EndWrite(region)
		}

		return
	}

	var buf []byte

	if l.fnTimestamp != nil {
		buf = l.fnTimestamp(buf)
		buf = append(buf, ' ')
	}

	buf = append(buf, msg...)
	buf = append(buf, '\n')

	region, errWrite := l.ingestor.TryWrite(uint32(len(buf))) //nolint:gosec
	if errWrite == nil {
		copy(region.Buf(), buf)
		l.ingestor.EndWrite(region)
	}
}

func (l *Logger) Print(args ...any) {
	estimatedMessageSizeInfo := helpers.GetEstimatedMessageSize("", args)

	l.logWithLabel(
		l.labelPrint(),
		uint32(estimatedMessageSizeInfo),
		args,
	)
}

func (l *Logger) PrintWithNoTimestamp(args ...any) {
	var buf []byte

	buf = helpers.AppendArgs(buf, args)
	buf = append(buf, '\n')

	region, errWrite := l.ingestor.TryWrite(uint32(len(buf))) //nolint:gosec
	if errWrite == nil {
		copy(region.Buf(), buf)

		l.ingestor.EndWrite(region)
	}
}

func (l *Logger) Printw(msg string, args ...any) {
	var buf []byte

	if l.fnTimestamp != nil {
		buf = l.fnTimestamp(buf)
		buf = append(buf, ' ')
	}

	buf = append(buf, msg...)
	buf = append(buf, '\n')
	buf = helpers.AppendArgs(buf, args)
	buf = append(buf, '\n')

	region, errWrite := l.ingestor.TryWrite(uint32(len(buf))) //nolint:gosec
	if errWrite == nil {
		copy(region.Buf(), buf)

		l.ingestor.EndWrite(region)
	}
}

func (l *Logger) Printf(format string, args ...any) {
	if l.withJSON {
		buf := make([]byte, 0, _PreallocationJSON)

		buf = l.appendJSON(
			buf,
			l.labelInfo(),
			"",
			0,
			fmt.Appendf(nil, format, args...),
		)

		region, errWrite := l.ingestor.TryWrite(uint32(len(buf))) //nolint:gosec
		if errWrite == nil {
			copy(region.Buf(), buf)

			l.ingestor.EndWrite(region)
		}
	} else {
		buf := make([]byte, 0, _PreallocationBuffer)

		if l.fnTimestamp != nil {
			buf = l.fnTimestamp(buf)
			buf = append(buf, ' ')
		}

		buf = fmt.Appendf(buf, format, args...)
		buf = append(buf, '\n')

		region, errWrite := l.ingestor.TryWrite(uint32(len(buf))) //nolint:gosec
		if errWrite == nil {
			copy(region.Buf(), buf)

			l.ingestor.EndWrite(region)
		}
	}
}

func (l *Logger) PrintRaw(msg []byte) {
	region, errWrite := l.ingestor.TryWrite(uint32(len(msg))) //nolint:gosec
	if errWrite == nil {
		copy(region.Buf(), msg)

		l.ingestor.EndWrite(region)
	}
}
