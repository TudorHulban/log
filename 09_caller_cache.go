package log

import (
	"runtime"
	"sync/atomic"
	"unsafe"
)

type callerCacheEntry struct {
	pc   uintptr
	file string
	line int
}

const cacheSize = 1024
const cacheMask = cacheSize - 1

var callerTable [cacheSize]unsafe.Pointer

func (l *Logger) slowPathCaller(pc uintptr, idx uintptr) (string, int) {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "", 0
	}
	file, line := fn.FileLine(pc)

	// Create new entry
	newEntry := &callerCacheEntry{
		pc:   pc,
		file: file,
		line: line,
	}

	// Atomic store (last-writer-wins, race-safe)
	atomic.StorePointer(&callerTable[idx], unsafe.Pointer(newEntry))

	return file, line
}

// file, line = l.getCallerData()

func (l *Logger) getCallerData(skip int) (file string, line int) {
	// 1. Get the PC at the specified skip level.
	// We use a small stack buffer.
	var pcs [1]uintptr

	// runtime.Callers(skip, ...) starts at the skip level.
	// If l.callerLevel is set correctly (e.g., 3 or 4),
	// it retrieves the PC of your user's code.
	n := runtime.Callers(skip, pcs[:])
	if n == 0 {
		return "", 0
	}
	pc := pcs[0]

	// 2. Hash the PC
	idx := (pc >> 4) & cacheMask
	ptr := atomic.LoadPointer(&callerTable[idx])

	if ptr != nil {
		entry := (*callerCacheEntry)(ptr)
		if entry.pc == pc {
			return entry.file, entry.line
		}
	}

	// 3. Slow Path: Resolve symbols only if PC not in cache
	return l.slowPathCaller(pc, idx)
}
