package loginfo

import (
	"bytes"
	"fmt"
	"strconv"
)

/*
Colour Styles
const (
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
) */

type lito struct {
	color int64
	style int64
}

func cWarn() func(word string) string {
	lito := lito{
		color: 33,
		style: 1,
	}
	return func(word string) string {
		b := new(bytes.Buffer)
		b.WriteString("\x1b[" + strconv.FormatInt(lito.style, 10) + ";" + strconv.FormatInt(lito.color, 10) + "m")
		return fmt.Sprintf("%s%v\x1b[0m", b.String(), word)
	}
}

func cDebug() func(word string) string {
	lito := lito{
		color: 34,
		style: 1,
	}
	return func(word string) string {
		b := new(bytes.Buffer)
		b.WriteString("\x1b[" + strconv.FormatInt(lito.style, 10) + ";" + strconv.FormatInt(lito.color, 10) + "m")
		return fmt.Sprintf("%s%v\x1b[0m", b.String(), word)
	}
}
