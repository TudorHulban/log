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
	if l.withJSON {
		buf := make([]byte, 0, 256)
		buf = l.appendJSON(buf, l.fnTimestamp(buf), l.labelInfo(), format, args...)

		_, _ = l.localWriter.Write(buf)
	} else {
		var arr [256]byte
		buf := arr[:0] // stack-allocated, no heap alloc

		buf = append(buf, l.fnTimestamp(buf)...)
		buf = append(buf, ' ')
		buf = fmt.Appendf(buf, format, args...)
		buf = append(buf, '\n')

		_, _ = l.localWriter.Write(buf)
	}
}
