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

	// timestamp
	if l.fnTimestamp != nil {
		buf = l.fnTimestamp(buf)
		buf = append(buf, ' ')
	}

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

	// args
	buf = helpers.AppendArgs(buf, args)
	buf = append(buf, '\n')

	// ingestion
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

	var buf []byte

	// timestamp
	if l.fnTimestamp != nil {
		buf = l.fnTimestamp(buf)
		buf = append(buf, ' ')
	}

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

	// formatted message (zero‑alloc)
	buf = helpers.Appendf(buf, format, args)
	buf = append(buf, '\n')

	// ingestion
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
