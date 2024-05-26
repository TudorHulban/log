package log

import (
	"fmt"
)

func (l *Logger) PrintMessage(msg string) {
	l.localWriter.Write(
		[]byte(
			l.withl.withTimestamp() + " " + msg + "\n",
		),
	)
}

func (l *Logger) Print(args ...any) {
	l.localWriter.Write(
		[]byte(
			l.withl.withTimestamp() + " " + fmt.Sprint(args...) + "\n",
		),
	)
}

func (l *Logger) Printf(format string, args ...any) {
	l.localWriter.Write(
		[]byte(
			l.withTimestamp() + " " + fmt.Sprintf(format, args...) + "\n",
		),
	)
}
