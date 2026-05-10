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
	region, errWrite := l.ingestor.TryWrite(uint32(len(msg)))
	if errWrite != nil {
		return
	}

	buf := region.Buf()[:0]

	if l.withJSON {
		buf = l.appendJSON(
			buf,
			l.labelPrint(),
			"",
			0,
			[]byte(msg),
		)

		copy(region.Buf(), buf)
		l.ingestor.EndWrite(region)

		return
	}

	if l.fnTimestamp != nil {
		buf = l.fnTimestamp(buf)
		buf = append(buf, ' ')
	}

	buf = append(buf, msg...)
	buf = append(buf, '\n')

	copy(region.Buf(), buf)
	l.ingestor.EndWrite(region)
}

func (l *Logger) Print(args ...any) {
	l.logWithLabel(
		l.labelPrint(),
		helpers.GetEstimatedMessageSize("", args),
		args,
	)
}

func (l *Logger) PrintWithNoTimestamp(args ...any) {
	region, errWrite := l.ingestor.TryWrite(
		helpers.GetEstimatedMessageSize("", args) + _DeltaEstimation,
	)
	if errWrite != nil {
		return
	}

	buf := region.Buf()[:0]

	buf = helpers.AppendArgs(buf, args)
	buf = append(buf, '\n')

	copy(region.Buf(), buf)
	l.ingestor.EndWrite(region)
}

func (l *Logger) Printw(msg string, args ...any) {
	region, errWrite := l.ingestor.TryWrite(
		uint32(len(msg)) + helpers.GetEstimatedMessageSize("", args) + _DeltaEstimation,
	)
	if errWrite != nil {
		return
	}

	buf := region.Buf()[:0]

	if l.fnTimestamp != nil {
		buf = l.fnTimestamp(buf)
		buf = append(buf, ' ')
	}

	buf = append(buf, msg...)
	buf = append(buf, '\n')
	buf = helpers.AppendArgs(buf, args)
	buf = append(buf, '\n')

	copy(region.Buf(), buf)
	l.ingestor.EndWrite(region)
}

func (l *Logger) Printf(format string, args ...any) {
	region, errWrite := l.ingestor.TryWrite(
		helpers.GetEstimatedMessageSize(format, args) + _DeltaEstimation,
	)
	if errWrite != nil {
		return
	}

	buffer := region.Buf()[:0]

	if l.withJSON {
		buffer = l.appendJSON(
			buffer,
			l.labelPrint(),
			"",
			0,
			fmt.Appendf(nil, format, args...),
		)

		copy(region.Buf(), buffer)
		l.ingestor.EndWrite(region)
	} else {
		if l.fnTimestamp != nil {
			buffer = l.fnTimestamp(buffer)
			buffer = append(buffer, ' ')
		}

		buffer = fmt.Appendf(buffer, format, args...)
		buffer = append(buffer, '\n')

		copy(region.Buf(), buffer)

		l.ingestor.EndWrite(region)
	}
}

func (l *Logger) PrintRaw(msg []byte) {
	region, errWrite := l.ingestor.TryWrite(uint32(len(msg))) //nolint:gosec
	if errWrite == nil {
		copy(region.Buf(), msg)

		l.ingestor.EndWrite(region)
	}
}
