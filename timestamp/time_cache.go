package timestamp

import (
	"strconv"
	"sync/atomic"
	"time"
)

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

var (
	timeCacheStandard  timeCache
	timeCacheYYYYMonth timeCache

	timeCacheRFC3339 timeCache
)

func buildRFC3339Cache(now time.Time) {
	nowNano := now.UnixNano()
	nowMs := nowNano / 1e6

	prev := timeCacheRFC3339.active.Load()
	if prev != nil && prev.valueMillisecond == nowMs {
		return // already current
	}

	next := new(timeBuf)
	next.valueMillisecond = nowMs
	nowDay := nowNano / nanosPerDay

	if prev != nil {
		*next = *prev
	}

	if prev == nil || nowDay != prev.valueDay {
		next.valueDay = nowDay
		year, month, day := now.Date()
		b := next.output[:0]
		b = strconv.AppendInt(b, int64(year), 10)

		b = append(b, '-')
		if month < 10 {
			b = append(b, '0')
		}

		b = strconv.AppendInt(b, int64(month), 10)

		b = append(b, '-')
		if day < 10 {
			b = append(b, '0')
		}

		b = strconv.AppendInt(b, int64(day), 10)
		b = append(b, 'T')
		next.dateLen = len(b)
	}

	hour, minute, sec := now.Clock()
	milli := now.Nanosecond() / 1e6
	b := next.output[next.dateLen:next.dateLen]

	if hour < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(hour), 10)

	b = append(b, ':')
	if minute < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(minute), 10)

	b = append(b, ':')
	if sec < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(sec), 10)

	b = append(b, '.')
	if milli < 100 {
		b = append(b, '0')
	}

	if milli < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(milli), 10)
	b = append(b, 'Z')

	next.length = next.dateLen + len(b)
	timeCacheRFC3339.active.Store(next)
}

func buildStandardCache(now time.Time) {
	nowNano := now.UnixNano()
	nowMs := nowNano / 1e6

	prev := timeCacheStandard.active.Load()
	if prev != nil && prev.valueMillisecond == nowMs {
		return
	}

	next := new(timeBuf)
	next.valueMillisecond = nowMs
	nowDay := nowNano / nanosPerDay

	if prev != nil {
		*next = *prev
	}

	if prev == nil || nowDay != prev.valueDay {
		next.valueDay = nowDay
		year, month, day := now.Date()
		b := next.output[:0]
		b = strconv.AppendInt(b, int64(year), 10)

		b = append(b, '/')
		if month < 10 {
			b = append(b, '0')
		}

		b = strconv.AppendInt(b, int64(month), 10)

		b = append(b, '/')
		if day < 10 {
			b = append(b, '0')
		}

		b = strconv.AppendInt(b, int64(day), 10)
		b = append(b, ' ')

		next.dateLen = len(b)
	}

	hour, minute, sec := now.Clock()
	milli := now.Nanosecond() / 1e6
	b := next.output[next.dateLen:next.dateLen]

	if hour < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(hour), 10)

	b = append(b, ':')
	if minute < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(minute), 10)

	b = append(b, ':')
	if sec < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(sec), 10)

	b = append(b, '.')
	if milli < 100 {
		b = append(b, '0')
	}

	if milli < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(milli), 10)

	next.length = next.dateLen + len(b)
	timeCacheStandard.active.Store(next)
}

func buildYYYYMonthCache(now time.Time) {
	nowNano := now.UnixNano()
	nowMs := nowNano / 1e6

	prev := timeCacheYYYYMonth.active.Load()
	if prev != nil && prev.valueMillisecond == nowMs {
		return
	}

	next := new(timeBuf)
	next.valueMillisecond = nowMs
	nowDay := nowNano / nanosPerDay

	if prev != nil {
		*next = *prev
	}

	if prev == nil || nowDay != prev.valueDay {
		next.valueDay = nowDay
		year, month, day := now.Date()
		b := next.output[:0]

		b = strconv.AppendInt(b, int64(year), 10)
		if month < 10 {
			b = append(b, '0')
		}

		b = strconv.AppendInt(b, int64(month), 10)

		b = append(b, ' ')
		if day < 10 {
			b = append(b, '0')
		}

		b = strconv.AppendInt(b, int64(day), 10)
		b = append(b, ' ')
		next.dateLen = len(b)
	}

	hour, minute, sec := now.Clock()
	milli := now.Nanosecond() / 1e6
	b := next.output[next.dateLen:next.dateLen]

	if hour < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(hour), 10)

	b = append(b, ':')
	if minute < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(minute), 10)

	b = append(b, ':')
	if sec < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(sec), 10)

	b = append(b, '.')
	if milli < 100 {
		b = append(b, '0')
	}

	if milli < 10 {
		b = append(b, '0')
	}

	b = strconv.AppendInt(b, int64(milli), 10)

	next.length = next.dateLen + len(b)
	timeCacheYYYYMonth.active.Store(next)
}
