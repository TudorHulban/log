package log

import (
	"strconv"

	"github.com/tudorhulban/log/helpers"
)

type fieldKind uint8

const (
	kindString fieldKind = iota
	kindInt
	kindBool
)

type field struct {
	key  string
	kind fieldKind
	sval string
	ival int
	bval bool
}

type Formatter struct {
	logger *Logger
	fields []field // immutable once created
}

func NewFormatter(logger *Logger) Formatter {
	return Formatter{
		logger: logger,
		fields: nil,
	}
}

func (f Formatter) WithString(key, value string) Formatter {
	result := Formatter{
		logger: f.logger,
		fields: make([]field, len(f.fields)+1),
	}

	copy(result.fields, f.fields)

	result.fields[len(f.fields)] = field{
		key:  key,
		kind: kindString,
		sval: value,
	}

	return result
}

func (f Formatter) WithInt(key string, value int) Formatter {
	result := Formatter{
		logger: f.logger,
		fields: make([]field, len(f.fields)+1),
	}

	copy(result.fields, f.fields)

	result.fields[len(f.fields)] = field{
		key:  key,
		kind: kindInt,
		ival: value,
	}

	return result
}

func (f Formatter) WithBool(key string, value bool) Formatter {
	result := Formatter{
		logger: f.logger,
		fields: make([]field, len(f.fields)+1),
	}

	copy(result.fields, f.fields)

	result.fields[len(f.fields)] = field{
		key:  key,
		kind: kindBool,
		bval: value,
	}

	return result
}

func (f Formatter) Print(args ...any) {
	region, errWrite := f.logger.ingestor.TryWrite(f.logger.estimatedMessageSize)
	if errWrite != nil {
		return
	}

	buf := region.Buf()[:0]

	if f.logger.fnTimestamp != nil {
		buf = f.logger.fnTimestamp(buf)
		buf = append(buf, ' ')
	}

	// encode fields
	for ix := range f.fields {
		fld := f.fields[ix]

		buf = append(buf, fld.key...)
		buf = append(buf, '=')

		switch fld.kind {
		case kindString:
			buf = append(buf, fld.sval...)

		case kindInt:
			buf = strconv.AppendInt(buf, int64(fld.ival), 10)

		case kindBool:
			if fld.bval {
				buf = append(buf, "true"...)
			} else {
				buf = append(buf, "false"...)
			}
		}

		buf = append(buf, ' ')
	}

	// encode args
	buf = helpers.AppendArgs(buf, args)
	buf = append(buf, '\n')

	copy(region.Buf(), buf)

	f.logger.ingestor.EndWrite(region)
}
