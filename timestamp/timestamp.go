package timestamp

import (
	"strconv"
	"time"
)

type Timestamp func() string

var TimestampNil = func() string {
	return ""
}

var TimestampNano = func() string {
	return strconv.Itoa(
		int(time.Now().UnixNano()),
	)
}

var TimestampYYYYMonth = func() string {
	now := time.Now()

	month := "0" + strconv.Itoa(
		int(now.Month()),
	)

	hour := "0" + strconv.Itoa(
		int(now.Hour()),
	)

	minute := "0" + strconv.Itoa(
		int(now.Minute()),
	)

	second := "0" + strconv.Itoa(
		int(now.Second()),
	)

	milisecond := "00" + strconv.Itoa(
		int(now.Nanosecond()/1000000), // with miliseconds performance was worse
	)

	return strconv.Itoa(int(now.Year())) +
		month[len(month)-2:] + " " +
		hour[len(hour)-2:] + ":" +
		minute[len(minute)-2:] + ":" +
		second[len(second)-2:] + "." +
		milisecond[len(milisecond)-3:]
}
