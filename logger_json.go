package log

import (
	"bytes"
	"fmt"
)

type paramsJSONWCaller struct {
	timestamp string
	file      string
	level     string
	message   string

	line int
}

func json(params *paramsJSONWCaller) []byte {
	var writer bytes.Buffer

	if _, errWrite := fmt.Fprintf(
		&writer,
		`{"timestamp":"%s","%s":"%s"}`,

		params.timestamp,
		params.level,
		params.message,
	); errWrite != nil {
		return nil
	}

	writer.WriteString("\n")

	return writer.Bytes()
}

func jsonWCaller(params *paramsJSONWCaller) []byte {
	var writer bytes.Buffer

	if _, errWrite := fmt.Fprintf(
		&writer,
		`{"timestamp":"%s","file":"%s","line":%d,"%s":"%s"}`,

		params.timestamp,
		params.file,
		params.line,
		params.level,
		params.message,
	); errWrite != nil {
		return nil
	}

	writer.WriteString("\n")

	return writer.Bytes()
}

func (*Logger) appendJSON(buf, ts []byte, level, format string, args ...any) []byte {
	buf = append(buf, '{')

	if ts != nil {
		buf = append(buf, `"ts":"`...)
		buf = append(buf, ts...)
		buf = append(buf, `",`...)
	}

	buf = append(buf, `"level":"`...)
	buf = append(buf, level...)
	buf = append(buf, `","msg":"`...)
	buf = fmt.Appendf(buf, format, args...)
	buf = append(buf, "\"}\n"...)

	return buf
}

// Formats
// {"ts":"2026-03-18T14:27:09.123Z","level":"info","msg":"user login"}
// {"ts":"2026-03-18T14:27:09.123Z","level":"info","msg":"user login","user_id":3847291,"ip":"10.44.12.189","device":"mobile","session_id":"sess_abc123xyz"}
// {"ts":"…","level":"warn","msg":"slow database query","duration_ms":342,"query":"SELECT …","rows":124}
// {"ts":"…","level":"error","msg":"payment failed","error":"card_declined","code":"AUTH_402","attempt":3}
