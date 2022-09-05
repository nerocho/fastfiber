package middleware

import (
	"time"

	winner_logger "github.com/bfmTech/logger-go"
	"github.com/gofiber/fiber/v2"
)

type AccessOptions struct {
	Logger      winner_logger.Logger
	LogResponse bool
}

//访问日志 中间件,官方的日志中间件不太灵活，先自己实现
func Access(ops *AccessOptions) fiber.Handler {

	return func(c *fiber.Ctx) (err error) {
		start := time.Now()
		err = c.Next()
		end := time.Now()

		log := &winner_logger.AccessLog{
			Method:    c.Method(),
			Status:    int32(c.Response().StatusCode()),
			BeginTime: start.Unix(),
			EndTime:   end.Unix(),
			Referer:   c.GetRespHeader("referer"),
			HttpHost:  string(c.Request().Host()),
			Interface: string(c.Request().URI().Path()),
			ReqQuery:  string(c.Request().URI().QueryString()),
			ReqBody:   c.Request().PostArgs().String(),
			// ResBody:   string(c.Response().Body()),
			ClientIp:  c.IP(),
			UserAgent: c.GetRespHeader("user-agent"),
			ReqId:     c.GetRespHeader("X-Request-ID"),
			Headers:   c.GetRespHeader("token"),
		}

		if ops.LogResponse {
			log.ResBody = string(c.Response().Body())
		}

		ops.Logger.Access(log)

		return err
	}
}
