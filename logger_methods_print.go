package log

import (
	"fmt"
)

func (l *Logger) PrintMessage(msg string) {
	l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " + msg + "\n",
		),
	)
}

func (l *Logger) Print(args ...any) {
	l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " +
				fmt.Sprint(args...) + "\n",
		),
	)
}

func (l *Logger) Printw(msg string, args ...any) {
	l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " + msg + "\n" +
				fmt.Sprint(args...) + "\n",
		),
	)
}

func (l *Logger) Printf(format string, args ...any) {
	l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " +
				fmt.Sprintf(format, args...) + "\n",
		),
	)
}
