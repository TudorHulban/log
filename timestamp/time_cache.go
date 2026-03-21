package timestamp

import (
	"strconv"
	"time"
)

type timeCache struct {
	currentTimestamp int64

	buf    [32]byte
	length int
}

var timeCacheStandard timeCache
var timeCacheYYYYMonth timeCache

func updateStandardTimeCache() {
	now := time.Now()
	nowTimestamp := now.UnixNano() / 1e6

	// update timestamp every millisecond. TTL = 1 millisecond.
	if nowTimestamp == timeCacheStandard.currentTimestamp {
		return
	}

	timeCacheStandard.currentTimestamp = nowTimestamp

	year, month, day := now.Date()
	hour, minute, sec := now.Clock()
	milli := now.Nanosecond() / 1e6

	std := timeCacheStandard.buf[:0]

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

	timeCacheStandard.length = len(std)
}

func updateYYYYMonthTimeCache() {
	now := time.Now()
	nowTimestamp := now.UnixNano() / 1e6

	// update timestamp every millisecond. TTL = 1 millisecond.
	if nowTimestamp == timeCacheYYYYMonth.currentTimestamp {
		return
	}

	timeCacheStandard.currentTimestamp = nowTimestamp

	year, month, day := now.Date()
	hour, minute, sec := now.Clock()
	milli := now.Nanosecond() / 1e6

	// -----------------------------
	// CUSTOM FORMAT
	// YYYYMM DD HH:MM:SS.mmm
	// -----------------------------
	custom := timeCacheYYYYMonth.buf[:0]

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

	timeCacheYYYYMonth.length = len(custom)
}
