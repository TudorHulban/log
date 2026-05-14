package log

import (
	"runtime"
	"strings"
	"testing"
)

func TestGetCallerData(t *testing.T) {
	l := &Logger{
		callerLevel: 2,
	}

	// 1. Verify File Suffix
	// We check that the file path returned contains the current test file name.
	file1, _ := l.getCallerData(2)
	if !strings.HasSuffix(file1, "logger_test.go") {
		t.Errorf("expected file to end with logger_test.go, got %q", file1)
	}

	// 2. Test Accuracy
	// We get the line number of the line immediately before our call.
	_, _, expectedLine, _ := runtime.Caller(0)
	_, line2 := l.getCallerData(2) // This is line 'expectedLine + 1'

	if line2 != expectedLine+1 {
		t.Errorf("expected line %d, got %d", expectedLine+1, line2)
	}

	// 3. Test the Cache (Fast Path)
	// Running this in a loop ensures the atomic pointer swap and
	// subsequent loads from the table are stable.
	for i := 0; i < 100; i++ {
		f, lgn := callHelper(l, 2)
		if !strings.HasSuffix(f, "logger_test.go") || lgn == 0 {
			t.Fatalf("Cache iteration %d failed: got %s:%d", i, f, lgn)
		}
	}
}

// callHelper provides a stable, repeatable call site for the cache test
func callHelper(l *Logger, level int) (string, int) {
	return l.getCallerData(level)
}

func TestCallerLevel(t *testing.T) {
	l := &Logger{}

	// Helper function to wrap the call
	wrapper := func(depth int) (string, int) {
		return l.getCallerData(depth)
	}

	// Level 2 should be the anonymous func (this line)
	f2, _ := wrapper(2)
	// Level 3 should be TestCallerLevel (the caller of the anonymous func)
	f3, _ := wrapper(3)

	if !strings.Contains(f2, "logger_test.go") {
		t.Errorf("Level 2 failed, got %s", f2)
	}
	if !strings.Contains(f3, "logger_test.go") {
		t.Errorf("Level 3 failed, got %s", f3)
	}
}
