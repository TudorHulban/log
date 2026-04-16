package log

import (
	"runtime"

	"github.com/tudorhulban/log/helpers"
)

func (l *Logger) logWithLabel(label string, args ...any) {
	if l.withJSON {
		var (
			file string
			line int
		)

		buf := make([]byte, 0, _PreallocationJSON)

		if l.withCaller {
			_, fileCaller, lineCaller, _ := runtime.Caller(l.callerLevel)
			file = fileCaller
			line = lineCaller
		}

		buf = l.appendJSON(
			buf,
			label,
			file,
			line,
			helpers.AppendArgs(nil, args),
		)

		buf = append(buf, '\n')

		region, errWrite := l.ingestor.TryWrite(uint32(len(buf))) //nolint:gosec
		if errWrite == nil {
			copy(region.Buf(), buf)
			l.ingestor.EndWrite(region)
		}

		return
	}

	// Non‑JSON path
	buf := make([]byte, 0, _PreallocationBuffer)

	if l.fnTimestamp != nil {
		buf = l.fnTimestamp(buf)
		buf = append(buf, ' ')
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(l.callerLevel)

		buf = append(buf, file...)
		buf = append(buf, ' ')
		buf = append(buf, 'L', 'i', 'n', 'e')
		buf = helpers.AppendInt(buf, line)
		buf = append(buf, ' ')
	}

	buf = append(buf, label...)
	buf = append(buf, delim...)
	buf = helpers.AppendArgs(buf, args)
	buf = append(buf, '\n')

	region, errWrite := l.ingestor.TryWrite(uint32(len(buf))) //nolint:gosec
	if errWrite == nil {
		copy(region.Buf(), buf)

		l.ingestor.EndWrite(region)
	}
}

func (l *Logger) logfWithLabel(label string, format string, args ...any) {
	if l.withJSON {
		var (
			file string
			line int
		)

		buf := make([]byte, 0, _PreallocationJSON)

		if l.withCaller {
			_, fileCaller, lineCaller, _ := runtime.Caller(l.callerLevel)
			file = fileCaller
			line = lineCaller
		}

		buf = l.appendJSON(
			buf,
			label,
			file,
			line,
			helpers.Appendf(nil, format, args...),
		)

		buf = append(buf, '\n')

		region, errWrite := l.ingestor.TryWrite(uint32(len(buf))) //nolint:gosec
		if errWrite == nil {
			copy(region.Buf(), buf)
			l.ingestor.EndWrite(region)
		}

		return
	}

	// Non‑JSON path
	buf := make([]byte, 0, _PreallocationBuffer)

	if l.fnTimestamp != nil {
		buf = l.fnTimestamp(buf)
		buf = append(buf, ' ')
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(l.callerLevel)

		buf = append(buf, file...)
		buf = append(buf, ' ')
		buf = append(buf, 'L', 'i', 'n', 'e')
		buf = helpers.AppendInt(buf, line)
		buf = append(buf, ' ')
	}

	buf = append(buf, label...)
	buf = append(buf, delim...)
	buf = helpers.Appendf(buf, format, args...)
	buf = append(buf, '\n')

	region, errWrite := l.ingestor.TryWrite(uint32(len(buf))) //nolint:gosec
	if errWrite == nil {
		copy(region.Buf(), buf)

		l.ingestor.EndWrite(region)
	}
}

func (l *Logger) logwWithLabel(label string, msg string, keysAndValues ...any) {
	if l.withJSON {
		var (
			file string
			line int
		)

		buf := make([]byte, 0, _PreallocationJSON)

		if l.withCaller {
			_, fileCaller, lineCaller, _ := runtime.Caller(l.callerLevel)
			file = fileCaller
			line = lineCaller
		}

		buf = l.appendJSON(
			buf,
			label,
			file,
			line,
			helpers.AppendArgs(nil, append([]any{msg}, keysAndValues...)),
		)

		buf = append(buf, '\n')

		region, errWrite := l.ingestor.TryWrite(uint32(len(buf))) //nolint:gosec
		if errWrite == nil {
			copy(region.Buf(), buf)
			l.ingestor.EndWrite(region)
		}

		return
	}

	// Non‑JSON path
	buf := make([]byte, 0, _PreallocationBuffer)

	if l.fnTimestamp != nil {
		buf = l.fnTimestamp(buf)
		buf = append(buf, ' ')
	}

	buf = append(buf, msg...)
	buf = append(buf, '\n')

	if l.withCaller {
		_, file, line, _ := runtime.Caller(l.callerLevel)

		buf = append(buf, file...)
		buf = append(buf, ' ')
		buf = append(buf, 'L', 'i', 'n', 'e')
		buf = helpers.AppendInt(buf, line)
		buf = append(buf, ' ')
	}

	buf = append(buf, label...)
	buf = append(buf, delim...)

	buf = helpers.AppendArgs(buf, keysAndValues)
	buf = append(buf, '\n')

	region, errWrite := l.ingestor.TryWrite(uint32(len(buf))) //nolint:gosec
	if errWrite == nil {
		copy(region.Buf(), buf)

		l.ingestor.EndWrite(region)
	}
}
