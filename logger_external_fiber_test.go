package log_test

import (
	"os"
	"testing"

	"github.com/TudorHulban/log"
	"github.com/TudorHulban/log/timestamp"
	fiberlog "github.com/gofiber/fiber/v2/log"
)

func TestFiber(t *testing.T) {
	l := log.NewLogger(
		&log.ParamsNewLogger{
			LoggerLevel:   log.LevelDEBUG,
			LoggerWriter:  os.Stdout,
			WithTimestamp: timestamp.TimestampNano,
		},
	)

	fiberlog.SetLogger(l)
	fiberlog.SetLevel(log.LevelDEBUG)

	// app := fiber.New()

	// l.Fatal(app.Listen(":3000"))
}
