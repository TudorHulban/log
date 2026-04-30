package helpers

import (
	"fmt"
)

func GetEstimatedMessageSize(format string, args []any) int {
	size := 0
	ai := 0
	flen := len(format)

	for i := 0; i < flen; i++ {
		c := format[i]

		if c != '%' {
			size++

			continue
		}

		// "%%"
		if i+1 < flen && format[i+1] == '%' {
			size++
			i++

			continue
		}

		i++
		if i >= flen {
			size++ // lone '%'
			break
		}

		if ai >= len(args) {
			size = size + 2 // "%x"

			continue
		}

		arg := args[ai]
		ai++

		switch format[i] {
		case 's':
			switch v := arg.(type) {
			case string:
				size = size + len(v)
			case []byte:
				size = size + len(v)
			default:
				size = size + len(fmt.Sprint(v))
			}

		case 'd':
			switch v := arg.(type) {
			case int:
				size = size + DigitsInt(v)
			case int64:
				size = size + digitsInt64(v)
			case int32:
				size = size + digitsInt64(int64(v))
			case uint:
				size = size + digitsUint64(uint64(v))
			case uint64:
				size = size + digitsUint64(v)
			default:
				size = size + len(fmt.Sprint(v))
			}

		case 't':
			size = size + 5 // true/false worst case

		case 'f':
			switch v := arg.(type) {
			case float64:
				size = size + float64Len(v)
			case float32:
				size = size + float64Len(float64(v))
			default:
				size = size + len(fmt.Sprint(v))
			}

		case 'v':
			size = size + len(fmt.Sprint(arg))

		default:
			size = size + 2 // "%x"
		}
	}

	return size
}

func DigitsInt(v int) int {
	if v == 0 {
		return 1
	}

	n := 0
	if v < 0 {
		n = 1
		v = -v
	}

	for v > 0 {
		v /= 10
		n++
	}

	return n
}

func digitsInt64(v int64) int {
	if v == 0 {
		return 1
	}

	n := 0
	if v < 0 {
		n = 1
		v = -v
	}

	for v > 0 {
		v /= 10
		n++
	}

	return n
}

func digitsUint64(v uint64) int {
	if v == 0 {
		return 1
	}

	n := 0

	for v > 0 {
		v /= 10
		n++
	}

	return n
}

func float64Len(v float64) int {
	// conservative upper bound (cheap + safe for sizing)
	// "-12345.6789"
	n := 0
	if v < 0 {
		n = 1
		v = -v
	}

	// integer part
	i := uint64(v)
	n += digitsUint64(i)

	// decimal part (assume ".+" even if not always present)
	return n + 6
}
