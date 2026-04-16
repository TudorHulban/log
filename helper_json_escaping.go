package log

import "unsafe"

// No per‑byte append.
//
// No bounds checks.
//
// No slice growth.
//
// Writes directly into the backing array.
//
// No worst‑case over‑allocation.
func appendJSONString(buf []byte, input string) []byte {
	writeIndex := len(buf)

	// Grow capacity only, without changing length
	if cap(buf)-len(buf) < 8 {
		newBuf := make([]byte, len(buf), len(buf)+len(input)+8)

		copy(newBuf, buf)
		buf = newBuf
	}

	backingArray := unsafe.Slice(unsafe.SliceData(buf), cap(buf))

	for ix := 0; ix < len(input); ix++ {
		c := input[ix]

		if c != '\\' && c != '"' && c != '\n' && c != '\r' && c != '\t' {
			if writeIndex == cap(buf) {
				newCap := cap(buf) + 64
				newBuf := make([]byte, writeIndex, newCap)
				copy(newBuf, backingArray[:writeIndex])
				buf = newBuf
				backingArray = unsafe.Slice(unsafe.SliceData(buf), cap(buf))
			}

			backingArray[writeIndex] = c
			writeIndex++

			continue
		}

		if writeIndex+2 > cap(buf) {
			newCap := cap(buf) + 64
			newBuf := make([]byte, writeIndex, newCap)

			copy(newBuf, backingArray[:writeIndex])
			buf = newBuf
			backingArray = unsafe.Slice(unsafe.SliceData(buf), cap(buf))
		}

		switch c {
		case '\\', '"':
			backingArray[writeIndex] = '\\'
			backingArray[writeIndex+1] = c
		case '\n':
			backingArray[writeIndex] = '\\'
			backingArray[writeIndex+1] = 'n'
		case '\r':
			backingArray[writeIndex] = '\\'
			backingArray[writeIndex+1] = 'r'
		case '\t':
			backingArray[writeIndex] = '\\'
			backingArray[writeIndex+1] = 't'
		}

		writeIndex += 2
	}

	return buf[:writeIndex]
}
