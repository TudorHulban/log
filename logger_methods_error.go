package log

import (
	"fmt"
	"runtime"
	"strconv"
)

func (l Logger) labelError() string {
	return ternary(
		l.withColor,

		colorError(logLevels[LevelERROR]),
		logLevels[LevelERROR],
	)
}

func (l Logger) Error(args ...any) {
	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		_, _ = l.localWriter.Write(
			[]byte(
				l.withTimestamp() + " " + file + " Line" + delim +
					strconv.FormatInt(int64(line), 10) + " " +
					l.labelError() +
					delim +
					fmt.Sprint(args...) + "\n",
			),
		)

		return
	}

	_, _ = l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " +
				l.labelError() +
				delim +
				fmt.Sprint(args...) + "\n",
		),
	)
}

func (l Logger) Errorf(format string, args ...any) {
	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		_, _ = l.localWriter.Write(
			[]byte(
				l.withTimestamp() + " " + file + " Line" + delim +
					strconv.FormatInt(int64(line), 10) + " " +
					l.labelError() +
					delim +
					fmt.Sprintf(format, args...) + "\n",
			),
		)

		return
	}

	_, _ = l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " +
				l.labelError() +
				delim +
				fmt.Sprintf(format, args...) + "\n",
		),
	)
}

func (l Logger) Errorw(format string, args ...any) {
	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		_, _ = l.localWriter.Write(
			[]byte(
				l.withTimestamp() + " " + file + " Line" + delim +
					strconv.FormatInt(int64(line), 10) + " " +
					l.labelError() +
					delim +
					fmt.Sprintf(format, args...) + "\n",
			),
		)

		return
	}

	_, _ = l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " +
				l.labelError() +
				delim +
				fmt.Sprintf(format, args...) + "\n",
		),
	)
}
