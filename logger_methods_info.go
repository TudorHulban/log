package log

import (
	"fmt"
	"runtime"
	"strconv"
)

func (l Logger) Info(args ...any) {
	if l.logLevel == LevelNONE {
		return
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		l.localWriter.Write(
			[]byte(
				l.withTimestamp() + " " + file + " Line" + delim +
					strconv.FormatInt(int64(line), 10) + " " +
					logLevels[1] + delim +
					fmt.Sprint(args...) + "\n",
			),
		)

		return
	}

	l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " + logLevels[1] + delim +
				fmt.Sprint(args...) + "\n",
		),
	)
}

func (l Logger) Infof(format string, args ...any) {
	if l.logLevel == LevelNONE {
		return
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		l.localWriter.Write(
			[]byte(
				l.withTimestamp() + " " + file + " Line" + delim +
					strconv.FormatInt(int64(line), 10) + " " +
					logLevels[1] + delim +
					fmt.Sprintf(format, args...) + "\n",
			),
		)

		return
	}

	l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " + logLevels[1] + delim +
				fmt.Sprintf(format, args...) + "\n",
		),
	)
}

func (l Logger) Infow(msg string, keysAndValues ...any) {
	if l.logLevel == LevelNONE {
		return
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		l.localWriter.Write(
			[]byte(
				l.withTimestamp() + " " + msg + "\n" + file + " Line" + delim +
					strconv.FormatInt(int64(line), 10) + " " +
					logLevels[1] + delim +
					fmt.Sprint(keysAndValues...) + "\n",
			),
		)

		return
	}

	l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " + logLevels[1] + delim +
				fmt.Sprint(keysAndValues...) + "\n",
		),
	)
}
