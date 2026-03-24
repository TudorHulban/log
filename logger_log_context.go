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

type LogContext struct {
	logger *Logger
	cfg    atomic.Pointer[formatterConfig]
}

func NewLogContext(logger *Logger) *LogContext {
	f := LogContext{logger: logger}
	f.cfg.Store(&formatterConfig{fields: nil})

	return &f
}

func makeFieldPtr(key string, value any) *field {
	fld := makeField(key, value)

	return &fld
}

func (f *LogContext) WithRoot(key string, value any) *LogContext {
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

func (ctx *LogContext) With(key string, value any) *Entry {
	e := entryPool.Get().(*Entry)

	e.formatter = ctx
	e.fields = e.fields[:0] // reset slice
	e.fields = append(e.fields, makeField(key, value))

	return e
}

func (f *LogContext) SetString(key, value string) *LogContext {
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

func (f *LogContext) SetInt(key string, value int) *LogContext {
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

func (f *LogContext) SetBool(key string, value bool) *LogContext {
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

func (f *LogContext) Clear() {
	old := f.cfg.Load()

	newCfg := &formatterConfig{
		root:   old.root, // keep root
		fields: nil,      // clear ephemeral fields
	}

	f.cfg.Store(newCfg)
}

func (f *LogContext) Reset() {
	f.cfg.Store(&formatterConfig{})
}

func (f *LogContext) Print(args ...any) {
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
		buf = appendField(buf, cfg.root)
	}

	// encode ephemeral fields
	for i := range cfg.fields {
		buf = appendField(buf, &cfg.fields[i])
	}

	buf = helpers.AppendArgs(buf, args)
	buf = append(buf, '\n')

	copy(region.Buf(), buf)
	f.logger.ingestor.EndWrite(region)
}
