package response

import (
	"github.com/gofiber/fiber/v2"
)

const (
	msg_success = "success"
	msg_fail    = "fail"
	msg_warn    = "warn"
)

type response struct {
	ctx  *fiber.Ctx
	code int
	msg  string
}

/**
MakeRes 返回消息到客户端
示例：
MakeRes(c).Code(somecode).Msg(somemsg).Send(somedata) 返回一个自定义消息
MakeRes(c).Msg(somemsg).Send() 只返回一个提示
MakeRes(c).SendSuccess(somedata) 简写返回一段数据
MakeRes(c).SendSuccess(somemsg,somedata) 简写返回一段数据并覆盖默认的消息提示
MakeRes(c).SendFail 同上
MakeRes(c).SendWarn 同上
*/
func MakeRes(ctx *fiber.Ctx) *response {
	return &response{
		ctx: ctx,
	}
}

//设置报文代码，且会自动根据代码设置httpcode
func (r *response) Code(code int) *response {
	//设置默认code
	if code == 0 {
		code = 200
	}

	r.code = code

	switch code {
	case 200, 400, 500:
		r.ctx.Status(code)
	default:
		r.ctx.Status(fiber.StatusOK)
	}

	return r
}

// 设置消息提示
func (r *response) Msg(msg string) *response {
	// 设置默认 消息
	if msg == "" {
		switch r.code {
		case 400:
			msg = msg_warn
		case 500:
			msg = msg_fail
		default:
			msg = msg_success
		}
	}
	r.msg = msg
	return r
}

//发送数据
func (r *response) Send(args ...interface{}) error {
	if args != nil {
		return r.ctx.JSON(buildRes(r.code, r.msg, args[0]))
	} else {
		return r.ctx.JSON(buildRes(r.code, r.msg, nil))
	}
}

func (r *response) SendSuccess(args ...interface{}) error {
	msg, data := makeArgs(&args)
	return r.Code(200).Msg(msg).Send(*data)
}

func (r *response) SendWarn(args ...interface{}) error {
	msg, data := makeArgs(&args)
	return r.Code(500).Msg(msg).Send(*data)
}

func (r *response) SendFail(args ...interface{}) error {
	msg, data := makeArgs(&args)
	return r.Code(400).Msg(msg).Send(*data)
}

func buildRes(code int, msg string, data interface{}) fiber.Map {
	if code == 0 {
		code = 200
	}
	if msg == "" {
		msg = msg_success
	}
	return fiber.Map{
		"code": code,
		"msg":  msg,
		"data": data,
	}
}

func makeArgs(args *[]interface{}) (string, *interface{}) {
	var msg string
	var data interface{}
	for i := 0; i < len(*args); i++ {
		switch arg := (*args)[i].(type) {
		case string:
			if i == 0 {
				msg = arg
			} else {
				data = arg
			}
		default:
			data = arg
		}
	}
	return msg, &data
}
