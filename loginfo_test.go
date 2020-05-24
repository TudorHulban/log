package loginfo

import (
	"bytes"
	"os"
	"testing"

	"log"

	"github.com/stretchr/testify/assert"
)

func createLogger(level int, t *testing.T) LogInfo {
	l, err := New(level, os.Stderr)
	if assert.Nil(t, err) {
		return l
	}
	return LogInfo{}
}

func Test1Logger(t *testing.T) {
	logger := createLogger(3, t)
	logger.Info("1")
	logger.Warn("2")
	logger.Debug("3")
}

// Benchmark_Logger-4   	  248170	      4614 ns/op	     339 B/op	       2 allocs/op - Flags: log.LstdFlags|log.Lmicroseconds|log.Lshortfile
// Benchmark_Logger-4   	 1208836	       952 ns/op	      79 B/op	       0 allocs/op - Flags: log.LstdFlags
func BenchmarkLogger_Info(b *testing.B) {
	logger, _ := New(1, &bytes.Buffer{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("1")
	}
}

func BenchmarkLogger_Debug(b *testing.B) {
	logger, _ := New(2, &bytes.Buffer{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("1")
	}
}

// Benchmark_SLogger-4   	 1502193	       838 ns/op	      60 B/op	       0 allocs/op - Flags: log.LstdFlags
func Benchmark_SLogger(b *testing.B) {
	logger := log.New(&bytes.Buffer{}, "", log.LstdFlags)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Print("1")
	}
}
