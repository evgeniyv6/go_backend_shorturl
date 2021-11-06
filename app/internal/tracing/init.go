package tracing

import (
	"fmt"
	"io"

	"github.com/evgeniyv6/go_backend_shorturl/app/internal/logger"

	"github.com/uber/jaeger-client-go/config"

	"github.com/opentracing/opentracing-go"
)

type (
	JaegerTracer struct {
		serviceName string
		logger      logger.ZapWrapper
	}
)

func NewJaegerTracer(name string, logger logger.ZapWrapper) JaegerTracer {
	return JaegerTracer{name, logger}
}

func (j *JaegerTracer) Init() (tracer opentracing.Tracer, closer io.Closer) {
	cfg := &config.Configuration{
		ServiceName: j.serviceName,
		Sampler:     &config.SamplerConfig{Type: "const", Param: 1},
		Reporter:    &config.ReporterConfig{LogSpans: true},
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(j.logger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}

	return tracer, closer
}
