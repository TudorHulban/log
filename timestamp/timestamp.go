package timestamp

import (
	"strconv"
	"time"
)

type Timestamp func(buf []byte) []byte

var TimestampNil = func(buf []byte) []byte {
	return nil
}

func TimestampNano(buf []byte) []byte {
	return strconv.AppendInt(buf[:0], time.Now().UnixNano(), 10)
}

func TimestampNano2() []byte {
	// 20 bytes is enough for any int64 in base 10
	var b [20]byte
	n := time.Now().UnixNano()

	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}

	return b[i:]
}

func TimestampYYYYMonth(buf []byte) string {
	now := time.Now()

	month := "0" + strconv.Itoa(
		int(now.Month()),
	)

	hour := "0" + strconv.Itoa(
		(now.Hour()),
	)

	minute := "0" + strconv.Itoa(
		(now.Minute()),
	)

	second := "0" + strconv.Itoa(
		(now.Second()),
	)

	millisecond := "00" + strconv.Itoa(
		(now.Nanosecond() / 1000000), // with miliseconds performance was worse
	)

	return strconv.Itoa(now.Year()) +
		month[len(month)-2:] + " " +
		hour[len(hour)-2:] + ":" +
		minute[len(minute)-2:] + ":" +
		second[len(second)-2:] + "." +
		millisecond[len(millisecond)-3:]
}
