package middleware

import (
	"github.com/gofiber/fiber/v2"

	"github.com/nerocho/fastfiber/utils/response"
)

// 404处理
func NotFound() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return response.MakeRes(c).Code(fiber.StatusNotFound).Msg("Not Found").Send()
	}
}
