package log

import (
	"fmt"
	"runtime"
	"strconv"
)

func (l Logger) Warn(args ...any) {
	if l.logLevel < 2 {
		return
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		l.localWriter.Write(
			[]byte(
				l.withTimestamp() + " " + file + " Line" + delim +
					strconv.FormatInt(int64(line), 10) + " " +
					colorWarn()(logLevels[2]) + delim +
					fmt.Sprint(args...) + "\n",
			),
		)

		return
	}

	l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " + colorWarn()(logLevels[2]) + delim +
				fmt.Sprint(args...) + "\n",
		),
	)
}

func (l Logger) Warnf(format string, args ...any) {
	if l.logLevel < 2 {
		return
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		l.localWriter.Write(
			[]byte(
				l.withTimestamp() + " " + file + " Line" + delim +
					strconv.FormatInt(int64(line), 10) + " " +
					colorWarn()(logLevels[2]) + delim +
					fmt.Sprintf(format, args...) + "\n",
			),
		)

		return
	}

	l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " + colorWarn()(logLevels[2]) + delim +
				fmt.Sprintf(format, args...) + "\n",
		),
	)
}

func (l Logger) Warnw(msg string, keysAndValues ...any) {
	if l.logLevel < 2 {
		return
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		l.localWriter.Write(
			[]byte(
				l.withTimestamp() + " " + file + " Line" + delim +
					strconv.FormatInt(int64(line), 10) + " " +
					colorWarn()(logLevels[2]) + delim +
					fmt.Sprint(keysAndValues...) + "\n",
			),
		)

		return
	}

	l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " + colorWarn()(logLevels[2]) + delim +
				fmt.Sprint(keysAndValues...) + "\n",
		),
	)
}
