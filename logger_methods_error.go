package log

import (
	"fmt"
	"runtime"
	"strconv"
)

func (l Logger) Error(args ...any) {
	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		l.localWriter.Write(
			[]byte(
				l.withTimestamp() + " " + file + " Line" + delim +
					strconv.FormatInt(int64(line), 10) + " " +
					colorDebug()(logLevels[4]) + delim +
					fmt.Sprint(args...) + "\n",
			),
		)

		return
	}

	l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " + colorDebug()(logLevels[4]) + delim +
				fmt.Sprint(args...) + "\n",
		),
	)
}

func (l Logger) Errorf(format string, args ...any) {
	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		l.localWriter.Write(
			[]byte(
				l.withTimestamp() + " " + file + " Line" + delim +
					strconv.FormatInt(int64(line), 10) + " " +
					colorDebug()(logLevels[4]) + delim +
					fmt.Sprintf(format, args...) + "\n",
			),
		)

		return
	}

	l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " + colorDebug()(logLevels[4]) + delim +
				fmt.Sprintf(format, args...) + "\n",
		),
	)
}

func (l Logger) Errorw(format string, args ...any) {
	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		l.localWriter.Write(
			[]byte(
				l.withTimestamp() + " " + file + " Line" + delim +
					strconv.FormatInt(int64(line), 10) + " " +
					colorDebug()(logLevels[4]) + delim +
					fmt.Sprintf(format, args...) + "\n",
			),
		)

		return
	}

	l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " + colorDebug()(logLevels[4]) + delim +
				fmt.Sprintf(format, args...) + "\n",
		),
	)
}
