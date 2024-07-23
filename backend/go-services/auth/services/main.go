package main

// grpc auth服务器主入口

import (
	"context"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/auth"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/database"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/redis"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	redis2 "github.com/redis/go-redis/v9"
	"github.com/willf/bloom"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net"
	"net/http"
	"os"
	"syscall"

	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	"github.com/BigNoseCattyHome/aorb/backend/utils/consul"
	"github.com/BigNoseCattyHome/aorb/backend/utils/extra/tracing"
	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
	"github.com/BigNoseCattyHome/aorb/backend/utils/prom"
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/oklog/run"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	tp, err := tracing.SetTraceProvider(config.AuthRpcServerName)

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

	// Configure Pyroscope
	//profiling.InitPyroscope("AorB.AuthService")

	log := logging.LogService(config.AuthRpcServerName)
	lis, err := net.Listen("tcp", config.Conf.Pod.PodIp+config.AuthRpcServerAddr)

	if err != nil {
		log.Panicf("Rpc %s listen happens error: %v", config.AuthRpcServerName, err)
	}

	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)

	reg := prom.Client
	reg.MustRegister(srvMetrics)

	// Create a new Bloom filter with a target false positive rate of 0.1%
	BloomFilter = bloom.NewWithEstimates(10000000, 0.001) // assuming we have 1 million users

	// Initialize BloomFilter from database
	collection := database.MongoDbClient.Database("aorb").Collection("users")
	findOptions := options.Find().SetProjection(bson.D{{"username", 1}})
	cur, err := collection.Find(context.Background(), bson.D{}, findOptions)
	var results []bson.M
	if err = cur.All(context.Background(), &results); err != nil {
		log.Fatal(err)
	}
	for _, result := range results {
		BloomFilter.AddString(result["username"].(string))
	}

	// Create a go routine to receive redis message and add it to BloomFilter
	go func() {
		pubSub := redis.Client.Subscribe(context.Background(), config.BloomRedisChannel)
		defer func(pubSub *redis2.PubSub) {
			err := pubSub.Close()
			if err != nil {
				log.Panicf("Closing redis pubsub happend error: %s", err)
			}
		}(pubSub)

		_, err := pubSub.ReceiveMessage(context.Background())
		if err != nil {
			log.Panicf("Reveiving message from redis happens error: %s", err)
			panic(err)
		}

		ch := pubSub.Channel()
		for msg := range ch {
			log.Infof("Add user name to BloomFilter: %s", msg.Payload)
			BloomFilter.AddString(msg.Payload)
		}
	}()

	s := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.ChainUnaryInterceptor(srvMetrics.UnaryServerInterceptor(grpcprom.WithExemplarFromContext(prom.ExtractContext))),
		grpc.ChainStreamInterceptor(srvMetrics.StreamServerInterceptor(grpcprom.WithExemplarFromContext(prom.ExtractContext))),
	)

	if err := consul.RegisterConsul(config.AuthRpcServerName, config.AuthRpcServerAddr); err != nil {
		log.Panicf("Rpc %s register consul happens error for: %v", config.AuthRpcServerName, err)
	}
	log.Infof("Rpc %s is running at %s now", config.AuthRpcServerName, config.AuthRpcServerAddr)

	var srv AuthServiceImpl
	auth.RegisterAuthServiceServer(s, srv)
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())

	srv.New()
	srvMetrics.InitializeMetrics(s)

	g := &run.Group{}
	g.Add(func() error {
		return s.Serve(lis)
	}, func(err error) {
		s.GracefulStop()
		s.Stop()
		log.Errorf("Rpc %s listen happens error for: %v", config.AuthRpcServerName, err)
	})

	httpSrv := &http.Server{Addr: config.Conf.Pod.PodIp + config.AuthMetrics}
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
			log.Errorf("Prometheus %s listen happens error for: %v", config.AuthRpcServerName, err)
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
