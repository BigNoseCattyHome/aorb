package main

import (
	"context"
	"github.com/BigNoseCattyHome/aorb/backend/go-services/vote/services"
	votePb "github.com/BigNoseCattyHome/aorb/backend/rpc/vote"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	"github.com/BigNoseCattyHome/aorb/backend/utils/consul"
	"github.com/BigNoseCattyHome/aorb/backend/utils/extra/tracing"
	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
	"github.com/BigNoseCattyHome/aorb/backend/utils/prom"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"net/http"
	"os"
	"syscall"
)

func main() {
	tp, err := tracing.SetTraceProvider(config.VoteRpcServerName)

	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Panicf("Error to set the trace")
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logging.Logger.WithFields(logrus.Fields{
				"err": err,
			}).Errorf("Error to set the trace")
		}
	}()

	//profiling.InitPyroscope("AorB.VoteService")

	log := logging.LogService(config.VoteRpcServerName)
	lis, err := net.Listen("tcp", config.Conf.Pod.PodIp+config.VoteRpcServerAddr)

	if err != nil {
		log.Panicf("Rpc %s listen happens error: %v", config.VoteRpcServerName, err)
	}

	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)
	reg := prom.Client
	reg.MustRegister(srvMetrics)

	s := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.ChainUnaryInterceptor(srvMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(prom.ExtractContext))),
		grpc.ChainStreamInterceptor(srvMetrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(prom.ExtractContext))),
	)

	if err := consul.RegisterConsul(config.VoteRpcServerName, config.VoteRpcServerAddr); err != nil {
		log.Panicf("Rpc %s register consul happens error for: %v", config.VoteRpcServerName, err)
	}
	log.Infof("Rpc %s is running at %s now", config.VoteRpcServerName, config.VoteRpcServerAddr)

	var srv services.VoteServiceImpl
	votePb.RegisterVoteServiceServer(s, srv)
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())
	defer services.CloseMQConn()
	if err := consul.RegisterConsul(config.VoteRpcServerName, config.VoteRpcServerAddr); err != nil {
		log.Panicf("Rpc %s register consul happens error for: %v", config.VoteRpcServerName, err)
	}
	srv.New()
	srvMetrics.InitializeMetrics(s)

	g := &run.Group{}
	g.Add(func() error {
		return s.Serve(lis)
	}, func(err error) {
		s.GracefulStop()
		s.Stop()
		log.Errorf("Rpc %s listen happens error for: %v", config.VoteRpcServerName, err)
	})

	httpSrv := &http.Server{Addr: config.Conf.Pod.PodIp + config.VoteMetrics}
	g.Add(func() error {
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
	}, func(error) {
		if err := httpSrv.Close(); err != nil {
			log.Errorf("Prometheus %s listen happens error for: %v", config.VoteRpcServerName, err)
		}
	})

	g.Add(run.SignalHandler(context.Background(), syscall.SIGINT, syscall.SIGTERM))

	if err := g.Run(); err != nil {
		log.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Error when runing http server")
		os.Exit(1)
	}
}
