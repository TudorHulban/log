package log

import (
	"fmt"
	"runtime"
	"strconv"
)

func (l Logger) labelWarn() string {
	return ternary(
		l.withColor,

		colorWarn(logLevels[LevelWARN]),
		logLevels[LevelWARN],
	)
}

func (l Logger) Warn(args ...any) {
	if l.logLevel < LevelWARN {
		return
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		_, _ = l.localWriter.Write(
			[]byte(
				l.withTimestamp() + " " + file + " Line" + delim +
					strconv.FormatInt(int64(line), 10) + " " +
					l.labelWarn() +
					delim +
					fmt.Sprint(args...) + "\n",
			),
		)

		return
	}

	_, _ = l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " +
				l.labelWarn() +
				delim +
				fmt.Sprint(args...) + "\n",
		),
	)
}

func (l Logger) Warnf(format string, args ...any) {
	if l.logLevel < LevelWARN {
		return
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		_, _ = l.localWriter.Write(
			[]byte(
				l.withTimestamp() + " " + file + " Line" + delim +
					strconv.FormatInt(int64(line), 10) + " " +
					l.labelWarn() +
					delim +
					fmt.Sprintf(format, args...) + "\n",
			),
		)

		return
	}

	_, _ = l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " +
				l.labelWarn() +
				delim +
				fmt.Sprintf(format, args...) + "\n",
		),
	)
}

func (l Logger) Warnw(msg string, keysAndValues ...any) {
	if l.logLevel < LevelWARN {
		return
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		_, _ = l.localWriter.Write(
			[]byte(
				l.withTimestamp() + " " + file + " Line" + delim +
					strconv.FormatInt(int64(line), 10) + " " +
					l.labelWarn() +
					delim +
					fmt.Sprint(keysAndValues...) + "\n",
			),
		)

		return
	}

	_, _ = l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " +
				l.labelWarn() +
				delim +
				fmt.Sprint(keysAndValues...) + "\n",
		),
	)
}
