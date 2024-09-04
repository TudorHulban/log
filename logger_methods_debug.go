package log

import (
	"fmt"
	"runtime"
	"strconv"
)

func (l Logger) labelDebug() string {
	return ternary(
		l.withColor,

		colorDebug(logLevels[LevelDEBUG]),
		logLevels[LevelDEBUG],
	)
}

func (l Logger) Debug(args ...any) {
	if l.logLevel < LevelDEBUG {
		return
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		l.localWriter.Write(
			[]byte(
				l.withTimestamp() +
					" " + file + " Line" + delim +
					strconv.FormatInt(int64(line), 10) + " " +
					l.labelDebug() +
					delim +
					fmt.Sprint(args...) + "\n",
			),
		)

		return
	}

	l.localWriter.Write(
		[]byte(
			l.withTimestamp() +
				" " +
				l.labelDebug() +
				delim +
				fmt.Sprint(args...) + "\n",
		),
	)
}

func (l Logger) Debugf(format string, args ...any) {
	if l.logLevel < LevelDEBUG {
		return
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		l.localWriter.Write(
			[]byte(
				l.withTimestamp() +
					" " + file + " Line" + delim +
					strconv.FormatInt(int64(line), 10) + " " +
					l.labelDebug() +
					delim +
					fmt.Sprintf(format, args...) + "\n",
			),
		)

		return
	}

	l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " +
				l.labelDebug() +
				delim +
				fmt.Sprintf(format, args...) + "\n",
		),
	)
}

func (l Logger) Debugw(msg string, keysAndValues ...any) {
	if l.logLevel < LevelDEBUG {
		return
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		l.localWriter.Write(
			[]byte(
				l.withTimestamp() +
					" " + msg + "\n" + file + " Line" + delim +
					strconv.FormatInt(int64(line), 10) + " " +
					l.labelDebug() +
					delim +
					fmt.Sprint(keysAndValues...) + "\n",
			),
		)

		return
	}

	l.localWriter.Write(
		[]byte(
			l.withTimestamp() +
				" " + msg + "\n" +
				l.labelDebug() +
				delim +
				fmt.Sprint(keysAndValues...) + "\n",
		),
	)
}
