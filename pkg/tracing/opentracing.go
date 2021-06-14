package tracing

import (
	"io"

	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/util/opentelemetry"
	"github.com/micro/micro/v3/util/opentelemetry/jaeger"
)

func SetupOpentracing(serviceName string) io.Closer {
	c, _ := config.Get("jaegeraddress")
	openTracer, closer, err := jaeger.New(
		opentelemetry.WithServiceName(serviceName),
		opentelemetry.WithTraceReporterAddress(c.String("localhost:6831")),
	)
	if err != nil {
		logger.Fatalf("Error configuring opentracing: %v", err)
	}
	logger.Infof("Configured jaeger to %s", c.String("localhost:6831"))

	opentelemetry.DefaultOpenTracer = openTracer
	return closer
}
