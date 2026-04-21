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

	// Pre-allocate capacity if needed for initial buffer
	if cap(buf)-len(buf) < 8 {
		newBuf := make([]byte, len(buf), len(buf)+len(input)+8)

		copy(newBuf, buf)
		buf = newBuf
	}

	backingArray := unsafe.Slice(unsafe.SliceData(buf), cap(buf))

	for i := 0; i < len(input); i++ {
		c := input[i]

		switch {
		// Short escapes for common cases (2 bytes each)
		case c == '\\' || c == '"':
			if writeIndex+2 > cap(buf) {
				newCap := cap(buf) + 64
				newBuf := make([]byte, writeIndex, newCap)
				copy(newBuf, backingArray[:writeIndex])
				buf = newBuf
				backingArray = unsafe.Slice(unsafe.SliceData(buf), cap(buf))
			}

			backingArray[writeIndex] = '\\'
			backingArray[writeIndex+1] = c
			writeIndex = writeIndex + 2

		case c == '\n':
			if writeIndex+2 > cap(buf) {
				newCap := cap(buf) + 64
				newBuf := make([]byte, writeIndex, newCap)
				copy(newBuf, backingArray[:writeIndex])
				buf = newBuf
				backingArray = unsafe.Slice(unsafe.SliceData(buf), cap(buf))
			}

			backingArray[writeIndex] = '\\'
			backingArray[writeIndex+1] = 'n'
			writeIndex = writeIndex + 2

		case c == '\r':
			if writeIndex+2 > cap(buf) {
				newCap := cap(buf) + 64
				newBuf := make([]byte, writeIndex, newCap)
				copy(newBuf, backingArray[:writeIndex])
				buf = newBuf
				backingArray = unsafe.Slice(unsafe.SliceData(buf), cap(buf))
			}

			backingArray[writeIndex] = '\\'
			backingArray[writeIndex+1] = 'r'
			writeIndex = writeIndex + 2

		case c == '\t':
			if writeIndex+2 > cap(buf) {
				newCap := cap(buf) + 64
				newBuf := make([]byte, writeIndex, newCap)
				copy(newBuf, backingArray[:writeIndex])
				buf = newBuf
				backingArray = unsafe.Slice(unsafe.SliceData(buf), cap(buf))
			}

			backingArray[writeIndex] = '\\'
			backingArray[writeIndex+1] = 't'
			writeIndex = writeIndex + 2

		// RFC 8259 §7: other control chars (0x00–0x1F) MUST be \uXXXX
		case c < 0x20:
			if writeIndex+6 > cap(buf) {
				newCap := cap(buf) + 64
				newBuf := make([]byte, writeIndex, newCap)

				copy(newBuf, backingArray[:writeIndex])
				buf = newBuf
				backingArray = unsafe.Slice(unsafe.SliceData(buf), cap(buf))
			}

			backingArray[writeIndex] = '\\'
			backingArray[writeIndex+1] = 'u'
			backingArray[writeIndex+2] = '0'
			backingArray[writeIndex+3] = '0'
			backingArray[writeIndex+4] = '0' + (c >> 4) // c>>4 ∈ {0,1}

			lo := c & 0xF

			if lo < 10 {
				backingArray[writeIndex+5] = '0' + lo
			} else {
				backingArray[writeIndex+5] = 'a' + lo - 10
			}

			writeIndex = writeIndex + 6

		// Regular character - no escaping
		default:
			if writeIndex == cap(buf) {
				newCap := cap(buf) + 64
				newBuf := make([]byte, writeIndex, newCap)

				copy(newBuf, backingArray[:writeIndex])
				buf = newBuf
				backingArray = unsafe.Slice(unsafe.SliceData(buf), cap(buf))
			}

			backingArray[writeIndex] = c
			writeIndex++
		}
	}

	return buf[:writeIndex]
}
