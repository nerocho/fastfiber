package redispool

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	winner_logger "github.com/bfmTech/logger-go"
	redis "github.com/go-redis/redis/v8"
)

type TracingHook struct{}

type RedisHookSpan string

const (
	_RedisSpan               = "_RedisSpan"
	tag                      = "[redisInfo]"
	startTime  RedisHookSpan = "_winnerRedisStartTime"
	requestId                = "request_id"
)

var _ redis.Hook = TracingHook{}
var wl winner_logger.Logger

func NewRedisTracingHook(logger winner_logger.Logger) redis.Hook {
	wl = logger

	return TracingHook{}
}

func (TracingHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, startTime, time.Now()), nil
}

func (TracingHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {

	_reqId := ctx.Value(requestId)
	reqId, ok := _reqId.(string)
	if !ok {
		reqId = ""
	}

	_st := ctx.Value(startTime)
	startTime, ok := _st.(time.Time)
	if !ok {
		return nil
	}

	costSeconds := strconv.FormatInt(time.Since(startTime).Milliseconds(), 10)

	if err := cmd.Err(); err != nil && !errors.Is(err, redis.Nil) {
		wl.Error(err)
	} else {
		wl.Info(tag, "request_id:"+reqId, "cmd:"+cmd.String(), "cost:"+costSeconds)
	}

	return nil
}

func (TracingHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return context.WithValue(ctx, startTime, time.Now()), nil
}

func (TracingHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {

	_reqId := ctx.Value(requestId)
	reqId, ok := _reqId.(string)
	if !ok {
		reqId = ""
	}

	_st := ctx.Value(startTime)
	startTime, ok := _st.(time.Time)
	if !ok {
		return nil
	}

	costSeconds := strconv.FormatInt(time.Since(startTime).Milliseconds(), 10)

	pipeline := make([]string, len(cmds))

	for idx, cmd := range cmds {
		pipeline[idx] = cmd.String()
	}

	wl.Info(tag, "request_id:"+reqId, "pipelines:["+strings.Join(pipeline, ",")+"]", "cost:"+costSeconds)

	return nil

}
