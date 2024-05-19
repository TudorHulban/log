package log

import (
	"bytes"
	"sync"
)

const delim = ": "

const (
	// NONE no logging
	NONE = 0
	// INFO level
	INFO = 1
	// WARN level
	WARN = 2
	// DEBUG level
	DEBUG = 3
)

var logLevels = []string{"NONE", "INFO", "WARN", "DEBUG"}

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}
