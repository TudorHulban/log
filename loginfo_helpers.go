package log

import (
	"bytes"
	"fmt"
	"strconv"
	"sync"
	"time"
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

func colorWarn() func(word string) string {
	return func(word string) string {
		buf := bufPool.Get().(*bytes.Buffer)
		defer bufPool.Put(buf)

		buf.Reset()

		buf.WriteString("\x1b[1;33m")
		return fmt.Sprintf("%s%v\x1b[0m", buf.String(), word)
	}
}

func colorDebug() func(word string) string {
	return func(word string) string {
		buf := bufPool.Get().(*bytes.Buffer)
		defer bufPool.Put(buf)

		buf.Reset()

		buf.WriteString("\x1b[1;34m")
		return fmt.Sprintf("%s%v\x1b[0m", buf.String(), word)
	}
}

func convertLevel(level int) int {
	switch {
	case level < 1:
		return 0
	case level == 1:
		return 1
	case level == 2:
		return 2
	default:
		return 3
	}
}

// timestamp provides time in format YYYYMonth HH24:Minutes:Seconds.Miliseconds
func timestamp() string {
	now := time.Now()

	theMonth := "0" + strconv.FormatInt(int64(now.Month()), 10)
	theMonth = theMonth[len(theMonth)-2:]

	theHour := "0" + strconv.FormatInt(int64(now.Hour()), 10)
	theHour = theHour[len(theHour)-2:]

	theMin := "0" + strconv.FormatInt(int64(now.Minute()), 10)
	theMin = theMin[len(theMin)-2:]

	theSec := "0" + strconv.FormatInt(int64(now.Second()), 10)
	theSec = theSec[len(theSec)-2:]

	theMilisec := "00" + strconv.FormatInt(int64(now.Nanosecond()/1000000), 10)
	theMilisec = theMilisec[len(theMilisec)-3:]

	return strconv.FormatInt(int64(now.Year()), 10) + theMonth + " " + theHour + ":" + theMin + ":" + theSec + "." + theMilisec
}
