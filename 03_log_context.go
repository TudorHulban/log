package log

import (
	"sync/atomic"
)

type formatterConfig struct {
	root   *field  // nil if no root
	fields []field // ephemeral fields
}

// LogContext acts as the root.
type LogContext struct {
	logger *Logger
	cfg    atomic.Pointer[formatterConfig]
}

func NewLogContext(logger *Logger) *LogContext {
	f := LogContext{
		logger: logger,
	}

	f.cfg.Store(&formatterConfig{fields: nil})

	return &f
}

func (ctx *LogContext) WithRoot(key string, value any) *LogContext {
	old := ctx.cfg.Load()

	// copy ephemeral fields
	newFields := make([]field, len(old.fields))
	copy(newFields, old.fields)

	// replace root
	newCfg := &formatterConfig{
		root:   makeFieldPtr(key, value),
		fields: newFields,
	}

	ctx.cfg.Store(newCfg)

	return ctx
}

func (ctx *LogContext) SetString(key, value string) *LogContext {
	old := ctx.cfg.Load()

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

	ctx.cfg.Store(newCfg)

	return ctx
}

func (ctx *LogContext) SetInt(key string, value int64) *LogContext {
	old := ctx.cfg.Load()

	newFields := make([]field, len(old.fields)+1)

	copy(newFields, old.fields)
	newFields[len(old.fields)] = field{
		key:      key,
		kind:     kindInt,
		valueInt: value,
	}

	ctx.cfg.Store(&formatterConfig{root: old.root, fields: newFields})

	return ctx
}

func (ctx *LogContext) SetBool(key string, value bool) *LogContext {
	old := ctx.cfg.Load()

	newFields := make([]field, len(old.fields)+1)

	copy(newFields, old.fields)
	newFields[len(old.fields)] = field{
		key:       key,
		kind:      kindBool,
		valueBool: value,
	}

	ctx.cfg.Store(&formatterConfig{root: old.root, fields: newFields})

	return ctx
}

func (ctx *LogContext) Clear() {
	old := ctx.cfg.Load()

	newCfg := &formatterConfig{
		root:   old.root, // keep root
		fields: nil,      // clear ephemeral fields
	}

	ctx.cfg.Store(newCfg)
}

func (ctx *LogContext) Reset() {
	ctx.cfg.Store(&formatterConfig{})
}

func (ctx *LogContext) With(key string, value any) *Entry {
	e, _ := entryPool.Get().(*Entry) //nolint:revive

	e.formatter = ctx
	e.fields = e.fields[:0] // reset slice
	e.fields = append(e.fields, makeField(key, value))

	return e
}

func (ctx *LogContext) WithString(key, value string) *Entry {
	e, _ := entryPool.Get().(*Entry) //nolint:revive

	e.formatter = ctx
	e.fields = e.fields[:0] // reset slice
	e.fields = append(
		e.fields,
		field{
			key:         key,
			valueString: value,
		},
	)

	return e
}
