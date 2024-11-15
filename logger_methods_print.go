package log

import (
	"fmt"
)

func (l *Logger) PrintMessage(msg string) {
	_, _ = l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " + msg + "\n",
		),
	)
}

func (l *Logger) Print(args ...any) {
	_, _ = l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " +
				fmt.Sprint(args...) + "\n",
		),
	)
}

func (l *Logger) Printw(msg string, args ...any) {
	_, _ = l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " + msg + "\n" +
				fmt.Sprint(args...) + "\n",
		),
	)
}

func (l *Logger) Printf(format string, args ...any) {
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
					fmt.Sprintf(format, args...)+"\n",
			),
		),
	)
}
