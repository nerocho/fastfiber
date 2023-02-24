package fastfiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
	"github.com/gookit/validate/locales/zhcn"
)

func init() {
	zhcn.RegisterGlobal()

	validate.Config(func(opt *validate.GlobalOption) {
		opt.StopOnError = false
		opt.SkipOnEmpty = false
	})
}

// 绑定请求参数并返回校验的错误（一个）
func BindAndOneErr(ctx *fiber.Ctx, ptr interface{}) error {

	if ctx.Route().Method == "GET" {
		if err := ctx.QueryParser(ptr); err != nil {
			return err
		}
	} else {
		if err := ctx.BodyParser(ptr); err != nil {
			return err
		}
	}

	v := validate.New(ptr)

	if !v.Validate() {
		return v.Errors.OneError()
	}

	return nil
}

// 绑定请求参数并返回校验的错误（所有）
func BindAndAllErr(ctx *fiber.Ctx, ptr interface{}) error {

	if ctx.Route().Method == "GET" {
		if err := ctx.QueryParser(ptr); err != nil {
			return err
		}
	} else {
		if err := ctx.BodyParser(ptr); err != nil {
			return err
		}
	}

	v := validate.New(ptr)

	if !v.Validate() {
		return v.Errors
	}

	return nil
}
