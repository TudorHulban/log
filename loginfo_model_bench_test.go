package log

import (
	"bytes"
	"testing"

	"log"
)

// Values on 8 cores with sync pool.
/*
goos: linux
goarch: amd64
pkg: github.com/TudorHulban/log
cpu: AMD Ryzen 5 3400G with Radeon Vega Graphics
BenchmarkLogger_Print-8   	 2143802	       542.9 ns/op	     115 B/op	     2 allocs/op
*/
func BenchmarkLogger_Print(b *testing.B) {
	logger := NewLogger(1, &bytes.Buffer{}, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Print("1")
	}
}

// BenchmarkLogger_InfoTrue-8   	  782294	      1322 ns/op	     550 B/op	       5 allocs/op
func BenchmarkLogger_InfoTrue(b *testing.B) {
	logger := NewLogger(1, &bytes.Buffer{}, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("1")
	}
}

// BenchmarkLogger_InfoFalse-8   	 2058907	       587 ns/op	     133 B/op	       2 allocs/op
func BenchmarkLogger_InfoFalse(b *testing.B) {
	logger := NewLogger(1, &bytes.Buffer{}, false)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("1")
	}
}

// BenchmarkLogger_Warn-8   	  641996	      1608 ns/op	     667 B/op	       9 allocs/op
func BenchmarkLogger_Warn(b *testing.B) {
	logger := NewLogger(3, &bytes.Buffer{}, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Warn("1")
	}
}

// BenchmarkLogger_Debug-8   	  656444	      1618 ns/op	     663 B/op	       9 allocs/op
func BenchmarkLogger_Debug(b *testing.B) {
	logger := NewLogger(3, &bytes.Buffer{}, true)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("1")
	}
}

// Benchmark_StandardLogger-8   	 4824270	       246 ns/op	      74 B/op	       0 allocs/op
func Benchmark_StandardLogger(b *testing.B) {
	logger := log.New(&bytes.Buffer{}, "", log.LstdFlags)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Print("1")
	}
}
