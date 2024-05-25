package log

import (
	"fmt"
	"runtime"
	"strconv"
)

func (l *Logger) Warn(args ...interface{}) {
	if l.logLevel < 2 {
		return
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		l.write(
			[]byte(
				timestamp() + " " + file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + colorWarn()(logLevels[2]) + delim + fmt.Sprint(args...) + "\n",
			),
		)

		return
	}

	l.write(
		[]byte(
			timestamp() + " " + colorWarn()(logLevels[2]) + delim + fmt.Sprint(args...) + "\n",
		),
	)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	if l.logLevel < 2 {
		return
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		l.write(
			[]byte(
				timestamp() + " " + file + " Line" + delim + strconv.FormatInt(int64(line), 10) + " " + colorWarn()(logLevels[2]) + delim + fmt.Sprintf(format, args...) + "\n",
			),
		)

		return
	}

	l.write(
		[]byte(
			timestamp() + " " + colorWarn()(logLevels[2]) + delim + fmt.Sprintf(format, args...) + "\n",
		),
	)

}
