package log

import (
	"fmt"
	"runtime"
	"strconv"
)

func (l *Logger) Info(args ...any) {
	if l.logLevel == 0 {
		return
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		l.write(
			[]byte(
				timestamp() + " " + file + " Line" + delim +
					strconv.FormatInt(int64(line), 10) + " " +
					logLevels[1] + delim + fmt.Sprint(args...) + "\n",
			),
		)

		return
	}

	l.write(
		[]byte(
			timestamp() + " " + logLevels[1] + delim +
				fmt.Sprint(args...) + "\n",
		),
	)
}

func (l *Logger) Infof(format string, args ...any) {
	if l.logLevel == 0 {
		return
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		l.write(
			[]byte(
				timestamp() + " " + file + " Line" + delim +
					strconv.FormatInt(int64(line), 10) + " " +
					logLevels[1] + delim + fmt.Sprintf(format, args...) + "\n",
			),
		)

		return
	}

	l.write(
		[]byte(
			timestamp() + " " + logLevels[1] + delim +
				fmt.Sprintf(format, args...) + "\n",
		),
	)
}
