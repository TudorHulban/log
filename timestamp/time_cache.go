package timestamp

import (
	"strconv"
	"time"
)

type timeCache struct {
	currentTimestamp int64

	// your custom format
	customBuf [32]byte
	customLen int

	// standard library format
	stdBuf [32]byte
	stdLen int
}

var tc timeCache

func updateTimeCache() {
	now := time.Now()
	nowTimestamp := now.UnixNano() / 1e6

	if nowTimestamp == tc.currentTimestamp { // update timestamp every millisecond. TTL = 1 millisecond.
		return
	}

	tc.currentTimestamp = nowTimestamp

	year, month, day := now.Date()
	hour, minute, sec := now.Clock()
	milli := now.Nanosecond() / 1e6

	// -----------------------------
	// CUSTOM FORMAT
	// YYYYMM DD HH:MM:SS.mmm
	// -----------------------------
	custom := tc.customBuf[:0]

	// YYYY
	custom = strconv.AppendInt(custom, int64(year), 10)

	// MM
	if month < 10 {
		custom = append(custom, '0')
	}

	custom = strconv.AppendInt(custom, int64(month), 10)
	custom = append(custom, ' ')

	// DD
	if day < 10 {
		custom = append(custom, '0')
	}

	custom = strconv.AppendInt(custom, int64(day), 10)
	custom = append(custom, ' ')

	// HH
	if hour < 10 {
		custom = append(custom, '0')
	}

	custom = strconv.AppendInt(custom, int64(hour), 10)
	custom = append(custom, ':')

	// MM
	if minute < 10 {
		custom = append(custom, '0')
	}

	custom = strconv.AppendInt(custom, int64(minute), 10)
	custom = append(custom, ':')

	// SS
	if sec < 10 {
		custom = append(custom, '0')
	}

	custom = strconv.AppendInt(custom, int64(sec), 10)
	custom = append(custom, '.')

	// mmm
	if milli < 100 {
		custom = append(custom, '0')
	}

	if milli < 10 {
		custom = append(custom, '0')
	}

	custom = strconv.AppendInt(custom, int64(milli), 10)

	tc.customLen = len(custom)

	// -----------------------------
	// STANDARD FORMAT
	// YYYY/MM/DD HH:MM:SS.mmm
	// -----------------------------
	std := tc.stdBuf[:0]

	// YYYY
	std = strconv.AppendInt(std, int64(year), 10)
	std = append(std, '/')

	// MM
	if month < 10 {
		std = append(std, '0')
	}

	std = strconv.AppendInt(std, int64(month), 10)
	std = append(std, '/')

	// DD
	if day < 10 {
		std = append(std, '0')
	}

	std = strconv.AppendInt(std, int64(day), 10)
	std = append(std, ' ')

	// HH
	if hour < 10 {
		std = append(std, '0')
	}

	std = strconv.AppendInt(std, int64(hour), 10)
	std = append(std, ':')

	// MM
	if minute < 10 {
		std = append(std, '0')
	}

	std = strconv.AppendInt(std, int64(minute), 10)
	std = append(std, ':')

	// SS
	if sec < 10 {
		std = append(std, '0')
	}

	std = strconv.AppendInt(std, int64(sec), 10)
	std = append(std, '.')

	// mmm
	if milli < 100 {
		std = append(std, '0')
	}

	if milli < 10 {
		std = append(std, '0')
	}

	std = strconv.AppendInt(std, int64(milli), 10)

	tc.stdLen = len(std)
}
