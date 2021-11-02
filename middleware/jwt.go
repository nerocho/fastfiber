package middleware

import (
	"errors"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/nerocho/fastfiber/utils/jwt"
)

type JwtConfig struct {
	Secret        string
	AppName       string
	BusinessLines string
	Expire        time.Duration
	Next          func(c *fiber.Ctx) bool
}

// jwt 验签中间件，默认自动从header中依次获取token
// Secret        string 秘钥
// Next          func(c *fiber.Ctx) bool
func JWT(cfg JwtConfig) fiber.Handler {
	if len(cfg.Secret) == 0 {
		log.Fatal("Jwt.Secret不能为空")
	}

	return func(c *fiber.Ctx) error {
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		token := c.Get("token")
		if token == "" {
			return errors.New("非法的请求,token不存在")
		}

		if userId, err := jwt.ParseToken(token, cfg.Secret); err != nil {
			return err
		} else {
			c.Locals("UserId", userId) // 设置用户id
		}

		return c.Next()
	}
}
