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

func (e *Entry) With(key string, value any) *Entry {
	e.fields = append(
		e.fields,
		makeField(key, value),
	)

	return e
}

func (e *Entry) Print(args ...any) {
	cfg := e.formatter.cfg.Load()

	region, err := e.formatter.logger.ingestor.TryWrite(e.formatter.logger.estimatedMessageSize)
	if err != nil {
		entryPool.Put(e)
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

	// return to pool
	entryPool.Put(e)
}
