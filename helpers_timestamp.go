package log

import (
	"strconv"
	"time"
)

// time format: YYYYMonth HH24:Minutes:Seconds.Miliseconds
func timestamp() string {
	now := time.Now()

	theMonth := "0" + strconv.Itoa(
		int(now.Month()),
	)

	theHour := "0" + strconv.Itoa(
		int(now.Hour()),
	)

	theMin := "0" + strconv.Itoa(
		int(now.Minute()),
	)

	theSec := "0" + strconv.Itoa(
		int(now.Second()),
	)

	theMilisec := "00" + strconv.Itoa(
		int(now.Nanosecond()/1000000),
	)

	return strconv.Itoa(int(now.Year())) +
		theMonth[len(theMonth)-2:] + " " +
		theHour[len(theHour)-2:] + ":" +
		theMin[len(theMin)-2:] + ":" +
		theSec[len(theSec)-2:] + "." +
		theMilisec[len(theMilisec)-3:]
}
