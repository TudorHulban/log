package log

import (
	"fmt"

	"github.com/tudorhulban/log/helpers"
)

func (l *Logger) PrintMessage(msg string) {
	region, errWrite := l.ingestor.TryWrite(l.estimatedMessageSize)
	if errWrite == nil {
		buf := region.Buf()[:0]

		if l.fnTimestamp != nil {
			buf = append(buf, l.fnTimestamp()...)
			buf = append(buf, ' ')
		}

		buf = append(buf, []byte(msg)...)
		buf = append(buf, '\n')

		copy(region.Buf(), buf)

		l.ingestor.EndWrite(region)
	}
}

func (l *Logger) Print(args ...any) {
	region, errWrite := l.ingestor.TryWrite(l.estimatedMessageSize)
	if errWrite == nil {
		buf := region.Buf()[:0]

		if l.fnTimestamp != nil {
			buf = append(buf, l.fnTimestamp()...)
			buf = append(buf, ' ')
		}

		buf = helpers.AppendArgs(buf, args)
		buf = append(buf, '\n')

		copy(region.Buf(), buf)

		l.ingestor.EndWrite(region)
	}
}

func (l *Logger) PrintWithNoTimestamp(args ...any) {
	region, errWrite := l.ingestor.TryWrite(l.estimatedMessageSize)
	if errWrite == nil {
		buf := region.Buf()[:0]

		buf = helpers.AppendArgs(buf, args)
		buf = append(buf, '\n')

		copy(region.Buf(), buf)

		l.ingestor.EndWrite(region)
	}
}

func (l *Logger) Printw(msg string, args ...any) {
	region, errWrite := l.ingestor.TryWrite(l.estimatedMessageSize)
	if errWrite == nil {
		buf := region.Buf()[:0]

		if l.fnTimestamp != nil {
			buf = append(buf, l.fnTimestamp()...)
			buf = append(buf, ' ')
		}

		buf = append(buf, []byte(msg)...)
		buf = append(buf, '\n')
		buf = helpers.AppendArgs(buf, args)
		buf = append(buf, '\n')

		copy(region.Buf(), buf)

		l.ingestor.EndWrite(region)
	}
}

func (l *Logger) Printf(format string, args ...any) {
	region, errWrite := l.ingestor.TryWrite(l.estimatedMessageSize)
	if errWrite == nil {
		buf := region.Buf()[:0]

		if l.withJSON {
			if l.fnTimestamp != nil {
				buf = l.appendJSON(
					buf,
					l.fnTimestamp(),
					l.labelInfo(),
					format, args...,
				)
			} else {
				buf = l.appendJSON(
					buf,
					nil,
					l.labelInfo(),
					format, args...,
				)
			}

			copy(region.Buf(), buf)

			l.ingestor.EndWrite(region)
		} else {
			if l.fnTimestamp != nil {
				buf = append(buf, l.fnTimestamp()...)
				buf = append(buf, ' ')
			}

			buf = fmt.Appendf(buf, format, args...)
			buf = append(buf, '\n')

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
