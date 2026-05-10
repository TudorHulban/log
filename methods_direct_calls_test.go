package log

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/log/timestamp"
)

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

func TestDirectCalls(t *testing.T) {
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
		keyValue = 1 // Based on your 'produces' comment, the value recorded was 1
	)

	l.Trace("111")
	l.Tracef("%d", value)
	l.Tracew("msg-trace", "key1", keyValue)

	l.Debug("222")
	l.Debugf("%d", value)
	l.Debugw("msg-debug", "key1", keyValue)

	l.Info("333")
	l.Infof("%d", value)
	l.Infow("msg-info", "key1", keyValue)

	l.Warn("444")
	l.Warnf("%d", value)
	l.Warnw("msg-info", "key1", keyValue)

	l.Error("555")
	l.Errorf("%d", value)
	l.Errorw("msg-info", "key1", keyValue)

	// Note: PrintMessage usually triggers an INFO level log
	l.PrintMessage("666")
	// Note: your newTestEntries skips lines without '{',
	// ensure PrintWithNoTimestamp still outputs JSON.
	l.PrintWithNoTimestamp("777")

	cancel()
	<-chIngestionEnd

	fmt.Println(writer.String())

	entries, errParse := newTestEntries(writer.String())
	require.NoError(t, errParse)
	require.NotEmpty(t, entries)

	// --- Assertions using the DSL ---

	// 1. Verify all (or most) entries have a timestamp (except the one from PrintWithNoTimestamp)
	// If PrintWithNoTimestamp removes the "ts" key, haveTimestamp() would return false.
	// Given your provided 'produces' lines all have "ts", we check general existence:
	require.True(t, entries.haveTimestamp(), "All log entries should have a valid RFC3339 timestamp")

	// 2. Count occurrences of log levels
	require.NoError(t, entries.hasKeyWithValue("level", "TRACE", 3))
	require.NoError(t, entries.hasKeyWithValue("level", "DEBUG", 3))
	require.NoError(t, entries.hasKeyWithValue("level", "INFO", 4)) // 3 calls + "created logger" msg
	require.NoError(t, entries.hasKeyWithValue("level", "WARN", 3))
	require.NoError(t, entries.hasKeyWithValue("level", "ERROR", 3))

	// 3. Verify specific messages
	require.NoError(t, entries.hasKeyWithValue("msg", "0", 5)) // Trace, Debug, Info, Warn, Error
	require.NoError(t, entries.hasKeyWithValue("msg", "1", 5)) // The formatted '%d' calls

	// 4. Verify structured logging key-values
	// Checking how many times "key1" appeared with value 1
	require.NoError(t, entries.hasKeyWithValue("key1", keyValue, 5))

	// 5. Verify complex matches (Level + Message + Key)
	// Verify the specific TRACE structured log
	require.NoError(t, entries.hasKeysWithValues(1,
		"level", "TRACE",
		"msg", "msg-trace",
		"key1", keyValue,
	))

	// Verify the specific ERROR structured log
	require.NoError(t, entries.hasKeysWithValues(1,
		"level", "ERROR",
		"msg", "msg-info",
		"key1", keyValue,
	))

	// 6. Verify caller information exists on most lines
	require.NoError(t, entries.hasKey("caller", 15))
}
