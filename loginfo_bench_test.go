package log

import (
	"bytes"
	"testing"

	"log"
)

// BenchmarkLogger_Print-4   	  595407	      1978 ns/op	     179 B/op	       3 allocs/op
func BenchmarkLogger_Print(b *testing.B) {
	logger := New(1, &bytes.Buffer{}, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Print("1")
	}
}

// BenchmarkLogger_Info-4   	  232419	      5516 ns/op	     578 B/op	       6 allocs/op
func BenchmarkLogger_InfoTrue(b *testing.B) {
	logger := New(1, &bytes.Buffer{}, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("1")
	}
}

// BenchmarkLogger_InfoFalse-4   	  590342	      2024 ns/op	     185 B/op	       3 allocs/op
func BenchmarkLogger_InfoFalse(b *testing.B) {
	logger := New(1, &bytes.Buffer{}, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("1")
	}
}

// BenchmarkLogger_Warn-4   	  167058	      7035 ns/op	     811 B/op	      11 allocs/op
func BenchmarkLogger_Warn(b *testing.B) {
	logger := New(3, &bytes.Buffer{}, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Warn("1")
	}
}

// BenchmarkLogger_Debug-4   	  142561	      7325 ns/op	     854 B/op	      11 allocs/op
func BenchmarkLogger_Debug(b *testing.B) {
	logger := New(3, &bytes.Buffer{}, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("1")
	}
}

// Benchmark_StandardLogger-4   	 1502193	       838 ns/op	      60 B/op	       0 allocs/op - Flags: log.LstdFlags
func Benchmark_StandardLogger(b *testing.B) {
	logger := log.New(&bytes.Buffer{}, "", log.LstdFlags)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Print("1")
	}
}
