package log

import (
	"runtime"

	"github.com/tudorhulban/log/helpers"
)

func (l Logger) labelDebug() string {
	if l.withColor {
		return colorDebug(logLevels[LevelDEBUG])
	} else {
		return logLevels[LevelDEBUG]
	}
}

func (l *Logger) Debug(args ...any) {
	if l.logLevel < LevelDEBUG {
		return
	}

	var buf []byte

	if l.withJSON {
		var file string
		var line int

		if l.withCaller {
			_, fileCaller, lineCaller, _ := runtime.Caller(l.callerLevel)
			file = fileCaller
			line = lineCaller
		}

		msg := helpers.AppendArgs(nil, args)

		// JSON append: timestamp, label, caller, args
		buf = l.appendJSON(
			buf,
			l.labelDebug(),
			file,
			line,
			string(msg),
		)

		buf = append(buf, '\n')

		region, errWrite := l.ingestor.TryWrite(uint32(len(buf)))
		if errWrite == nil {
			copy(region.Buf(), buf)

			l.ingestor.EndWrite(region)
		}

		return
	}

	// Non‑JSON path
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

	buf = append(buf, l.labelDebug()...)
	buf = append(buf, delim...)
	buf = helpers.AppendArgs(buf, args)
	buf = append(buf, '\n')

	region, errWrite := l.ingestor.TryWrite(uint32(len(buf)))
	if errWrite == nil {
		copy(region.Buf(), buf)

		l.ingestor.EndWrite(region)
	}
}

func (l *Logger) Debugf(format string, args ...any) {
	if l.logLevel < LevelDEBUG {
		return
	}

	// JSON path
	if l.withJSON {
		var buf []byte

		var file string
		var line int

		if l.withCaller {
			_, fileCaller, lineCaller, _ := runtime.Caller(1)
			file = fileCaller
			line = lineCaller
		}

		buf = l.appendJSON(
			buf,
			l.labelDebug(),
			file,
			line,

			string((helpers.Appendf(nil, format, args...))),
		)

		buf = append(buf, '\n')

		region, errWrite := l.ingestor.TryWrite(uint32(len(buf)))
		if errWrite == nil {
			copy(region.Buf(), buf)

			l.ingestor.EndWrite(region)
		}

		return
	}

	// Non‑JSON path
	var buf []byte

	if l.fnTimestamp != nil {
		buf = l.fnTimestamp(buf)
		buf = append(buf, ' ')
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		buf = append(buf, file...)
		buf = append(buf, ' ')
		buf = append(buf, 'L', 'i', 'n', 'e')
		buf = helpers.AppendInt(buf, line)
		buf = append(buf, ' ')
	}

	buf = append(buf, l.labelDebug()...)
	buf = append(buf, delim...)

	buf = helpers.Appendf(buf, format, args)
	buf = append(buf, '\n')

	region, errWrite := l.ingestor.TryWrite(uint32(len(buf)))
	if errWrite == nil {
		copy(region.Buf(), buf)

		l.ingestor.EndWrite(region)
	}
}

func (l *Logger) Debugw(msg string, keysAndValues ...any) {
	if l.logLevel < LevelDEBUG {
		return
	}

	var buf []byte

	// timestamp
	if l.fnTimestamp != nil {
		buf = l.fnTimestamp(buf)
		buf = append(buf, ' ')
	}

	// message
	buf = append(buf, msg...)
	buf = append(buf, '\n')

	// caller info
	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		buf = append(buf, file...)
		buf = append(buf, ' ')
		buf = append(buf, 'L', 'i', 'n', 'e')
		buf = helpers.AppendInt(buf, line)
		buf = append(buf, ' ')
	}

	// label
	buf = append(buf, l.labelDebug()...)
	buf = append(buf, delim...)

	// structured key/value pairs
	buf = helpers.AppendArgs(buf, keysAndValues)
	buf = append(buf, '\n')

	// ingestion
	region, errWrite := l.ingestor.TryWrite(uint32(len(buf)))
	if errWrite == nil {
		copy(region.Buf(), buf)
		l.ingestor.EndWrite(region)
	}
}
