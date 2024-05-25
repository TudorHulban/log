package log

import (
	"fmt"
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

func colorWarn() func(word string) string {
	return func(word string) string {
		return fmt.Sprintf(
			"%s%v\x1b[0m",
			"\x1b[1;33m",
			word,
		)
	}
}

func colorDebug() func(word string) string {
	return func(word string) string {
		return fmt.Sprintf(
			"%s%v\x1b[0m",
			"\x1b[1;34m",
			word,
		)
	}
}
