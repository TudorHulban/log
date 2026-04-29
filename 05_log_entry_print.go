package log

import "github.com/tudorhulban/log/helpers"

func (e *Entry) Msg(msg string) {
	if Level(e.formatter.logger.logLevel.Load()) > e.level {
		return
	}

	cfg := e.formatter.cfg.Load()
	logger := e.formatter.logger

	region, errWrite := logger.ingestor.TryWrite(
		uint32(len(msg) + e.estimateFieldsSize() + 128),
	)
	if errWrite != nil {
		entryPool.Put(e)

		return
	}

	buf := region.Buf()[:0]

	// JSON MODE
	if logger.withJSON {
		buf = append(buf, '{')

		// timestamp
		if logger.fnTimestamp != nil {
			buf = append(buf, `"ts":`...)
			buf = helpers.AppendJSON_Quoted(
				buf,
				string(logger.fnTimestamp(nil)),
			)
			buf = append(buf, ',')
		}

		// level
		buf = append(buf, `"level":`...)
		buf = helpers.AppendJSON_Quoted(buf, e.level.String())
		buf = append(buf, ',')

		// root field
		if cfg.root != nil {
			fld := cfg.root

			buf = append(buf, '"')
			buf = append(buf, fld.key...)
			buf = append(buf, '"', ':')

			switch fld.kind {
			case kindString:
				buf = helpers.AppendJSON_Quoted(buf, fld.valueString)
			case kindInt:
				buf = helpers.AppendInt(buf, fld.valueNumeric)
			case kindBool:
				buf = helpers.AppendBool(buf, fld.valueBool)
			}

			buf = append(buf, ',')
		}

		// context fields
		for ix := range cfg.fields {
			fld := &cfg.fields[ix]

			buf = append(buf, '"')
			buf = append(buf, fld.key...)
			buf = append(buf, '"', ':')

			switch fld.kind {
			case kindString:
				buf = helpers.AppendJSON_Quoted(buf, fld.valueString)
			case kindInt:
				buf = helpers.AppendInt(buf, fld.valueNumeric)
			case kindBool:
				buf = helpers.AppendBool(buf, fld.valueBool)
			}

			buf = append(buf, ',')
		}

		// entry fields
		for ix := range e.fields {
			fld := &e.fields[ix]

			buf = append(buf, '"')
			buf = append(buf, fld.key...)
			buf = append(buf, '"', ':')

			switch fld.kind {
			case kindString:
				buf = helpers.AppendJSON_Quoted(buf, fld.valueString)
			case kindInt:
				buf = helpers.AppendInt(buf, fld.valueNumeric)
			case kindBool:
				buf = helpers.AppendBool(buf, fld.valueBool)
			}

			buf = append(buf, ',')
		}

		// message
		buf = append(buf, `"message":`...)
		buf = helpers.AppendJSON(buf, []byte(msg))
		buf = append(buf, '}', '\n')

		copy(region.Buf(), buf)
		logger.ingestor.EndWrite(region)
		entryPool.Put(e)

		return
	}

	// TEXT MODE (fast path)
	if logger.fnTimestamp != nil {
		buf = logger.fnTimestamp(buf)
		buf = append(buf, ' ')
	}

	// root
	if cfg.root != nil {
		fld := cfg.root

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

		buf = append(buf, ' ')
	}

	// context fields
	for ix := range cfg.fields {
		fld := &cfg.fields[ix]

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

		buf = append(buf, ' ')
	}

	// entry fields
	for ix := range e.fields {
		fld := &e.fields[ix]

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

		buf = append(buf, ' ')
	}

	buf = helpers.AppendJSON(buf, []byte(msg))
	buf = append(buf, '\n')

	copy(region.Buf(), buf)
	logger.ingestor.EndWrite(region)

	entryPool.Put(e)
}
