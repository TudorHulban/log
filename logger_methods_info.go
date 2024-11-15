package log

import (
	"fmt"
	"runtime"
	"strconv"
)

func (l Logger) labelInfo() string {
	return ternary(
		l.withColor,

		colorInfo(logLevels[LevelINFO]),
		logLevels[LevelINFO],
	)
}

func (l Logger) Info(args ...any) {
	if l.logLevel == LevelNONE {
		return
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		_, _ = l.localWriter.Write(
			ternary(
				l.withJSON,

				jsonWCaller(
					&paramsJSONWCaller{
						timestamp: l.withTimestamp(),
						file:      file,
						line:      line,
						level:     l.labelInfo(),
						message:   fmt.Sprint(args...),
					},
				),

				[]byte(
					l.withTimestamp()+" "+file+" Line"+delim+
						strconv.FormatInt(int64(line), 10)+" "+
						l.labelInfo()+
						delim+
						fmt.Sprint(args...)+"\n",
				),
			),
		)

		return
	}

	_, _ = l.localWriter.Write(
		ternary(
			l.withJSON,

			json(
				&paramsJSONWCaller{
					timestamp: l.withTimestamp(),
					level:     l.labelInfo(),
					message:   fmt.Sprint(args...),
				},
			),

			[]byte(
				l.withTimestamp()+" "+
					l.labelInfo()+
					delim+
					fmt.Sprint(args...)+"\n",
			),
		),
	)
}

func (l Logger) Infof(format string, args ...any) {
	if l.logLevel == LevelNONE {
		return
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		_, _ = l.localWriter.Write(
			ternary(
				l.withJSON,

				jsonWCaller(
					&paramsJSONWCaller{
						timestamp: l.withTimestamp(),
						file:      file,
						line:      line,
						level:     l.labelInfo(),
						message:   fmt.Sprintf(format, args...),
					},
				),

				[]byte(
					l.withTimestamp()+" "+file+" Line"+delim+
						strconv.FormatInt(int64(line), 10)+" "+
						l.labelInfo()+
						delim+
						fmt.Sprintf(format, args...)+"\n",
				),
			),
		)

		return
	}

	_, _ = l.localWriter.Write(
		ternary(
			l.withJSON,

			json(
				&paramsJSONWCaller{
					timestamp: l.withTimestamp(),
					level:     l.labelInfo(),
					message:   fmt.Sprintf(format, args...),
				},
			),

			[]byte(
				l.withTimestamp()+" "+
					l.labelInfo()+
					delim+
					fmt.Sprintf(format, args...)+"\n",
			),
		),
	)
}

func (l Logger) Infow(msg string, keysAndValues ...any) {
	if l.logLevel == LevelNONE {
		return
	}

	if l.withCaller {
		_, file, line, _ := runtime.Caller(1)

		_, _ = l.localWriter.Write(
			[]byte(
				l.withTimestamp() + " " + msg + "\n" + file + " Line" + delim +
					strconv.FormatInt(int64(line), 10) + " " +
					logLevels[1] + delim +
					fmt.Sprint(keysAndValues...) + "\n",
			),
		)

		return
	}

	_, _ = l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " +
				logLevels[1] +
				delim +
				fmt.Sprint(keysAndValues...) + "\n",
		),
	)
}
