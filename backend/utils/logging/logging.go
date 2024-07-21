package logging

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var hostname string

func init() {
	// 获取主机名
	hostname, _ = os.Hostname()
	fmt.Println(hostname)

	// 获取当前工作目录
	curDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	// 根据配置文件设置日志级别
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

	// 设置日志文件路径
	filePath := path.Join(curDir, config.Conf.Log.LogPath, "aorb.log")

	// 创建日志文件夹
	dir := path.Dir(filePath)
	if err := os.MkdirAll(dir, os.FileMode(0755)); err != nil {
		panic(err)
	}

	// 打开日志文件
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	// 设置日志格式为 JSON
	log.SetFormatter(&log.JSONFormatter{})
	// 添加日志钩子
	log.AddHook(logTraceHook{})
	// 设置日志输出到文件和标准输出
	log.SetOutput(io.MultiWriter(f, os.Stdout))

	// 初始化 Logger，包含主机名和 Pod IP
	Logger = log.WithFields(log.Fields{
		"Hostname": hostname,
		"Pod":      config.Conf.Pod.PodIp,
	})
}

// logTraceHook 是一个日志钩子，用于在日志中添加 trace 和 span 信息
type logTraceHook struct{}

func (t logTraceHook) Levels() []log.Level {
	return log.AllLevels
}

func (t logTraceHook) Fire(entry *log.Entry) error {
	ctx := entry.Context
	if ctx == nil {
		return nil
	}

	span := trace.SpanFromContext(ctx)
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

// LogService 返回一个带有服务名称的日志记录器
func LogService(name string) *log.Entry {
	return Logger.WithFields(log.Fields{
		"Service": name,
	})
}

// SetSpanError 记录 span 错误并设置状态为错误
func SetSpanError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, "Internal Error")
}

// SetSpanErrorWithDesc 记录 span 错误并设置状态为错误，带有描述信息
func SetSpanErrorWithDesc(span trace.Span, err error, desc string) {
	span.RecordError(err)
	span.SetStatus(codes.Error, desc)
}

// SetSpanWithHostname 在 span 中设置主机名和 Pod IP
func SetSpanWithHostname(span trace.Span) {
	span.SetAttributes(attribute.String("hostname", hostname))
	span.SetAttributes(attribute.String("podIP", config.Conf.Pod.PodIp))
}
