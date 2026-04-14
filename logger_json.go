package log

import (
	"strconv"

	"github.com/tudorhulban/log/helpers"
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

func appendAny(buf []byte, value any) []byte {
	switch x := value.(type) {
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

func appendEscapedJSON(buf []byte, s string) []byte {
	for i := 0; i < len(s); i++ {
		switch c := s[i]; c {
		case '\\':
			buf = append(buf, '\\', '\\')
		case '"':
			buf = append(buf, '\\', '"')
		case '\n':
			buf = append(buf, '\\', 'n')
		case '\r':
			buf = append(buf, '\\', 'r')
		case '\t':
			buf = append(buf, '\\', 't')
		default:
			buf = append(buf, c)
		}
	}

	return buf
}

// appendJSON builds a JSON log entry.
// `msg` must already be a fully formatted string.
// Caller info is included only when file != "" and line > 0.
func (l *Logger) appendJSON(buf []byte, level, file string, line int, msg string) []byte {
	buf = append(buf, '{')

	// timestamp
	if l.fnTimestamp != nil {
		buf = append(buf, `"ts":"`...)
		buf = l.fnTimestamp(buf)
		buf = append(buf, `",`...)
	}

	// level
	buf = append(buf, `"level":"`...)
	buf = append(buf, level...)
	buf = append(buf, `",`...)

	// caller info
	if len(file) > 0 && line > 0 {
		buf = append(buf, `"caller":"`...)
		buf = append(buf, file...)
		buf = append(buf, `","line":`...)
		buf = helpers.AppendInt(buf, line)
		buf = append(buf, ',')

	}

	// message
	buf = append(buf, `"msg":"`...)
	buf = appendEscapedJSON(buf, msg)
	buf = append(buf, `"}`...)

	return buf
}

// Formats
// {"ts":"2026-03-18T14:27:09.123Z","level":"info","msg":"user login"}
// {"ts":"2026-03-18T14:27:09.123Z","level":"info","msg":"user login","user_id":3847291,"ip":"10.44.12.189","device":"mobile","session_id":"sess_abc123xyz"}
// {"ts":"…","level":"warn","msg":"slow database query","duration_ms":342,"query":"SELECT …","rows":124}
// {"ts":"…","level":"error","msg":"payment failed","error":"card_declined","code":"AUTH_402","attempt":3}
