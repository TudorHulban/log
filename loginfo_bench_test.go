package log

import (
	"bytes"
	"testing"

	"log"
)

// goos: linux
// goarch: amd64
// pkg: github.com/TudorHulban/log
// cpu: AMD Ryzen 5 5600U with Radeon Graphics
// BenchmarkLogger_Print-12    	 3268773	       352.3 ns/op	      86 B/op	       2 allocs/op
func BenchmarkLogger_Print(b *testing.B) {
	logger := NewLogger(1, &bytes.Buffer{}, true)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Print("1")
	}
}

// BenchmarkLogger_InfoTrue-12    	 1000000	      1099 ns/op	     654 B/op	       5 allocs/op
func BenchmarkLogger_InfoTrue(b *testing.B) {
	logger := NewLogger(1, &bytes.Buffer{}, true)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("1")
	}
}

// BenchmarkLogger_InfoFalse-12    	 3145372	       374.3 ns/op	      92 B/op	       2 allocs/op
func BenchmarkLogger_InfoFalse(b *testing.B) {
	logger := NewLogger(1, &bytes.Buffer{}, false)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Info("1")
	}
}

// BenchmarkLogger_Warn-12    	  993535	      1396 ns/op	     735 B/op	       9 allocs/op
func BenchmarkLogger_Warn(b *testing.B) {
	logger := NewLogger(3, &bytes.Buffer{}, true)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Warn("1")
	}
}

// BenchmarkLogger_Debug-12    	  963235	      1538 ns/op	     748 B/op	       9 allocs/op
func BenchmarkLogger_Debug(b *testing.B) {
	logger := NewLogger(3, &bytes.Buffer{}, true)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Debug("1")
	}
}

// Benchmark_StandardLogger-12    	 6683350	       186.4 ns/op	      53 B/op	       0 allocs/op
func Benchmark_StandardLogger(b *testing.B) {
	logger := log.New(&bytes.Buffer{}, "", log.LstdFlags)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		logger.Print("1")
	}
}
