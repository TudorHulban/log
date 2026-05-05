package log

import (
	"bytes"
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"
	"unicode"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/log/timestamp"
)

func TestArenalog_MultipleFields_AllLevels(t *testing.T) {
	// 1. Create a buffer to capture output
	var buf bytes.Buffer

	// 2. Setup Ingestor with the buffer as the writer
	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		&buf,
	)
	require.NoError(t, errCrIngestor)
	require.NotNil(t, ingestor)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	// 3. Setup Logger
	logger, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelTrace,

			WithFatalWriter: &buf,
			WithTimestamp:   timestamp.TimestampRFC3339,
			WithJSON:        true,
		},
	)
	require.NoError(t, errCrLogger)

	// 4. Create Context with Base Fields
	logContext := NewLogContext(logger).
		WithRoot("service", "auth").
		SetInt("req_id", 12345).
		SetBool("cache_hit", true)

	// 5. Generate Logs for All Levels

	// --- TRACE ---
	logContext.
		WithString("area", "trace-area").
		Trace().
		WithString("user", "arena-trace").
		Msg("trace message")

	// --- DEBUG ---
	logContext.
		WithString("area", "debug-area").
		Debug().
		WithString("user", "arena-debug").
		WithInt("attempt", 1).
		Msg("debug message")

	// --- INFO ---
	logContext.
		WithString("area", "info-area").
		Info().
		WithString("user", "arena-info").
		WithFloat("some_float", 1.1137).
		WithBool("success", true).
		Msg("info message")

	// --- ERROR ---
	logContext.
		WithString("area", "error-area").
		Error().
		WithString("user", "arena-error").
		WithString("error_detail", "something failed").
		Msg("error message")

	// 6. Stop Ingestor and Wait for flush
	cancel()
	<-chIngestionEnd

	// 7. Parse and Assert Output
	output := buf.String()

	// Split by newline and filter out empty strings
	linesRaw := strings.Split(output, "\n")
	require.GreaterOrEqual(t, len(linesRaw), 3)

	var linesJSON []string

	for _, line := range linesRaw {
		trimmed := strings.TrimFunc(line, unicode.IsSpace)

		idx := strings.IndexByte(trimmed, '{')
		if idx >= 0 {
			linesJSON = append(linesJSON, trimmed[idx:])
		}
	}

	// We expect 5 lines: 1 Init message + 4 Log messages
	require.Equal(t,
		5,
		len(linesJSON),

		"Expected 5 JSON log lines, got %d. Output:\n%s",
		len(linesJSON),
		output,
	)

	// Helper to parse JSON line
	parseLog := func(line string) map[string]any {
		var m map[string]any

		require.NoError(t,
			json.Unmarshal([]byte(line), &m),

			"Failed to parse JSON: %s",
			line,
		)

		return m
	}

	// --- Line 0: Initialization Message ---
	initLog := parseLog(linesJSON[0])
	require.Equal(t, "INFO", initLog["level"])
	require.Equal(t, "created logger, level TRACE", initLog["msg"])

	// --- Line 1: TRACE ---
	traceLog := parseLog(linesJSON[1])
	require.Equal(t, "TRACE", traceLog["level"])
	require.Equal(t, "trace message", traceLog["msg"])
	require.Equal(t, "auth", traceLog["service"])
	require.Equal(t, float64(12345), traceLog["req_id"])
	require.Equal(t, true, traceLog["cache_hit"])
	require.Equal(t, "trace-area", traceLog["area"])
	require.Equal(t, "arena-trace", traceLog["user"])

	// --- Line 2: DEBUG ---
	debugLog := parseLog(linesJSON[2])
	require.Equal(t, "DEBUG", debugLog["level"])
	require.Equal(t, "debug message", debugLog["msg"])
	require.Equal(t, "debug-area", debugLog["area"])
	require.Equal(t, "arena-debug", debugLog["user"])
	require.Equal(t, float64(1), debugLog["attempt"])

	// --- Line 3: INFO ---
	infoLog := parseLog(linesJSON[3])
	require.Equal(t, "INFO", infoLog["level"])
	require.Equal(t, "info message", infoLog["msg"])
	require.Equal(t, "info-area", infoLog["area"])
	require.Equal(t, "arena-info", infoLog["user"])
	require.Equal(t, true, infoLog["success"])

	// Float comparison with tolerance
	infoFloat := infoLog["some_float"].(float64)
	require.InDelta(t, 1.1137, infoFloat, 0.0001)

	// --- Line 4: ERROR ---
	errorLog := parseLog(linesJSON[4])
	require.Equal(t, "ERROR", errorLog["level"])
	require.Equal(t, "error message", errorLog["msg"])
	require.Equal(t, "error-area", errorLog["area"])
	require.Equal(t, "arena-error", errorLog["user"])
	require.Equal(t, "something failed", errorLog["error_detail"])

	// Verify Timestamp exists and is valid RFC3339 on one of the entries
	ts, couldCast := traceLog["ts"].(string)
	require.True(t,
		couldCast,
		"Timestamp missing",
	)

	_, errParse := time.Parse(time.RFC3339, ts)
	require.NoError(t,
		errParse,

		"Invalid timestamp format: %s",
		ts,
	)
}
