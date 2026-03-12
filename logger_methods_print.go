package log

import (
	"fmt"
)

func (l *Logger) PrintMessage(msg string) {
	var arr [256]byte

	buf := arr[:0] // stack-allocated, no heap alloc

	buf = append(buf, l.fnTimestamp(buf)...)
	buf = append(buf, ' ')
	buf = append(buf, []byte(msg)...)
	buf = append(buf, '\n')

	_, _ = l.localWriter.Write(buf)
}

func (l *Logger) Print(args ...any) {
	var arr [256]byte

	buf := arr[:0] // stack-allocated, no heap alloc

	buf = append(buf, l.fnTimestamp(buf)...)
	buf = append(buf, ' ')
	buf = append(buf, fmt.Sprint(args...)...)
	buf = append(buf, '\n')

	_, _ = l.localWriter.Write(buf)
}

func (l *Logger) Printw(msg string, args ...any) {
	var arr [256]byte

	buf := arr[:0] // stack-allocated, no heap alloc

	buf = append(buf, l.fnTimestamp(buf)...)
	buf = append(buf, ' ')
	buf = append(buf, []byte(msg)...)
	buf = append(buf, '\n')
	buf = append(buf, fmt.Sprint(args...)...)
	buf = append(buf, '\n')

	_, _ = l.localWriter.Write(buf)
}

func (l *Logger) Printf(format string, args ...any) {
	l.buf = l.buf[:0] // reset without allocating

	if l.withJSON {
		l.buf = l.appendJSON(l.buf, l.fnTimestamp(l.buf), l.labelInfo(), format, args...)

		_, _ = l.localWriter.Write(l.buf)
	} else {
		l.buf = append(l.buf, l.fnTimestamp(l.buf)...)
		l.buf = append(l.buf, ' ')
		l.buf = fmt.Appendf(l.buf, format, args...)
		l.buf = append(l.buf, '\n')

		_, _ = l.localWriter.Write(l.buf)
	}
}
