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
	level     Level
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
func (e *Entry) WithInt(key string, value int64) *Entry {
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

func (e *Entry) Trace() *Entry {
	e.level = LevelTrace

	return e
}

func (e *Entry) Debug() *Entry {
	e.level = LevelDebug

	return e
}

func (e *Entry) Info() *Entry {
	e.level = LevelInfo

	return e
}

func (e *Entry) Warn() *Entry {
	e.level = LevelWarn

	return e
}

func (e *Entry) Error() *Entry {
	e.level = LevelError

	return e
}

func (e *Entry) Fatal() *Entry {
	e.level = LevelFatal

	return e
}

func (e *Entry) Panic() *Entry {
	e.level = LevelPanic

	return e
}

func (e *Entry) estimateFieldsSize() uint32 {
	var result uint32

	for _, f := range e.fields {
		// key
		result = result + uint32(len(`{"key":"`)+len(f.key)+len(`","value":""}`))

		switch f.kind {
		case kindString:
			result = result + uint32(len(f.valueString))

		case kindBool:
			if f.valueBool {
				result = result + 4 // true
			} else {
				result = result + 5 // false
			}

		case kindInt:
			result = result + helpers.DigitsInt(f.valueNumeric)
		}
	}

	return result
}
