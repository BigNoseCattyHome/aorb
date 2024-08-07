package tracing

import (
	"context"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	trace2 "go.opentelemetry.io/otel/trace"
)

var Tracer trace2.Tracer

func SetTraceProvider(name string) (*trace.TracerProvider, error) {
	client := otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(config.Conf.Tracing.EndPoint),
		otlptracehttp.WithInsecure(),
	)
	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Can not init otel !")
		return nil, err
	}

	var sampler trace.Sampler
	if config.Conf.Tracing.State == "disable" {
		sampler = trace.NeverSample()
	} else {
		sampler = trace.TraceIDRatioBased(config.Conf.Tracing.Sampler)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String(name),
			),
		),
		trace.WithSampler(sampler),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	Tracer = otel.Tracer(name)
	return tp, nil
}
