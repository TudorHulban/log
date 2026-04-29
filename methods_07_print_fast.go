package log

import "github.com/tudorhulban/log/helpers"

// - PrintFast / PrintMessageFast / PrintWithNoTimestamp / Printw / Printf
//   Use a pre-allocated arena buffer of estimatedMessageSize.
//   Fast, zero alloc.
//   If the message exceeds estimatedMessageSize it is silently truncated.
//   Use when throughput matters and message size is predictable.

func (l *Logger) PrintFast(args ...any) {
	region, errWrite := l.ingestor.TryWrite(l.estimatedMessageSizeOverall)
	if errWrite == nil {
		buf := region.Buf()[:0]

		if l.fnTimestamp != nil {
			buf = l.fnTimestamp(buf)
			buf = append(buf, ' ')
		}

		buf = helpers.AppendArgs(buf, args...)
		buf = append(buf, '\n')

		copy(region.Buf(), buf)

		l.ingestor.EndWrite(region)
	}
}

func (l *Logger) PrintMessageFast(msg string) {
	region, errWrite := l.ingestor.TryWrite(l.estimatedMessageSizeOverall)
	if errWrite == nil {
		buf := region.Buf()[:0]

		if l.fnTimestamp != nil {
			buf = l.fnTimestamp(buf)
			buf = append(buf, ' ')
		}

		buf = append(buf, msg...)
		buf = append(buf, '\n')

		copy(region.Buf(), buf)

		l.ingestor.EndWrite(region)
	}
}

func (l *Logger) PrintWithNoTimestampFast(args ...any) {
	region, errWrite := l.ingestor.TryWrite(l.estimatedMessageSizeOverall)
	if errWrite == nil {
		buf := region.Buf()[:0]

		buf = helpers.AppendArgs(buf, args...)
		buf = append(buf, '\n')

		copy(region.Buf(), buf)

		l.ingestor.EndWrite(region)
	}
}

func (l *Logger) PrintwFast(msg string, args ...any) {
	region, errWrite := l.ingestor.TryWrite(l.estimatedMessageSizeOverall)
	if errWrite == nil {
		buf := region.Buf()[:0]

		if l.fnTimestamp != nil {
			buf = l.fnTimestamp(buf)
			buf = append(buf, ' ')
		}

		buf = append(buf, msg...)
		buf = append(buf, '\n')
		buf = helpers.AppendArgs(buf, args...)
		buf = append(buf, '\n')

		copy(region.Buf(), buf)

		l.ingestor.EndWrite(region)
	}
}

func (l *Logger) PrintfFast(format string, args ...any) {
	region, errWrite := l.ingestor.TryWrite(l.estimatedMessageSizeOverall)
	if errWrite == nil {
		buf := region.Buf()[:0]

		if l.withJSON {
			msg := helpers.Appendf(nil, format, args)

			buf = l.appendJSON(
				buf,
				l.labelInfo(),
				"",
				0,
				msg,
			)

			buf = append(buf, '\n')

			copy(region.Buf(), buf)

			l.ingestor.EndWrite(region)
		} else {
			if l.fnTimestamp != nil {
				buf = l.fnTimestamp(buf)
				buf = append(buf, ' ')
			}

			buf = helpers.Appendf(buf, format, args)
			buf = append(buf, '\n')

			copy(region.Buf(), buf)

			l.ingestor.EndWrite(region)
		}
	}
}
