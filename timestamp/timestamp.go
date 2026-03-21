package timestamp

import (
	"strconv"
	"time"
)

// For no timestamp do not add a timestamp function.
type Timestamp func() []byte

// TimestampNano provides true nanosecond‑accurate timestamps.
// On Linux time.Now() costs ~40–70 ns by itself.
// UnixNano() + AppendInt adds ~10–15 ns.
//
// Due to above, cost is around 150 ns.
func TimestampNano() []byte {
	return strconv.AppendInt(nil, time.Now().UnixNano(), 10)
}

func TimestampStandard() []byte {
	updateStandardTimeCache()

	return append(
		[]byte(nil),
		timeCacheStandard.buf[:timeCacheStandard.length]...,
	)
}

func TimestampYYYYMonth() []byte {
	updateYYYYMonthTimeCache()

	return append(
		[]byte(nil),
		timeCacheYYYYMonth.buf[:timeCacheYYYYMonth.length]...,
	)
}

func TimestampRFC3339() []byte {
	return []byte(time.Now().UTC().Format(time.RFC3339))
}

func TimestampRFC3339Nano() []byte {
	return []byte(time.Now().UTC().Format(time.RFC3339Nano))
}
