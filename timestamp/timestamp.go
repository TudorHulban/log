package timestamp

import (
	"strconv"
	"time"
)

// For no timestamp do not add a timestamp function.
type Timestamp func(appendTo []byte) []byte

// TimestampNano provides true nanosecond‑accurate timestamps.
// On Linux time.Now() costs ~40–70 ns by itself.
// UnixNano() + AppendInt adds ~10–15 ns.
//
// Due to above, cost is around 150 ns.
func TimestampNano(appendTo []byte) []byte {
	return strconv.AppendInt(appendTo, time.Now().UnixNano(), 10)
}

func TimestampStandard(appendTo []byte) []byte {
	updateStandardTimeCache()

	buf := timeCacheStandard.active.Load()

	return append(appendTo, buf.output[:buf.length]...)
}

func TimestampYYYYMonth(appendTo []byte) []byte {
	updateYYYYMonthTimeCache()

	buf := timeCacheYYYYMonth.active.Load()

	return append(appendTo, buf.output[:buf.length]...)
}

func TimestampRFC3339(appendTo []byte) []byte {
	return append(
		appendTo,
		[]byte(time.Now().UTC().Format(time.RFC3339))...,
	)
}

func TimestampRFC3339Nano(appendTo []byte) []byte {
	return append(
		appendTo,
		[]byte(time.Now().UTC().Format(time.RFC3339Nano))...,
	)
}
