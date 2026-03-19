package log

import (
	"fmt"

	"github.com/tudorhulban/log/helpers"
)

func (l *Logger) PrintMessage(msg string) {
	var arr [256]byte

	buf := arr[:0] // stack-allocated, no heap alloc

	buf = append(buf, l.fnTimestamp(buf)...)
	buf = append(buf, ' ')
	buf = append(buf, []byte(msg)...)
	buf = append(buf, '\n')

	region, errWrite := l.ingestor.TryWrite(uint32(len(buf)))
	if errWrite == nil {
		copy(region.Buf(), buf)

		l.ingestor.EndWrite(region)
	}
}

func (l *Logger) Print(args ...any) {
	var arr [256]byte

	buf := arr[:0] // stack-allocated, no heap alloc

	if l.fnTimestamp != nil {
		buf = append(buf, l.fnTimestamp(buf)...)
	}

	buf = append(buf, ' ')
	buf = helpers.AppendArgs(buf, args)
	buf = append(buf, '\n')

	region, errWrite := l.ingestor.TryWrite(uint32(len(buf)))
	if errWrite == nil {
		copy(region.Buf(), buf)

		l.ingestor.EndWrite(region)
	}
}

func (l *Logger) PrintWithNoTimestamp(args ...any) {
	var arr [256]byte

	buf := arr[:0] // stack-allocated, no heap alloc

	buf = helpers.AppendArgs(buf, args)
	buf = append(buf, '\n')

	region, errWrite := l.ingestor.TryWrite(uint32(len(buf)))
	if errWrite == nil {
		copy(region.Buf(), buf)

		l.ingestor.EndWrite(region)
	}
}

func (l *Logger) PrintFast(msg []byte) {
	region, errWrite := l.ingestor.TryWrite(uint32(len(msg)))
	if errWrite == nil {
		copy(region.Buf(), msg)

		l.ingestor.EndWrite(region)
	}
}

func (l *Logger) Printw(msg string, args ...any) {
	var arr [256]byte

	buf := arr[:0] // stack-allocated, no heap alloc

	if l.fnTimestamp != nil {
		buf = append(buf, l.fnTimestamp(buf)...)
		buf = append(buf, ' ')
	}

	buf = append(buf, []byte(msg)...)
	buf = append(buf, '\n')
	buf = helpers.AppendArgs(buf, args)
	buf = append(buf, '\n')

	region, errWrite := l.ingestor.TryWrite(uint32(len(buf)))
	if errWrite == nil {
		copy(region.Buf(), buf)

		l.ingestor.EndWrite(region)
	}
}

func (l *Logger) Printf(format string, args ...any) {
	var arr [256]byte

	buf := arr[:0] // stack-allocated, no heap alloc

	if l.withJSON {
		buf = l.appendJSON(
			buf,
			l.fnTimestamp(buf),
			l.labelInfo(),
			format, args...,
		)

		_, _ = l.ingestor.Write(buf)
	} else {
		if l.fnTimestamp != nil {
			buf = append(buf, l.fnTimestamp(buf)...)
			buf = append(buf, ' ')
		}

		buf = fmt.Appendf(buf, format, args...)
		buf = append(buf, '\n')

		_, _ = l.ingestor.Write(buf)
	}
}
