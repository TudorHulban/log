package log

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/log/timestamp"
)

type logEntry struct {
	Level string          `json:"level"`
	Msg   string          `json:"msg"`
	Key1  json.RawMessage `json:"key1,omitempty"`
}

// produces:
// {"ts":"2026-05-06T17:04:09+03:00","level":"INFO","caller":"/mnt/tmpfs.ramdisk/log/methods_02_info.go","line":18,"msg":"0"}
// {"ts":"2026-05-06T17:04:09+03:00","level":"ERROR","caller":"/mnt/tmpfs.ramdisk/log/methods_04_error.go","line":43,"msg":"msg-info","key1":1}
// {"ts":"2026-05-06T17:04:09+03:00","level":"INFO","msg":"created logger, level TRACE"}
// {"ts":"2026-05-06T17:04:09+03:00","level":"INFO","caller":"/mnt/tmpfs.ramdisk/log/methods_02_info.go","line":30,"msg":"1"}
// {"ts":"2026-05-06T17:04:09+03:00","level":"TRACE","caller":"/mnt/tmpfs.ramdisk/log/methods_00_trace.go","line":18,"msg":"0"}
// {"ts":"2026-05-06T17:04:09+03:00","level":"INFO","caller":"/mnt/tmpfs.ramdisk/log/methods_02_info.go","line":47,"msg":"msg-info","key1":1}
// {"ts":"2026-05-06T17:04:09+03:00","level":"TRACE","caller":"/mnt/tmpfs.ramdisk/log/methods_00_trace.go","line":30,"msg":"1"}
// {"ts":"2026-05-06T17:04:09+03:00","level":"WARN","caller":"/mnt/tmpfs.ramdisk/log/methods_03_warn.go","line":18,"msg":"0"}
// {"ts":"2026-05-06T17:04:09+03:00","level":"TRACE","caller":"/mnt/tmpfs.ramdisk/log/methods_00_trace.go","line":47,"msg":"msg-trace","key1":1}
// {"ts":"2026-05-06T17:04:09+03:00","level":"WARN","caller":"/mnt/tmpfs.ramdisk/log/methods_03_warn.go","line":30,"msg":"1"}
// {"ts":"2026-05-06T17:04:09+03:00","level":"DEBUG","caller":"/mnt/tmpfs.ramdisk/log/methods_01_debug.go","line":18,"msg":"0"}
// {"ts":"2026-05-06T17:04:09+03:00","level":"WARN","caller":"/mnt/tmpfs.ramdisk/log/methods_03_warn.go","line":43,"msg":"msg-info","key1":1}
// {"ts":"2026-05-06T17:04:09+03:00","level":"DEBUG","caller":"/mnt/tmpfs.ramdisk/log/methods_01_debug.go","line":30,"msg":"1"}
// {"ts":"2026-05-06T17:04:09+03:00","level":"ERROR","caller":"/mnt/tmpfs.ramdisk/log/methods_04_error.go","line":18,"msg":"0"}
// {"ts":"2026-05-06T17:04:09+03:00","level":"DEBUG","caller":"/mnt/tmpfs.ramdisk/log/methods_01_debug.go","line":47,"msg":"msg-debug","key1":1}
// {"ts":"2026-05-06T17:04:09+03:00","level":"ERROR","caller":"/mnt/tmpfs.ramdisk/log/methods_04_error.go","line":30,"msg":"1"}

func TestDirectCallsx(t *testing.T) {
	var writer bytes.Buffer

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		&writer,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	l, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelTrace,

			WithFatalWriter: &writer,
			WithTimestamp:   timestamp.TimestampRFC3339Bucharest,
			WithCaller:      true,
			WithColor:       false,
			WithJSON:        true,
		},
	)
	require.NoError(t, errCrLogger)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	const (
		value    = 1
		keyValue = 777
	)

	l.Trace("0")
	l.Tracef("%d", value)
	l.Tracew("msg-trace", "key1", keyValue)

	l.Debug("0")
	l.Debugf("%d", value)
	l.Debugw("msg-debug", "key1", keyValue)

	l.Info("0")
	l.Infof("%d", value)
	l.Infow("msg-info", "key1", keyValue)

	l.Warn("0")
	l.Warnf("%d", value)
	l.Warnw("msg-info", "key1", keyValue)

	l.Error("0")
	l.Errorf("%d", value)
	l.Errorw("msg-info", "key1", keyValue)

	l.PrintMessage("0")
	l.PrintWithNoTimestamp("0")

	cancel()
	<-chIngestionEnd

	out := writer.String()
	require.NotEmpty(t, out)

	rawLines := strings.Split(strings.TrimSpace(out), "\n")

	var linesJSON []string

	for _, line := range rawLines {
		trimmed := strings.TrimSpace(line)

		idx := strings.IndexByte(trimmed, '{')
		if idx >= 0 {
			linesJSON = append(linesJSON, trimmed[idx:])
		}
	}

	require.NotEmpty(t, linesJSON)

	type cases struct {
		plain     bool
		formatted bool
		withKV    bool
	}

	expected := map[string]*cases{
		"TRACE": {},
		"DEBUG": {},
		"INFO":  {},
		"WARN":  {},
		"ERROR": {},
	}

	for _, line := range linesJSON {
		var e logEntry

		errUnmarshal := json.Unmarshal([]byte(line), &e)
		require.NoError(t, errUnmarshal)

		if e.Msg == "created logger, level TRACE" {
			continue
		}

		c, ok := expected[e.Level]
		require.True(t, ok, "unexpected level in line: %s", line)

		switch {
		case e.Msg == "0":
			require.False(t, c.plain, "duplicate plain case for level %s", e.Level)
			c.plain = true

		case e.Msg == "1":
			require.False(t, c.formatted, "duplicate formatted case for level %s", e.Level)
			c.formatted = true

		case strings.HasPrefix(e.Msg, "msg-"):
			require.False(t, c.withKV, "duplicate withKV case for level %s", e.Level)
			require.NotEmpty(t, e.Key1, "withKV case must have key1 for level %s", e.Level)

			var v int

			errKV := json.Unmarshal(e.Key1, &v)
			require.NoError(t, errKV)
			require.Equal(t, keyValue, v, "invalid key1 value for level %s", e.Level)

			c.withKV = true

		default:
			t.Fatalf("unexpected msg for level %s: %q", e.Level, e.Msg)
		}
	}

	for lvl, c := range expected {
		require.True(t, c.plain, "missing plain case for level %s", lvl)
		require.True(t, c.formatted, "missing formatted case for level %s", lvl)
		require.True(t, c.withKV, "missing withKV case for level %s", lvl)
	}
}
