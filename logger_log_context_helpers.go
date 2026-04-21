package log

import (
	"fmt"
	"strconv"

	"github.com/tudorhulban/log/helpers"
)

func makeField(key string, value any) field {
	switch v := value.(type) {
	case string:
		return field{
			key:         key,
			kind:        kindString,
			valueString: v,
		}

	case int:
		return field{
			key:          key,
			kind:         kindInt,
			valueNumeric: v,
		}

	case bool:
		return field{
			key:       key,
			kind:      kindBool,
			valueBool: v,
		}

	// You can add more typed cases here:
	// case uint:
	// case float64:
	// case error:
	// etc.

	default:
		// Fallback: convert to string once
		return field{
			key:         key,
			kind:        kindString,
			valueString: fmt.Sprint(v),
		}
	}
}

func appendField(buf []byte, fld *field) []byte {
	buf = append(buf, fld.key...)
	buf = append(buf, '=')

	switch fld.kind {
	case kindString:
		buf = append(buf, fld.valueString...)

	case kindInt:
		buf = helpers.AppendInt(buf, fld.valueNumeric)

	case kindBool:
		buf = helpers.AppendBool(buf, fld.valueBool)
	}

	return append(buf, ' ')
}

// appendJSONEscaped writes s into buf with JSON string escaping but without
// the surrounding quotes. Used by appendArgsQuotedJSON.
func appendJSONEscaped(buf []byte, s string) []byte {
	const hex = "0123456789abcdef"

	for i := 0; i < len(s); i++ {
		c := s[i]

		switch c {
		case '\\', '"':
			buf = append(buf, '\\', c)
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
			if c < 0x20 {
				buf = append(buf, '\\', 'u', '0', '0', hex[c>>4], hex[c&0xF])
			} else {
				buf = append(buf, c)
			}
		}
	}

	return buf
}

func appendQuotedJSON(buf []byte, s string) []byte {
	buf = append(buf, '"')
	const hex = "0123456789abcdef"

	for i := 0; i < len(s); i++ {
		c := s[i]

		switch c {
		case '\\', '"':
			buf = append(buf, '\\', c)
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
			if c < 0x20 {
				buf = append(buf, '\\', 'u', '0', '0', hex[c>>4], hex[c&0xF])
			} else {
				buf = append(buf, c)
			}
		}
	}

	return append(buf, '"')
}

// appendArgsQuotedJSON writes `"arg0 arg1 …"` directly into buf without any
// intermediate string allocation. It replaces the ArgsToString(args) +
// appendQuotedJSON pattern in the JSON hot path.
func appendArgsQuotedJSON(buf []byte, args []any) []byte {
	buf = append(buf, '"')

	for i, arg := range args {
		if i > 0 {
			buf = append(buf, ' ')
		}

		switch v := arg.(type) {
		case string:
			buf = appendJSONEscaped(buf, v)
		case []byte:
			buf = appendJSONEscaped(buf, string(v)) // string([]byte) is stack-alloc'd by the compiler for small slices
		case int:
			buf = helpers.AppendInt(buf, v)
		case int64:
			buf = strconv.AppendInt(buf, v, 10)
		case int32:
			buf = strconv.AppendInt(buf, int64(v), 10)
		case uint:
			buf = strconv.AppendUint(buf, uint64(v), 10)
		case uint64:
			buf = helpers.AppendUint64(buf, v)
		case float64:
			buf = strconv.AppendFloat(buf, v, 'f', -1, 64)
		case float32:
			buf = strconv.AppendFloat(buf, float64(v), 'f', -1, 32)
		case bool:
			buf = helpers.AppendBool(buf, v)
		case error:
			buf = appendJSONEscaped(buf, v.Error())
		case nil:
			buf = append(buf, 'n', 'u', 'l', 'l')
		default:
			buf = appendJSONEscaped(buf, fmt.Sprint(v))
		}
	}

	buf = append(buf, '"')

	return buf
}
