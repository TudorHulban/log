package log

import "github.com/tudorhulban/log/helpers"

type Request struct {
	formatter *Formatter
	fields    []field // per-request, owned by this Entry
}

func (req *Request) With(key string, value any) *Request {
	req.fields = append(
		req.fields,
		makeField(key, value),
	)

	return req
}

func (e *Request) Print(args ...any) {
	cfg := e.formatter.cfg.Load()

	region, err := e.formatter.logger.ingestor.TryWrite(e.formatter.logger.estimatedMessageSize)
	if err != nil {
		return
	}

	buf := region.Buf()[:0]

	if e.formatter.logger.fnTimestamp != nil {
		buf = e.formatter.logger.fnTimestamp(buf)
		buf = append(buf, ' ')
	}

	if cfg.root != nil {
		buf = appendField(buf, cfg.root)
	}

	for i := range cfg.fields {
		buf = appendField(buf, &cfg.fields[i])
	}

	for i := range e.fields {
		buf = appendField(buf, &e.fields[i])
	}

	buf = helpers.AppendArgs(buf, args)
	buf = append(buf, '\n')

	copy(region.Buf(), buf)
	e.formatter.logger.ingestor.EndWrite(region)
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
		if fld.valueBool {
			buf = append(buf, "true"...)
		} else {
			buf = append(buf, "false"...)
		}
	}

	return append(buf, ' ')
}
