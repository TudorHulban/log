package arena_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/tudorhulban/log/arena"
)

func TestHowToUse(t *testing.T) {
	rawLogger := arena.NewRawLogger(arena.Size100K, &bytes.Buffer{})

	rawLogger.StartIngestion(context.Background())

	// rawLogger.Write("xxx")
}
