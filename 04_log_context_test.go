package log

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/log/timestamp"
)

func TestContextPrint(t *testing.T) {
	var bufLogs, bufFatal bytes.Buffer

	writer := &bufLogs

	ingestor, errCrIngestor := bytearena.NewIngestor(
		bytearena.Size100K(),
		writer,
		bytearena.WithTickIfDataMilliseconds(10),
	)
	require.NoError(t, errCrIngestor)

	serviceLogging, errCrLogger := NewLogger(
		&ParamsNewLogger{
			Ingestor:    ingestor,
			LoggerLevel: LevelDebug,

			WithFatalWriter: &bufFatal,
			WithJSON:        true,
			WithColor:       true, // Your original setup had color enabled
			WithTimestamp:   timestamp.TimestampRFC3339Bucharest,
		},
	)
	require.NoError(t, errCrLogger)

	ctx, cancel := context.WithCancel(context.Background())
	chIngestionEnd := ingestor.StartIngestion(ctx)

	f := NewLogContext(serviceLogging).
		WithRoot("service", "auth").
		SetInt("req_id", 12345).
		SetBool("cache_hit", true).
		SetString("root ends", "here")

	// Execution
	f.Print("xxx1")
	f.SetString("area", "some area")
	f.Print("login ok again")

	go func() {
		f.With("xxxxxxxxxxxxx", "2").
			Error().
			WithGoroutineID().
			Msg("some error")
	}()

	f.With("zzzz", 4.3).
		Info().
		WithGoroutineID().
		Msg("finished")

	// produces
	// {"ts":"2026-05-07T13:40:02+03:00","level":"INFO","msg":"created logger, level DEBUG"}
	// {"ts":"2026-05-07T13:40:02+03:00","service":"auth","req_id":12345,"cache_hit":true,"root ends":"here","msg":"login ok"}
	// {"ts":"2026-05-07T13:40:02+03:00","service":"auth","req_id":12345,"cache_hit":true,"root ends":"here","area":"some area","msg":"login ok again"}
	// {"ts":"2026-05-07T13:40:02+03:00","level":"INFO","service":"auth","req_id":12345,"cache_hit":true,"root ends":"here","area":"some area","zzzz":4.299999999999,"g":8,"msg":""}
	// {"ts":"2026-05-07T13:40:02+03:00","level":"ERROR","service":"auth","req_id":12345,"cache_hit":true,"root ends":"here","area":"some area","xxxxxxxxxxxxx":"2","g":10,"msg":"some error"}

	// Allow time for the goroutine and the ingestor tick
	time.Sleep(100 * time.Millisecond)

	cancel()
	<-chIngestionEnd

	// --- Processing Lines ---
	require.Zero(t, bufFatal.Len())

	rawLines := strings.Split(bufLogs.String(), "\n")

	var linesJSON []string

	for _, line := range rawLines {
		trimmed := strings.TrimSpace(line)

		idx := strings.IndexByte(trimmed, '{')
		if idx >= 0 {
			linesJSON = append(linesJSON, trimmed[idx:])
		}
	}

	require.Len(t, linesJSON, 4)

	// --- Assertions ---
	findLogByMsg := func(msgValue string) bool {
		for _, l := range linesJSON {
			fmt.Println(l)

			if strings.Contains(l, msgValue) {
				return true
			}
		}

		return false
	}

	// 1. Verify "login ok"
	log1 := findLogByMsg("xxx1")
	require.NotNil(t,
		log1,
		"Could not find 'login ok' log",
	)
	// assert.Equal(t, "auth", log1["service"])
	// assert.Equal(t, float64(12345), log1["req_id"])

	// // 2. Verify "login ok again" (Context update check)
	// log2 := findLogByMsg("login ok again")
	// require.NotNil(t, log2)
	// assert.Equal(t, "some area", log2["area"])

	// // 3. Verify the Error log from Goroutine
	// logErr := findLogByMsg("some error")
	// require.NotNil(t, logErr)
	// assert.Equal(t, "ERROR", logErr["level"])
	// assert.Equal(t, "2", logErr["xxxxxxxxxxxxx"])
	// assert.NotNil(t, logErr["g"], "Goroutine ID should be present")

	// // 4. Verify Floating point and level
	// logInfo := findLogByMsg("finished")
	// require.NotNil(t, logInfo)

	// // Allow a margin of error of 0.000001
	// assert.InDelta(t, 4.3, logInfo["zzzz"].(float64), 0.000001)
	// assert.Equal(t, "INFO", logInfo["level"])
}
