package log

import (
	"bytes"
	"fmt"
	"sync"
)

/*
Colour Styles
	//Text Colour
	Black   = 30
	Red     = 31
	Green   = 32
	Yellow  = 33
	Blue    = 34
	Magenta = 35
	Cyan    = 36
	White   = 37
	Grey    = 90

	//Style
	Bold      = 1
	Dim       = 2
	Italic    = 3
	Underline = 4
	Blinkslow = 5
	Blinkfast = 6
	Inverse   = 7
	Hidden    = 8
	Strikeout = 9
*/

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func cWarn() func(word string) string {
	return func(word string) string {
		b := new(bytes.Buffer)
		b.WriteString("\x1b[1;33m")
		return fmt.Sprintf("%s%v\x1b[0m", b.String(), word)
	}
}

func cDebug() func(word string) string {
	return func(word string) string {
		b := new(bytes.Buffer)
		b.WriteString("\x1b[1;34m")
		return fmt.Sprintf("%s%v\x1b[0m", b.String(), word)
	}
}
