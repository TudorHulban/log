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

func (l *Logger) PrintMessage(msg string) {
	var buf []byte

	if l.fnTimestamp != nil {
		buf = l.fnTimestamp(buf)
		buf = append(buf, ' ')
	}

	buf = append(buf, []byte(msg)...)
	buf = append(buf, '\n')

	region, errWrite := l.ingestor.TryWrite(uint32(len(buf))) //nolint:gosec
	if errWrite == nil {
		copy(region.Buf(), buf)

		l.ingestor.EndWrite(region)
	}
}

func (l *Logger) Print(args ...any) {
	var buf []byte

	if l.fnTimestamp != nil {
		buf = l.fnTimestamp(buf)
		buf = append(buf, ' ')
	}

	buf = helpers.AppendArgs(buf, args)
	buf = append(buf, '\n')

	region, errWrite := l.ingestor.TryWrite(uint32(len(buf))) //nolint:gosec
	if errWrite == nil {
		copy(region.Buf(), buf)

		l.ingestor.EndWrite(region)
	}
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

	buf = append(buf, []byte(msg)...)
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
	var buf []byte

	if l.withJSON {
		buf = l.appendJSON(
			buf,
			l.labelInfo(),
			"",
			0,

			// string(helpers.Appendf(nil, format, args)),
			string(fmt.Appendf(nil, format, args...)),
		)
	} else {
		if l.fnTimestamp != nil {
			buf = l.fnTimestamp(buf)
			buf = append(buf, ' ')
		}

		buf = fmt.Appendf(buf, format, args...)
		buf = append(buf, '\n')
	}

	region, errWrite := l.ingestor.TryWrite(uint32(len(buf))) //nolint:gosec
	if errWrite == nil {
		copy(region.Buf(), buf)

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
