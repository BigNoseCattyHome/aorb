package main

// grpc auth服务器主入口

import (
	"context"
	"net"
	"net/http"
	"os"
	"syscall"

	authController "github.com/BigNoseCattyHome/aorb/backend/go-services/auth/handlers"
	authRpc "github.com/BigNoseCattyHome/aorb/backend/rpc/auth"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	"github.com/BigNoseCattyHome/aorb/backend/utils/consul"
	"github.com/BigNoseCattyHome/aorb/backend/utils/extra/tracing"
	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
	"github.com/BigNoseCattyHome/aorb/backend/utils/prom"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

var log = logging.LogService(config.AuthRpcServerName) // 使用logging库，添加字段日志AuthRpcServer

func main() {
	// 设置分布式跟踪提供者
	tp, err := tracing.SetTraceProvider(config.AuthRpcServerName)
	if err != nil {
		// 如果设置跟踪失败，记录错误并 panic
		log.WithFields(logrus.Fields{
			"err": err,
		}).Panicf("Error to set the trace")
	}
	// 确保在 main 函数结束时关闭跟踪提供者
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.WithFields(logrus.Fields{
				"err": err,
			}).Errorf("Error to set the trace")
		}
	}()

	// 创建 TCP 监听器
	lis, err := net.Listen("tcp", config.Conf.Pod.PodIp+config.AuthRpcServerAddr)
	if err != nil {
		// 如果监听失败，记录错误并 panic
		log.Panicf("Rpc %s listen happens error: %v", config.AuthRpcServerName, err)
	}

	// 初始化 Prometheus 监控指标
	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)
	// 注册监控指标
	reg := prom.Client
	reg.MustRegister(srvMetrics)

	// 创建 gRPC 服务器并配置拦截器
	s := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.ChainUnaryInterceptor(srvMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(prom.ExtractContext))),
		grpc.ChainStreamInterceptor(srvMetrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(prom.ExtractContext))),
	)

	// 将服务注册到 Consul
	if err := consul.RegisterConsul(config.AuthRpcServerName, config.AuthRpcServerAddr); err != nil {
		log.Panicf("Rpc %s register consul happens error for: %v", config.AuthRpcServerName, err)
	}
	log.Infof("Rpc %s is running at %s now", config.AuthRpcServerName, config.AuthRpcServerAddr)

	// 注册 gRPC 服务和健康检查服务
	var srv authController.AuthServiceImpl
	authRpc.RegisterAuthServiceServer(s, srv)
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())
	defer authController.CloseMQConn()

	// 再次注册到 Consul
	if err := consul.RegisterConsul(config.AuthRpcServerName, config.AuthRpcServerAddr); err != nil {
		log.Panicf("Rpc %s register consul happens error for: %v", config.AuthRpcServerName, err)
	}
	srv.New()
	srvMetrics.InitializeMetrics(s)

	// 使用 run.Group 管理 goroutines
	g := &run.Group{}
	g.Add(func() error {
		// 启动 gRPC 服务
		return s.Serve(lis)
	}, func(error) {
		// 优雅停止 gRPC 服务
		s.GracefulStop()
		s.Stop()
		log.Errorf("Rpc %s listen happens error for: %v", config.AuthRpcServerName, err)
	})

	// 启动 HTTP 服务器以提供 Prometheus 指标
	httpSrv := &http.Server{
		Addr: config.Conf.Pod.PodIp + config.Metrics,
	}
	g.Add(func() error {
		// 设置 HTTP 处理函数
		m := http.NewServeMux()
		m.Handle("/metrics", promhttp.HandlerFor(
			reg,
			promhttp.HandlerOpts{
				EnableOpenMetrics: true,
			},
		))
		httpSrv.Handler = m
		log.Infof("Promethus now running")
		return httpSrv.ListenAndServe()
	}, func(err error) {
		// 关闭 HTTP 服务器
		if err := httpSrv.Close(); err != nil {
			log.Errorf("Prometheus %s listen happens error for: %v", config.AuthRpcServerName, err)
		}
	})

	// 添加信号处理函数，处理 SIGINT 和 SIGTERM 信号
	g.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))

	// 运行所有 goroutines 并等待它们完成
	if err := g.Run(); err != nil {
		log.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Error when runing http server")
		os.Exit(1)
	}
}
