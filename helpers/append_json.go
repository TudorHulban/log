package helpers

import (
	"fmt"
	"strconv"
)

// AppendJSON writes s into buf with JSON string escaping but without
// the surrounding quotes. Used by appendArgsQuotedJSON.
func AppendJSON(buf []byte, word []byte) []byte {
	const hex = "0123456789abcdef"

	for ix := range word {
		char := word[ix]

		switch char {
		case '\\', '"':
			buf = append(buf, '\\', char)
		case '\n':
			buf = append(buf, '\\', 'n')
		case '\r':
			buf = append(buf, '\\', 'r')
		case '\t':
			buf = append(buf, '\\', 't')
		case '\b':
			buf = append(buf, '\\', 'b')
		case '\f':
			buf = append(buf, '\\', 'f')

		default:
			if char < 0x20 {
				buf = append(buf, '\\', 'u', '0', '0', hex[char>>4], hex[char&0xF])
			} else {
				buf = append(buf, char)
			}
		}
	}

	return buf
}

func AppendJSON_Quoted(buf []byte, text string) []byte {
	buf = append(buf, '"')

	const hex = "0123456789abcdef"

	for ix := 0; ix < len(text); ix++ {
		char := text[ix]

		switch char {
		case '\\', '"':
			buf = append(buf, '\\', char)
		case '\n':
			buf = append(buf, '\\', 'n')
		case '\r':
			buf = append(buf, '\\', 'r')
		case '\t':
			buf = append(buf, '\\', 't')
		case '\b':
			buf = append(buf, '\\', 'b')
		case '\f':
			buf = append(buf, '\\', 'f')

		default:
			if char < 0x20 {
				buf = append(buf, '\\', 'u', '0', '0', hex[char>>4], hex[char&0xF])
			} else {
				buf = append(buf, char)
			}
		}
	}

	return append(buf, '"')
}

// AppendJSON_Arguments writes `"arg0 arg1 …"` directly into buf without any
// intermediate string allocation. It replaces the ArgsToString(args) +
// appendQuotedJSON pattern in the JSON hot path.
func AppendJSON_Arguments(buf []byte, args []any) []byte {
	buf = append(buf, '"')

	for ix, arg := range args {
		if ix > 0 {
			buf = append(buf, ' ')
		}

		switch argument := arg.(type) {
		case string:
			buf = AppendJSON(buf, []byte(argument))
		case []byte:
			buf = AppendJSON(buf, argument)
		case int:
			buf = AppendInt(buf, argument)
		case int64:
			buf = strconv.AppendInt(buf, argument, 10)
		case int32:
			buf = strconv.AppendInt(buf, int64(argument), 10)
		case uint:
			buf = strconv.AppendUint(buf, uint64(argument), 10)
		case uint64:
			buf = AppendUint64(buf, argument)
		case float64:
			buf = strconv.AppendFloat(buf, argument, 'f', -1, 64)
		case float32:
			buf = strconv.AppendFloat(buf, float64(argument), 'f', -1, 32)
		case bool:
			buf = AppendBool(buf, argument)
		case error:
			buf = AppendJSON(buf, []byte(argument.Error()))
		case nil:
			buf = append(buf, 'n', 'u', 'l', 'l')

		default:
			buf = AppendJSON(buf, []byte(fmt.Sprint(argument)))
		}
	}

	buf = append(buf, '"')

	return buf
}

func AppendJSON_Formatted(dst []byte, format string, args ...any) []byte {
	dst = append(dst, '"')

	ai := 0
	flen := len(format)

	for i := 0; i < flen; i++ {
		c := format[i]

		if c != '%' {
			dst = append(dst, c)
			continue
		}

		if i+1 < flen && format[i+1] == '%' {
			dst = append(dst, '%')
			i++

			continue
		}

		i++
		if i >= flen || ai >= len(args) {
			dst = append(dst, '%')
			if i < flen {
				dst = append(dst, format[i])
			}

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
				dst = append(dst, fmt.Sprint(v)...)
			}

		case 'd':
			switch v := arg.(type) {
			case int:
				dst = AppendInt(dst, v)
			case int64:
				dst = strconv.AppendInt(dst, v, 10)
			default:
				dst = append(dst, fmt.Sprint(v)...)
			}

		default:
			dst = append(dst, '%', format[i])
		}
	}

	return append(dst, '"')
}
