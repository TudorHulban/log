package log

import (
	"runtime"

	"github.com/tudorhulban/log/helpers"
)

func (l *Logger) DebugFast(args ...any) {
	if l.logLevel < LevelDEBUG {
		return
	}

	region, errWrite := l.ingestor.TryWrite(l.estimatedMessageSize)
	if errWrite != nil {
		return
	}

	buf := region.Buf()[:0]

	if l.withJSON {
		var file string
		var line int

		if l.withCaller {
			_, fileCaller, lineCaller, _ := runtime.Caller(l.callerLevel)

			file = fileCaller
			line = lineCaller
		}

		buf = l.appendJSON(
			buf,
			l.labelDebug(),
			file,
			line,
			(helpers.AppendArgs(nil, args)),
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
		_, file, line, _ := runtime.Caller(1)

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

	copy(region.Buf(), buf)

	l.ingestor.EndWrite(region)
}
