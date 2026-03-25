package log

import (
	"sync"

	"github.com/tudorhulban/log/helpers"
)

var entryPool = sync.Pool{
	New: func() any {
		return &Entry{
			fields: make([]field, 0, 8), // small reusable buffer
		}
	},
}

// Entry is not safe for concurrent use.
// Each goroutine should obtain its own Entry via Formatter.With.
type Entry struct {
	formatter *LogContext
	fields    []field // per-request, owned by this Entry
}

// With allocates.
//
// Passing as any, the Go compiler must box the concrete value into an interface.
// The boxing causes a heap allocation when the value does not fit in a pointer word
// or the compiler cannot prove it escapes to the stack only.
func (e *Entry) With(key string, value any) *Entry {
	e.fields = append(
		e.fields,
		makeField(key, value),
	)

	return e
}

// WithString appends a string field without boxing value into any.
func (e *Entry) WithString(key, value string) *Entry {
	e.fields = append(
		e.fields,
		field{
			key:         key,
			kind:        kindString,
			valueString: value,
		},
	)

	return e
}

// WithInt appends an int field without boxing value into any.
func (e *Entry) WithInt(key string, value int) *Entry {
	e.fields = append(
		e.fields,
		field{
			key:          key,
			kind:         kindInt,
			valueNumeric: value,
		},
	)

	return e
}

// WithBool appends a bool field without boxing value into any.
func (e *Entry) WithBool(key string, value bool) *Entry {
	e.fields = append(
		e.fields,
		field{
			key:       key,
			kind:      kindBool,
			valueBool: value,
		},
	)

	return e
}

func (e *Entry) Print(args ...any) {
	cfg := e.formatter.cfg.Load()
	logger := e.formatter.logger

	region, errWrite := logger.ingestor.TryWrite(logger.estimatedMessageSize)
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
			buf = appendQuotedJSON(
				buf,
				string(logger.fnTimestamp(nil)),
			)
			buf = append(buf, ',')
		}

		// root field
		if cfg.root != nil {
			fld := cfg.root

			buf = append(buf, '"')
			buf = append(buf, fld.key...)
			buf = append(buf, '"', ':')

			switch fld.kind {
			case kindString:
				buf = appendQuotedJSON(buf, fld.valueString)
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
				buf = appendQuotedJSON(buf, fld.valueString)
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
				buf = appendQuotedJSON(buf, fld.valueString)
			case kindInt:
				buf = helpers.AppendInt(buf, fld.valueNumeric)
			case kindBool:
				buf = helpers.AppendBool(buf, fld.valueBool)
			}

			buf = append(buf, ',')
		}

		// message
		buf = append(buf, `"msg":`...)
		buf = appendArgsQuotedJSON(buf, args)
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

	buf = helpers.AppendArgs(buf, args)
	buf = append(buf, '\n')

	copy(region.Buf(), buf)
	logger.ingestor.EndWrite(region)

	entryPool.Put(e)
}
