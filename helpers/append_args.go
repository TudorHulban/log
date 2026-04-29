package helpers

import (
	"fmt"
	"strconv"
)

// AppendArgs appends each arg to dst without reflection or fmt.
// Covers the types that appear in practice. Falls back to fmt only for
// exotic types — still no alloc on the hot path.
func AppendArgs(destination []byte, args ...any) []byte {
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

// Appendf produces a []byte message without allocations.
func Appendf(dst []byte, format string, args []any) []byte {
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
			switch v := arg.(type) {
			case string:
				dst = append(dst, v...)
			case []byte:
				dst = append(dst, v...)
			case int:
				dst = AppendInt(dst, v)
			case int64:
				dst = strconv.AppendInt(dst, v, 10)
			case int32:
				dst = strconv.AppendInt(dst, int64(v), 10)
			case uint:
				dst = strconv.AppendUint(dst, uint64(v), 10)
			case uint64:
				dst = AppendUint64(dst, v)
			case float64:
				dst = strconv.AppendFloat(dst, v, 'f', -1, 64)
			case float32:
				dst = strconv.AppendFloat(dst, float64(v), 'f', -1, 32)
			case bool:
				dst = AppendBool(dst, v)
			case error:
				dst = AppendError(dst, v)
			case nil:
				dst = append(dst, "null"...)
			default:
				dst = append(dst, fmt.Sprint(v)...)
			}

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
