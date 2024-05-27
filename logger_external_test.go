package log_test

// File details how to use logger.

import (
	"bytes"
	"os"
	"testing"

	"github.com/TudorHulban/log"
	"github.com/TudorHulban/log/timestamp"
	"github.com/gofiber/fiber/v2"
	fiberlog "github.com/gofiber/fiber/v2/log"
)

type T struct {
	l log.Logger
}

func TestExternal(t *testing.T) {
	obj := T{
		l: log.NewLogger(
			&log.ParamsNewLogger{
				LoggerLevel:  log.LevelDEBUG,
				LoggerWriter: new(bytes.Buffer),

				WithTimestamp: timestamp.TimestampNano,
				WithCaller:    true,
			},
		),
	}

	obj.l.Info("xxx")
	obj.l.Debug("yyy")

	// assert.Contains(t, output.String(), "xxx") - race condition
}

func TestFiber(t *testing.T) {
	l := log.NewLogger(
		&log.ParamsNewLogger{
			LoggerLevel:   log.LevelDEBUG,
			LoggerWriter:  os.Stdout,
			WithTimestamp: timestamp.TimestampNano,
		},
	)

	fiberlog.SetLogger(l) // does not work!

	app := fiber.New()

	app.Use(l)

	l.Fatal(app.Listen(":3000"))
}
