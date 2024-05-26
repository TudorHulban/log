package log

import (
	"fmt"
	"runtime"
	"strconv"
)

func (l *Logger) Debug(args ...any) {
	if l.logLevel < 3 {
		return
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		l.localWriter.Write(
			[]byte(
				l.withTimestamp() + " " + file + " Line" + delim +
					strconv.FormatInt(int64(line), 10) + " " +
					colorDebug()(logLevels[3]) + delim + fmt.Sprint(args...) + "\n",
			),
		)

		return
	}

	l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " + colorDebug()(logLevels[3]) + delim + fmt.Sprint(args...) + "\n",
		),
	)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.logLevel < 3 {
		return
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		l.localWriter.Write(
			[]byte(
				l.withTimestamp() + " " + file + " Line" + delim +
					strconv.FormatInt(int64(line), 10) + " " +
					colorDebug()(logLevels[3]) + delim + fmt.Sprintf(format, args...) + "\n",
			),
		)

		return
	}

	l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " + colorDebug()(logLevels[3]) + delim + fmt.Sprintf(format, args...) + "\n",
		),
	)

}
