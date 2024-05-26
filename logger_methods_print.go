package log

import (
	"fmt"
)

func (l *Logger) PrintMessage(msg string) {
	l.write(
		[]byte(
			timestamp() + " " + msg + "\n",
		),
	)
}

func (l *Logger) PrintLocal(msg string) {
	l.localWriter.Write(
		[]byte(
			timestamp() + " " + msg + "\n",
		),
	)
}

func (l *Logger) Print(args ...any) {
	l.write(
		[]byte(
			timestamp() + " " + fmt.Sprint(args...) + "\n",
		),
	)
}

func (l *Logger) Printf(format string, args ...any) {
	l.write(
		[]byte(
			timestamp() + " " + fmt.Sprintf(format, args...) + "\n",
		),
	)
}
