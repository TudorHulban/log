package log_test

// File details how to use logger.

import (
	"bytes"
	"os"
	"sync"
	"testing"

	"github.com/TudorHulban/log"
	safewriter "github.com/TudorHulban/log/safe-writer"
	"github.com/TudorHulban/log/timestamp"
	"github.com/stretchr/testify/require"
)

type T struct {
	l *log.Logger
}

func TestSimpleExternal(t *testing.T) {
	var output bytes.Buffer

	obj := T{
		l: log.NewLogger(
			&log.ParamsNewLogger{
				LoggerLevel:  log.LevelDEBUG,
				LoggerWriter: &output,

				WithTimestamp: timestamp.TimestampNano,
				WithCaller:    true,
				WithColor:     true,
			},
		),
	}

	msg1 := "xxxxx"

	obj.l.Info(msg1)

	require.Contains(t,
		output.String(),
		msg1,
	)
}

func TestMultiExternal(t *testing.T) {
	output := safewriter.NewSafeWriter(os.Stdout)

	obj := T{
		l: log.NewLogger(
			&log.ParamsNewLogger{
				LoggerLevel:  log.LevelDEBUG,
				LoggerWriter: output,

				WithTimestamp: timestamp.TimestampNano,
				WithCaller:    true,
				WithColor:     true,
			},
		),
	}

	var wg sync.WaitGroup

	noWorkers := 5

	wg.Add(noWorkers)

	work := func(text any) {
		obj.l.Info(text)

		wg.Done()
	}

	for i := 0; i < noWorkers; i++ {
		go work(i)
	}

	wg.Wait()
}
