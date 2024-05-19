package log

import (
	"strconv"
	"time"
)

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
