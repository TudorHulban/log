package log_test

// File details how to use logger.

import (
	"bytes"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/log"
	safewriter "github.com/tudorhulban/log/safe-writer"
	"github.com/tudorhulban/log/timestamp"
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

	for ix := range noWorkers {
		go work(ix)
	}

	wg.Wait()
}
