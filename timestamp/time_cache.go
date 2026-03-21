package timestamp

import (
	"strconv"
	"sync/atomic"
	"time"
)

const nanosPerDay = int64(24 * 60 * 60 * 1e9)

type timeBuf struct {
	valueMillisecond int64 // millisecond-epoch at which this buf was built
	valueDay         int64 // Unix-day at which the date prefix was built

	length  int // total byte length of the formatted timestamp
	dateLen int // byte length of the date prefix inside data

	output [32]byte // the formatted timestamp bytes
}

type timeCache struct {
	active atomic.Pointer[timeBuf]
}

var timeCacheStandard timeCache
var timeCacheYYYYMonth timeCache

// gateStandard and gateYYYYMonth serialize the one writer per millisecond.
// When multiple goroutines observe a stale valueMillisecond simultaneously, exactly one
// wins the CAS and does the update. The rest load the freshly stored pointer.
var (
	gateStandard  atomic.Int64
	gateYYYYMonth atomic.Int64
)

func updateStandardTimeCache() {
	now := time.Now()
	nowMillisecond := now.UnixNano() / 1e6

	// update timestamp every millisecond. TTL = 1 millisecond.
	if !gateStandard.CompareAndSwap(gateStandard.Load(), nowMillisecond) {
		return
	}

	previous := timeCacheStandard.active.Load()

	next := new(timeBuf)
	next.valueMillisecond = nowMillisecond

	if previous != nil {
		*next = *previous
	}

	// rebuild date prefix only when the day changes.
	nowDay := now.UnixNano() / nanosPerDay

	if previous == nil || nowDay != previous.valueDay {
		next.valueDay = nowDay

		year, month, day := now.Date()

		std := next.output[:0]

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

		next.dateLen = len(std)
	}

	hour, minute, sec := now.Clock()
	milli := now.Nanosecond() / 1e6

	std := next.output[next.dateLen:next.dateLen]

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

	next.length = next.dateLen + len(std)

	timeCacheStandard.active.Store(next)
}

func updateYYYYMonthTimeCache() {
	now := time.Now()
	nowMs := now.UnixNano() / 1e6

	// update timestamp every millisecond. TTL = 1 millisecond.
	if !gateYYYYMonth.CompareAndSwap(gateYYYYMonth.Load(), nowMs) {
		return
	}

	previous := timeCacheYYYYMonth.active.Load()

	next := new(timeBuf)
	next.valueMillisecond = nowMs

	if previous != nil {
		*next = *previous
	}

	// rebuild date prefix only when the day changes.
	nowDay := now.UnixNano() / nanosPerDay

	if previous == nil || nowDay != previous.valueDay {
		next.valueDay = nowDay

		year, month, day := now.Date()

		custom := next.output[:0]

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

		next.dateLen = len(custom)
	}

	hour, minute, sec := now.Clock()
	milli := now.Nanosecond() / 1e6

	custom := next.output[next.dateLen:next.dateLen]

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

	next.length = next.dateLen + len(custom)

	timeCacheYYYYMonth.active.Store(next)
}
