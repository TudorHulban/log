package log

import (
	"runtime"
	"strconv"

	"github.com/tudorhulban/log/helpers"
)

func (l *Logger) logWithLabel(label string, estimatedMessageSize uint32, args []any) {
	region, errWrite := l.ingestor.TryWrite(estimatedMessageSize + _DeltaEstimation)
	if errWrite != nil {
		return
	}

	buf := region.Buf()[:0]

	if l.withJSON {
		var (
			file string
			line int
		)

		if l.withCaller {
			_, fileCaller, lineCaller, _ := runtime.Caller(int(l.callerLevel))
			file = fileCaller
			line = lineCaller
		}

		buf = l.appendJSON(
			buf,
			label,
			file,
			line,
			helpers.AppendArgs(nil, args...),
		)

		buf = append(buf, '\n')

		copy(region.Buf(), buf)
		l.ingestor.EndWrite(region)

		return
	}

	// Non‑JSON path
	if l.fnTimestamp != nil {
		buf = l.fnTimestamp(buf)
		buf = append(buf, ' ')
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(int(l.callerLevel))

		buf = append(buf, file...)
		buf = append(buf, ' ')
		buf = append(buf, 'L', 'i', 'n', 'e')
		buf = strconv.AppendInt(buf, int64(line), 10)
		buf = append(buf, ' ')
	}

	buf = append(buf, label...)
	buf = append(buf, delim...)
	buf = helpers.AppendArgs(buf, args...)
	buf = append(buf, '\n')

	copy(region.Buf(), buf)

	l.ingestor.EndWrite(region)
}

func (l *Logger) logfWithLabel(label, format string, estimatedMessageSize uint32, args []any) {
	region, errWrite := l.ingestor.TryWrite(estimatedMessageSize + _DeltaEstimation)
	if errWrite != nil {
		return
	}

	buf := region.Buf()[:0]

	if l.withJSON {
		var (
			file string
			line int
		)

		if l.withCaller {
			_, fileCaller, lineCaller, _ := runtime.Caller(int(l.callerLevel))
			file = fileCaller
			line = lineCaller
		}

		buf = l.appendJSON(
			buf,
			label,
			file,
			line,
			helpers.Appendf(nil, format, args),
		)

		buf = append(buf, '\n')

		copy(region.Buf(), buf)
		l.ingestor.EndWrite(region)

		return
	}

	// Non‑JSON path
	if l.fnTimestamp != nil {
		buf = l.fnTimestamp(buf)
		buf = append(buf, ' ')
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(int(l.callerLevel))

		buf = append(buf, file...)
		buf = append(buf, ' ')
		buf = append(buf, 'L', 'i', 'n', 'e')
		buf = strconv.AppendInt(buf, int64(line), 10)
		buf = append(buf, ' ')
	}

	buf = append(buf, label...)
	buf = append(buf, delim...)
	buf = helpers.Appendf(buf, format, args)
	buf = append(buf, '\n')

	copy(region.Buf(), buf)

	l.ingestor.EndWrite(region)
}

func (l *Logger) logwWithLabel(label, msg string, estimatedMessageSize uint32, keysAndValues ...any) {
	region, errWrite := l.ingestor.TryWrite(estimatedMessageSize + _DeltaEstimation)
	if errWrite != nil {
		return
	}

	buf := region.Buf()[:0]

	if l.withJSON {
		var (
			file string
			line int
		)

		if l.withCaller {
			_, fileCaller, lineCaller, _ := runtime.Caller(int(l.callerLevel))
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

		copy(region.Buf(), buf)
		l.ingestor.EndWrite(region)

		return
	}

	// Non‑JSON path
	if l.fnTimestamp != nil {
		buf = l.fnTimestamp(buf)
		buf = append(buf, ' ')
	}

	buf = append(buf, msg...)
	buf = append(buf, '\n')

	if l.withCaller {
		_, file, line, _ := runtime.Caller(int(l.callerLevel))

		buf = append(buf, file...)
		buf = append(buf, ' ')
		buf = append(buf, 'L', 'i', 'n', 'e')
		buf = strconv.AppendInt(buf, int64(line), 10)
		buf = append(buf, ' ')
	}

	buf = append(buf, label...)
	buf = append(buf, delim...)

	buf = helpers.AppendArgs(buf, keysAndValues)
	buf = append(buf, '\n')

	copy(region.Buf(), buf)

	l.ingestor.EndWrite(region)
}
