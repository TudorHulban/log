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
	var (
		file string
		line int
	)

	if l.withCaller {
		_, fileCaller, lineCaller, _ := runtime.Caller(int(l.callerLevel))
		file = fileCaller
		line = lineCaller
	}

	if l.withJSON {
		region, errWrite := l.ingestor.TryWrite(
			estimatedMessageSize +
				l.estimateJSONOverhead(len(msg), file, line, keysAndValues),
		)
		if errWrite != nil {
			return
		}

		buf := region.Buf()[:0]

		buf = l.appendJSONKV(
			buf,
			label,
			file,
			line,
			[]byte(msg),
			keysAndValues...,
		)

		buf = append(buf, '\n')

		copy(region.Buf(), buf)

		l.ingestor.EndWrite(region)

		return
	}

	// Non‑JSON path
	region, errWrite := l.ingestor.TryWrite(
		estimatedMessageSize +
			l.estimateJSONOverhead(len(msg), file, line, keysAndValues),
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
	buf = append(buf, ' ')

	if l.withCaller {
		buf = append(buf, file...)
		buf = append(buf, ' ')
		buf = append(buf, 'L', 'i', 'n', 'e')
		buf = strconv.AppendInt(buf, int64(line), 10)
		buf = append(buf, ' ')
	}

	buf = append(buf, label...)
	buf = append(buf, delim...)

	buf = helpers.AppendKeyValues(buf, keysAndValues...)
	buf = append(buf, '\n')

	copy(region.Buf(), buf)

	l.ingestor.EndWrite(region)
}

func (l *Logger) estimateJSONOverhead(msgLen int, file string, line int, kv []any) uint32 {
	var size uint32 = 64 // base JSON overhead

	// timestamp
	if l.fnTimestamp != nil {
		size += 32 // worst case timestamp length
	}

	// level
	size = size + 5 + 10 // 5 - level label, 10 - for JSON

	// caller info
	if len(file) > 0 && line > 0 {
		size += uint32(len(file)) + 20 // "caller":"...","line":123,
	}

	// message field: "msg":"<escaped>"
	// worst case: every char becomes \u00XX (6 bytes)
	size += uint32(msgLen) * 2
	size += 10 // field name + quotes + comma

	// key/value pairs
	for i := 0; i < len(kv); i += 2 {
		key := kv[i]
		val := kv[i+1]

		// key
		switch k := key.(type) {
		case string:
			size += uint32(len(k))*2 + 4
		case []byte:
			size += uint32(len(k))*2 + 4
		default:
			size += 16
		}

		// value
		switch v := val.(type) {
		case string:
			size += uint32(len(v))*2 + 4
		case []byte:
			size += uint32(len(v))*2 + 4
		case int, int32, int64, uint, uint64:
			size += 20
		case float32, float64:
			size += 32
		case bool:
			size += 5
		case nil:
			size += 4
		default:
			size += 32
		}
	}

	// newline
	size++

	return size
}
