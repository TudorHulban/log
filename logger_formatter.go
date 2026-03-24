package log

import (
	"sync/atomic"

	"github.com/tudorhulban/log/helpers"
)

type fieldKind uint8

const (
	kindString fieldKind = iota
	kindInt
	kindBool
)

type field struct {
	key          string
	valueString  string
	valueNumeric int

	kind fieldKind

	valueBool bool
}

type formatterConfig struct {
	root   *field  // nil if no root
	fields []field // ephemeral fields
}

type Formatter struct {
	logger *Logger
	cfg    atomic.Pointer[formatterConfig]
}

func NewFormatter(logger *Logger) *Formatter {
	f := &Formatter{logger: logger}
	f.cfg.Store(&formatterConfig{fields: nil})

	return f
}

func makeFieldPtr(key string, value any) *field {
	fld := makeField(key, value)

	return &fld
}

func (f *Formatter) WithRoot(key string, value any) *Formatter {
	old := f.cfg.Load()

	// copy ephemeral fields
	newFields := make([]field, len(old.fields))
	copy(newFields, old.fields)

	// replace root
	newCfg := &formatterConfig{
		root:   makeFieldPtr(key, value),
		fields: newFields,
	}

	f.cfg.Store(newCfg)

	return f
}

func (f *Formatter) WithString(key, value string) *Formatter {
	old := f.cfg.Load()

	newFields := make([]field, len(old.fields)+1)
	copy(newFields, old.fields)

	newFields[len(old.fields)] = field{
		key:         key,
		kind:        kindString,
		valueString: value,
	}

	newCfg := &formatterConfig{
		root:   old.root,  // keep root
		fields: newFields, // updated ephemeral fields
	}

	f.cfg.Store(newCfg)

	return f
}

func (f *Formatter) WithInt(key string, value int) *Formatter {
	old := f.cfg.Load()

	newFields := make([]field, len(old.fields)+1)
	copy(newFields, old.fields)
	newFields[len(old.fields)] = field{
		key:          key,
		kind:         kindInt,
		valueNumeric: value,
	}

	f.cfg.Store(&formatterConfig{root: old.root, fields: newFields})

	return f
}

func (f *Formatter) WithBool(key string, value bool) *Formatter {
	old := f.cfg.Load()

	newFields := make([]field, len(old.fields)+1)
	copy(newFields, old.fields)
	newFields[len(old.fields)] = field{
		key:       key,
		kind:      kindBool,
		valueBool: value,
	}

	f.cfg.Store(&formatterConfig{root: old.root, fields: newFields})

	return f
}

func (f *Formatter) ClearFields() {
	old := f.cfg.Load()

	newCfg := &formatterConfig{
		root:   old.root, // keep root
		fields: nil,      // clear ephemeral fields
	}

	f.cfg.Store(newCfg)
}

func (f *Formatter) Reset() {
	f.cfg.Store(&formatterConfig{})
}

func (f *Formatter) Print(args ...any) {
	cfg := f.cfg.Load() // atomic read

	region, err := f.logger.ingestor.TryWrite(f.logger.estimatedMessageSize)
	if err != nil {
		return
	}

	buf := region.Buf()[:0]

	if f.logger.fnTimestamp != nil {
		buf = f.logger.fnTimestamp(buf)
		buf = append(buf, ' ')
	}

	// encode root field first (if present)
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
			if fld.valueBool {
				buf = append(buf, "true"...)
			} else {
				buf = append(buf, "false"...)
			}
		}

		buf = append(buf, ' ')
	}

	// encode ephemeral fields
	for i := range cfg.fields {
		fld := &cfg.fields[i]

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

		buf = append(buf, ' ')
	}

	buf = helpers.AppendArgs(buf, args)
	buf = append(buf, '\n')

	copy(region.Buf(), buf)
	f.logger.ingestor.EndWrite(region)
}
