package log

import (
	"sync/atomic"

	"github.com/tudorhulban/log/helpers"
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

func (ctx *LogContext) SetInt(key string, value int) *LogContext {
	old := ctx.cfg.Load()

	newFields := make([]field, len(old.fields)+1)
	copy(newFields, old.fields)
	newFields[len(old.fields)] = field{
		key:          key,
		kind:         kindInt,
		valueNumeric: value,
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

func (ctx *LogContext) Print(args ...any) {
	cfg := ctx.cfg.Load() // atomic read

	region, err := ctx.logger.ingestor.TryWrite(ctx.logger.estimatedMessageSizeOverall)
	if err != nil {
		return
	}

	buf := region.Buf()[:0]

	if ctx.logger.fnTimestamp != nil {
		buf = ctx.logger.fnTimestamp(buf)
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
	ctx.logger.ingestor.EndWrite(region)
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
