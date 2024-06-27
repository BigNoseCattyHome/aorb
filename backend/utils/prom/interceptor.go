package prom

import (
	"context"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/otel/trace"
)

func ExtractContext(ctx context.Context) prometheus.Labels {
	if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
		return prometheus.Labels{
			"traceId": span.TraceID().String(),
			"spanId":  span.SpanID().String(),
			"podId":   config.Conf.Pod.PodIpAddress,
		}
	}
	return nil
}
