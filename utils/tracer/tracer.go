package tracer

import (
	"io"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go/config"
)

func NewJaegerTracer(serviceName, agentHostPort string) (opentracing.Tracer, io.Closer, error) {
	// jaeger client 配置项
	cfg := &config.Configuration{
		ServiceName: serviceName,
		// 固定采样
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		// 刷新缓冲区的频率，上报的 agent 地址等
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  agentHostPort,
		},
	}
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		return nil, nil, err
	}
	// 设置全局 tracer 对象
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer, nil
}
