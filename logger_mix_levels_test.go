package log

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tudorhulban/bytearena"
	"github.com/tudorhulban/log/timestamp"
)

func TestLevelsMatrix(t *testing.T) {
	type tc struct {
		description string
		level       Level
		shouldSee   map[string]string // msg → level
		shouldSkip  []string          // msgs that must not appear
	}

	tests := []tc{
		{
			description: "1. DEBUG threshold → DEBUG, INFO, WARN, ERROR, PRINT emitted",
			level:       LevelDEBUG,
			shouldSee: map[string]string{
				"dbg": `"level":"DEBUG"`,
				"inf": `"level":"INFO"`,
				"wrn": `"level":"WARN"`,
				"err": `"level":"ERROR"`,
				"prt": `"level":"PRINT"`,
			},
			shouldSkip: nil,
		},
		{
			description: "2. INFO threshold → INFO, WARN, ERROR, PRINT emitted; DEBUG suppressed",
			level:       LevelINFO,
			shouldSee: map[string]string{
				"inf": `"level":"INFO"`,
				"wrn": `"level":"WARN"`,
				"err": `"level":"ERROR"`,
				"prt": `"level":"PRINT"`,
			},
			shouldSkip: []string{"dbg"},
		},
		{
			description: "3. WARN threshold → WARN, ERROR, PRINT emitted; DEBUG, INFO suppressed",
			level:       LevelWARN,
			shouldSee: map[string]string{
				"wrn": `"level":"WARN"`,
				"err": `"level":"ERROR"`,
				"prt": `"level":"PRINT"`,
			},
			shouldSkip: []string{"dbg", "inf"},
		},
		{
			description: "4. ERROR threshold → ERROR, PRINT emitted; DEBUG, INFO, WARN suppressed",
			level:       LevelERROR,
			shouldSee: map[string]string{
				"err": `"level":"ERROR"`,
				"prt": `"level":"PRINT"`,
			},
			shouldSkip: []string{"dbg", "inf", "wrn"},
		},
		{
			description: "5. NONE threshold → only PRINT emitted",
			level:       LevelNONE,
			shouldSee: map[string]string{
				"prt": `"level":"PRINT"`,
			},
			shouldSkip: []string{"dbg", "inf", "wrn", "err"},
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.description,
			func(t *testing.T) {
				var writer bytes.Buffer

				ingestor, errCrIngestor := bytearena.NewIngestor(
					bytearena.Size100K(),
					&writer,
				)
				require.NoError(t, errCrIngestor)
				require.NotNil(t, ingestor)

				l, errCrLogger := NewLogger(
					&ParamsNewLogger{
						Ingestor:      ingestor,
						LoggerLevel:   tt.level,
						WithTimestamp: timestamp.TimestampRFC3339Bucharest,
						WithCaller:    true,
						WithColor:     false,
						WithJSON:      true,
					},
				)
				require.NoError(t, errCrLogger)

				ctx, cancel := context.WithCancel(context.Background())
				chIngestionEnd := ingestor.StartIngestion(ctx)

				// Emit all levels
				l.Debug("dbg")
				l.Info("inf")
				l.Warn("wrn")
				l.Error("err")
				l.Print("prt")

				cancel()
				<-chIngestionEnd

				out := writer.String()
				require.NotEmpty(t, out)

				lines := strings.Split(strings.TrimSpace(out), "\n")
				require.NotEmpty(t, lines)

				// Verify expected entries appear
				for msg, lvl := range tt.shouldSee {
					found := false

					for _, ln := range lines {
						if strings.Contains(ln, `"msg":"`+msg+`"`) &&
							strings.Contains(ln, lvl) &&
							strings.Contains(ln, `"ts":`) &&
							strings.Contains(ln, `"caller":`) {
							found = true

							break
						}
					}

					require.True(t,
						found,

						"expected to see msg=%s level=%s in\n%s",
						msg,
						lvl,
						out,
					)
				}

				// Verify suppressed entries do not appear
				for _, msg := range tt.shouldSkip {
					for _, ln := range lines {
						require.NotContains(t,
							ln,

							`"msg":"`+msg+`"`, "msg=%s must be suppressed",
							msg,
						)
					}
				}
			},
		)
	}
}
