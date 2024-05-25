package log

import (
	"fmt"
)

func renderPrint(msg string) string {
	return timestamp() + " " + msg + "\n"
}

func (l *Logger) PrintMessage(msg string) {
	l.write(
		[]byte(
			renderPrint(msg),
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
