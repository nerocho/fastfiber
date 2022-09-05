package fastfiber

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"

	"github.com/nerocho/fastfiber/utils/response"
)

// 获取环境变量
func GetEnv(envName string, defaultValues ...string) string {
	e := os.Getenv(envName)
	if e == "" && len(defaultValues) > 0 {
		e = defaultValues[0]
	}
	return e
}

// 安全运行goroutine
func GoSafe(fn func()) {
	go runSafe(fn)
}

func runSafe(fn func()) {
	defer func() {
		if err := recover(); err != nil {
			switch v := err.(type) {
			case string:
				Logger.Error(errors.New(v))
			case error:
				Logger.Error(v)
			default:
			}
		}
	}()

	fn()
}

// 带时限的异步执行 返回true为超时、false为未超时
func TaskWithTimeout(task func() error, duration time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	done := make(chan struct{}, 1)

	GoSafe(func() {
		task()
		done <- struct{}{}
	})

	select {
	case <-done:
		return false
	case <-ctx.Done():
		return true
	}
}

// 错误处理
func ErrorHandler(ctx *fiber.Ctx, err error) error {
	var code int
	switch e := err.(type) {
	case *fiber.Error: // 判断是否为框架错误
		code = e.Code
	case validate.Errors: //如果是验证类错误，覆盖状态码
		code = 400
	default:
		code = fiber.StatusInternalServerError
	}
	return response.MakeRes(ctx).Code(code).Msg(err.Error()).Send()
}
