package log

import (
	"strconv"

	"github.com/tudorhulban/log/helpers"
)

// appendJSON builds a JSON log entry.
// `msg` must already be a fully formatted string.
// Caller info is included only when file != "" and line > 0.
func (l *Logger) appendJSON(buffer []byte, level, file string, line int, msg []byte) []byte {
	buffer = append(buffer, '{')

	// timestamp
	if l.fnTimestamp != nil {
		buffer = append(buffer, `"ts":"`...)
		buffer = l.fnTimestamp(buffer)
		buffer = append(buffer, `",`...)
	}

	// level
	buffer = append(buffer, `"level":"`...)
	buffer = append(buffer, level...)
	buffer = append(buffer, `",`...)

	// caller info
	if len(file) > 0 && line > 0 {
		buffer = append(buffer, `"caller":"`...)
		buffer = append(buffer, file...)
		buffer = append(buffer, `","line":`...)
		buffer = strconv.AppendInt(buffer, int64(line), 10)
		buffer = append(buffer, ',')
	}

	// message
	buffer = append(buffer, `"msg":"`...)
	buffer = helpers.AppendJSON(buffer, msg)
	buffer = append(buffer, `"}`...)

	return buffer
}

func (l *Logger) appendJSONRoot(buffer []byte, cfg *formatterConfig, file string, line int, msg []byte) []byte {
	buffer = append(buffer, '{')

	// timestamp
	if l.fnTimestamp != nil {
		buffer = append(buffer, `"ts":"`...)
		buffer = l.fnTimestamp(buffer)
		buffer = append(buffer, `",`...)
	}

	// root
	if cfg.root != nil {
		fld := cfg.root

		buffer = append(buffer, '"')
		buffer = append(buffer, fld.key...)
		buffer = append(buffer, '"', ':')

		switch fld.kind {
		case kindString:
			buffer = helpers.AppendJSON_Quoted(
				buffer,
				[]byte(fld.valueString),
			)
		case kindInt:
			buffer = strconv.AppendInt(buffer, fld.valueInt, 10)
		case kindBool:
			buffer = strconv.AppendBool(buffer, fld.valueBool)
		case kindFloat:
			buffer = helpers.AppendFloat(buffer, fld.valueFloat, _PrecisionFloat)
		}

		buffer = append(buffer, ',')
	}

	// context fields
	for ix := range cfg.fields {
		fld := &cfg.fields[ix]

		buffer = append(buffer, '"')
		buffer = append(buffer, fld.key...)
		buffer = append(buffer, '"', ':')

		switch fld.kind {
		case kindString:
			buffer = helpers.AppendJSON_Quoted(
				buffer,
				[]byte(fld.valueString),
			)
		case kindInt:
			buffer = strconv.AppendInt(buffer, fld.valueInt, 10)
		case kindBool:
			buffer = strconv.AppendBool(buffer, fld.valueBool)
		case kindFloat:
			buffer = helpers.AppendFloat(buffer, fld.valueFloat, _PrecisionFloat)
		}

		buffer = append(buffer, ',')
	}

	// caller info
	if len(file) > 0 && line > 0 {
		buffer = append(buffer, `"caller":"`...)
		buffer = append(buffer, file...)
		buffer = append(buffer, `","line":`...)
		buffer = strconv.AppendInt(buffer, int64(line), 10)
		buffer = append(buffer, ',')
	}

	// message
	buffer = append(buffer, `"msg":"`...)
	buffer = helpers.AppendJSON(buffer, msg)
	buffer = append(buffer, `"}`...)

	return buffer
}

// Formats
// {"ts":"2026-03-18T14:27:09.123Z","level":"info","msg":"user login"}
// {"ts":"2026-03-18T14:27:09.123Z","level":"info","msg":"user login","user_id":3847291,"ip":"10.44.12.189","device":"mobile","session_id":"sess_abc123xyz"}
// {"ts":"…","level":"warn","msg":"slow database query","duration_ms":342,"query":"SELECT …","rows":124}
// {"ts":"…","level":"error","msg":"payment failed","error":"card_declined","code":"AUTH_402","attempt":3}
