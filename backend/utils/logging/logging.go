package logging

import (
	"fmt"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constans/config"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"io"
	"os"
	"path"
)

var hostname string

func init() {
	hostname, _ = os.Hostname()

	fmt.Println(hostname)

	switch config.Conf.Log.LoggerLevel {
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "WARN", "WARNING":
		log.SetLevel(log.WarnLevel)
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "FATAL":
		log.SetLevel(log.FatalLevel)
	case "TRACE":
		log.SetLevel(log.TraceLevel)
	}

	filePath := path.Join("/Users", "jinziguan", "desktop", "aorb", "log", "aorb.log")
	dir := path.Dir(filePath)
	if err := os.MkdirAll(dir, os.FileMode(0755)); err != nil {
		panic(err)
	}

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	log.SetFormatter(&log.JSONFormatter{})
	log.AddHook(logTraceHook{})
	log.SetOutput(io.MultiWriter(f, os.Stdout))

	Logger = log.WithFields(log.Fields{
		"Hostname": hostname,
		"Pod":      config.Conf.Pod.PodIpAddress,
	})
}

type logTraceHook struct{}

func (t logTraceHook) Levels() []log.Level { return log.AllLevels }

func (t logTraceHook) Fire(entry *log.Entry) error {
	ctx := entry.Context
	if ctx == nil {
		return nil
	}

	span := trace.SpanFromContext(ctx)
	//if !span.IsRecording() {
	//	return nil
	//}

	sCtx := span.SpanContext()
	if sCtx.HasTraceID() {
		entry.Data["trace_id"] = sCtx.TraceID().String()
	}
	if sCtx.HasSpanID() {
		entry.Data["span_id"] = sCtx.SpanID().String()
	}

	return nil
}

var Logger *log.Entry

func LogService(name string) *log.Entry {
	return Logger.WithFields(log.Fields{
		"Service": name,
	})
}

func SetSpanError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, "Internal Error")
}

func SetSpanErrorWithDesc(span trace.Span, err error, desc string) {
	span.RecordError(err)
	span.SetStatus(codes.Error, desc)
}

func SetSpanWithHostname(span trace.Span) {
	span.SetAttributes(attribute.String("hostname", hostname))
	span.SetAttributes(attribute.String("podIP", config.Conf.Pod.PodIpAddress))
}
