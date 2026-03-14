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

func TimestampStandard(buf []byte) []byte {
	updateTimeCache()

	return append(buf, tc.stdBuf[:tc.stdLen]...)
}

func TimestampYYYYMonth(buf []byte) []byte {
	updateTimeCache()

	return append(buf, tc.customBuf[:tc.customLen]...)
}
