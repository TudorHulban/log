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
	updateRFC3339TimeCache() // Fast path: update cache (cheap CAS if already fresh)

	// Load cached value (non-nil after update)
	buf := timeCacheStandard.active.Load()
	if buf != nil {
		return append(appendTo, buf.output[:buf.length]...)
	}

	// Fallback: should almost never happen (cache cold start)
	return time.Now().UTC().AppendFormat(appendTo, rfc3339MilliLayout)
}

func TimestampRFC3339Nano(appendTo []byte) []byte {
	return time.Now().UTC().AppendFormat(appendTo, time.RFC3339Nano)
}

func TimestampRFC3339Bucharest(appendTo []byte) []byte {
	loc, _ := time.LoadLocation("Europe/Bucharest")

	return time.Now().In(loc).AppendFormat(appendTo, time.RFC3339)
}
