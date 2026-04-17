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
