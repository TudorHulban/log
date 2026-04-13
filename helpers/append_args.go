package helpers

import (
	"fmt"
	"strconv"
)

func AppendInt(destination []byte, value int) []byte {
	// Enough space for int64 decimal representation
	var buf [20]byte

	i := len(buf)

	neg := value < 0
	if neg {
		value = -value
	}

	for {
		i--
		buf[i] = byte('0' + value%10)

		value = value / 10
		if value == 0 {
			break
		}
	}

	if neg {
		i--
		buf[i] = '-'
	}

	return append(destination, buf[i:]...)
}

func AppendUint64(destination []byte, value uint64) []byte {
	var buf [20]byte

	i := len(buf)

	for {
		i--
		buf[i] = byte('0' + value%10)

		value = value / 10
		if value == 0 {
			break
		}
	}

	return append(destination, buf[i:]...)
}

func AppendFloat(destination []byte, value float64, prec int) []byte {
	// Handle NaN and Inf explicitly
	if value != value {
		return append(destination, 'n', 'a', 'n')
	}

	if value > 1e308 {
		return append(destination, 'i', 'n', 'f')
	}

	if value < -1e308 {
		return append(destination, '-', 'i', 'n', 'f')
	}

	// Sign
	if value < 0 {
		destination = append(destination, '-')
		value = -value
	}

	// Integer part
	intPart := uint64(value)
	destination = AppendUint64(destination, intPart)

	// Fractional part
	if prec > 0 {
		destination = append(destination, '.')
		fractional := value - float64(intPart)

		for range prec {
			fractional *= 10
			digit := uint64(fractional)
			destination = append(destination, byte('0'+digit))
			fractional -= float64(digit)
		}
	}

	return destination
}

func AppendBool(destination []byte, value bool) []byte {
	if value {
		return append(destination, 't', 'r', 'u', 'e')
	}

	return append(destination, 'f', 'a', 'l', 's', 'e')
}

func AppendError(destination []byte, err error) []byte {
	if err == nil {
		return append(destination, 'n', 'i', 'l')
	}

	return append(destination, err.Error()...)
}

// AppendArgs appends each arg to dst without reflection or fmt.
// Covers the types that appear in practice. Falls back to fmt only for
// exotic types — still no alloc on the hot path.
func AppendArgs(destination []byte, args []any) []byte {
	for ix, arg := range args {
		if ix > 0 {
			destination = append(destination, ' ')
		}

		switch value := arg.(type) {
		case string:
			destination = append(destination, value...)
		case []byte:
			destination = append(destination, value...)
		case int:
			destination = AppendInt(destination, value)
		case int64:
			destination = strconv.AppendInt(destination, value, 10)
		case int32:
			destination = strconv.AppendInt(destination, int64(value), 10)
		case uint:
			destination = strconv.AppendUint(destination, uint64(value), 10)
		case uint64:
			destination = AppendUint64(destination, value)
		case float64:
			destination = strconv.AppendFloat(destination, value, 'f', -1, 64)
		case float32:
			destination = strconv.AppendFloat(destination, float64(value), 'f', -1, 32)
		case bool:
			destination = AppendBool(destination, value)
		case error:
			destination = AppendError(destination, value)
		case nil:
			destination = append(destination, "null"...)

		default:
			// Exotic types only — keeps hot path clean.
			destination = append(destination, fmt.Sprint(value)...)
		}
	}

	return destination
}

func toString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case []byte:
		return string(x)
	default:
		return fmt.Sprint(x)
	}
}

func ArgsToString(args []any) string {
	// Fast path: no args
	if len(args) == 0 {
		return ""
	}

	// We build into a byte slice and convert once at the end.
	// This avoids multiple string concatenations.
	buf := make([]byte, 0, 64) // small initial cap; grows if needed

	for i := range args {
		if i > 0 {
			buf = append(buf, ' ')
		}

		switch v := args[i].(type) {
		case string:
			buf = append(buf, v...)

		case int:
			buf = AppendInt(buf, v)

		case bool:
			buf = AppendBool(buf, v)

		case []byte:
			buf = append(buf, v...)

		default:
			// Fallback: convert to string once
			buf = append(buf, []byte(toString(v))...)
		}
	}

	return string(buf)
}

func Appendf(dst []byte, format string, args ...any) []byte {
	ai := 0
	flen := len(format)

	for i := 0; i < flen; i++ {
		c := format[i]

		if c != '%' {
			dst = append(dst, c)
			continue
		}

		// Handle "%%"
		if i+1 < flen && format[i+1] == '%' {
			dst = append(dst, '%')
			i++
			continue
		}

		// Move to verb
		i++
		if i >= flen {
			// malformed trailing '%'
			dst = append(dst, '%')
			break
		}

		if ai >= len(args) {
			// no argument provided
			dst = append(dst, '%', format[i])
			continue
		}

		arg := args[ai]
		ai++

		switch format[i] {
		case 's':
			switch v := arg.(type) {
			case string:
				dst = append(dst, v...)
			case []byte:
				dst = append(dst, v...)
			default:
				dst = append(dst, fmt.Sprint(v)...) // exotic fallback
			}

		case 'd':
			switch v := arg.(type) {
			case int:
				dst = AppendInt(dst, v)
			case int64:
				dst = AppendInt(dst, int(v))
			case int32:
				dst = AppendInt(dst, int(v))
			case uint:
				dst = AppendUint64(dst, uint64(v))
			case uint64:
				dst = AppendUint64(dst, v)
			default:
				dst = append(dst, fmt.Sprint(v)...)
			}

		case 'v':
			// direct passthrough to your existing fast path
			dst = AppendArgs(dst, []any{arg})

		case 't':
			switch v := arg.(type) {
			case bool:
				dst = AppendBool(dst, v)
			default:
				dst = append(dst, fmt.Sprint(v)...)
			}

		case 'f':
			switch v := arg.(type) {
			case float64:
				dst = strconv.AppendFloat(dst, v, 'f', -1, 64)
			case float32:
				dst = strconv.AppendFloat(dst, float64(v), 'f', -1, 32)
			default:
				dst = append(dst, fmt.Sprint(v)...)
			}

		default:
			// unsupported verb → literal fallback
			dst = append(dst, '%', format[i])
		}
	}

	return dst
}
