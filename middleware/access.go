package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/nerocho/fastfiber"
)

type AcccessOptions struct {
	LogResponse bool
}

//访问日志 中间件,官方的日志中间件不太灵活，先自己实现
func Access(ops *AcccessOptions) fiber.Handler {
	accessLogger := fastfiber.Logger.With().Str("type", "access").Logger()

	return func(c *fiber.Ctx) (err error) {
		start := time.Now()
		err = c.Next()
		latency := time.Since(start).Milliseconds()

		ev := accessLogger.Info()
		ev.Int("duration", int(latency))
		ev.Str("method", string(c.Request().Header.Method()))
		ev.Str("hostname", c.Hostname())
		ev.Str("url", c.Path())
		ev.Int("status", c.Response().StatusCode())
		ev.Str("ip", c.IP())
		ev.Str("queryString", c.Request().URI().QueryArgs().String())
		ev.Str("postString", c.Request().PostArgs().String())
		if ops.LogResponse {
			ev.Bytes("response", c.Response().Body())
		}
		ev.Msg("")
		return err
	}
}
