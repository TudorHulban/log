package log

import (
	"strconv"
)

func appendMsg(buf []byte, format string, args ...any) []byte {
	argIdx := 0
	position := 0

	for ix := 0; ix < len(format); ix++ {
		if format[ix] == '%' && ix+1 < len(format) {
			// flush literal segment safely
			if position < ix {
				buf = appendJSONString(buf, format[position:ix])
			}

			ix++

			switch format[ix] {
			case 's':
				buf = appendJSONString(buf, args[argIdx].(string))

			case 'd':
				buf = strconv.AppendInt(buf, int64(args[argIdx].(int)), 10)

			case 'v':
				buf = appendAny(buf, args[argIdx])

			default:
				buf = append(buf, '%', format[ix])
			}

			argIdx++
			position = ix + 1
		}
	}

	// flush remaining literal
	if position < len(format) {
		buf = appendJSONString(buf, format[position:])
	}

	return buf
}

func appendAny(buf []byte, v any) []byte {
	switch x := v.(type) {
	case string:
		return appendJSONString(buf, x)

	case int:
		return strconv.AppendInt(buf, int64(x), 10)

	case int64:
		return strconv.AppendInt(buf, x, 10)

	case bool:
		return strconv.AppendBool(buf, x)

	case float64:
		return strconv.AppendFloat(buf, x, 'f', -1, 64)

	default:
		return append(buf, "<unsupported>"...)
	}
}

func (l *Logger) appendJSON(buf []byte, level, format string, args ...any) []byte {
	buf = append(buf, '{')

	if l.fnTimestamp != nil {
		buf = append(buf, `"ts":"`...)
		buf = l.fnTimestamp(buf)
		buf = append(buf, `",`...)
	}

	buf = append(buf, `"level":"`...)
	buf = append(buf, level...)
	buf = append(buf, `","msg":"`...)
	buf = appendMsg(buf, format, args...)
	buf = append(buf, "\"}\n"...)

	return buf
}

// Formats
// {"ts":"2026-03-18T14:27:09.123Z","level":"info","msg":"user login"}
// {"ts":"2026-03-18T14:27:09.123Z","level":"info","msg":"user login","user_id":3847291,"ip":"10.44.12.189","device":"mobile","session_id":"sess_abc123xyz"}
// {"ts":"…","level":"warn","msg":"slow database query","duration_ms":342,"query":"SELECT …","rows":124}
// {"ts":"…","level":"error","msg":"payment failed","error":"card_declined","code":"AUTH_402","attempt":3}
