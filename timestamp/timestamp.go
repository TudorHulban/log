package timestamp

import (
	"strconv"
	"time"
)

type Timestamp func(buf []byte) []byte

var TimestampNil = func(buf []byte) []byte {
	return nil
}

// TimestampNano provides true nanosecond‑accurate timestamps.
// On Linux time.Now() costs ~40–70 ns by itself.
// UnixNano() + AppendInt adds ~10–15 ns.
//
// Due to above, cost is around 150 ns.
func TimestampNano(buf []byte) []byte {
	return strconv.AppendInt(buf[:0], time.Now().UnixNano(), 10)
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
