package log

import (
	"context"
	"io"

	"github.com/gofiber/fiber/v2/log"
)

func (l Logger) SetLevel(level log.Level) {}

func (l Logger) SetOutput(newOutput io.Writer) {}

func (l Logger) WithContext(ctx context.Context) log.CommonLogger {
	return l
}
