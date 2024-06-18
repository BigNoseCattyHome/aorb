package grpc

import (
	"fmt"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constans/config"
	logger "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"time"
)

func Connect(serviceName string) (conn *grpc.ClientConn) {
	kacp := keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             1 * time.Second,  // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: false,            // send pings even without active streams
	}
	consulAddr := config.Conf.Consul.Port + ":" + config.Conf.Consul.Port
	conn, err := grpc.Dial(
		fmt.Sprintf("consul://%s/%s?wait=15s", consulAddr, config.Conf.Consul.AnonymityName+serviceName),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithKeepaliveParams(kacp),
	)
	logger.Debugf("connect")
	if err != nil {
		logger.WithFields(logger.Fields{
			"service": config.Conf.Consul.AnonymityName + serviceName,
			"err":     err,
		}).Errorf("Cannot connect to %v service", config.Conf.Consul.AnonymityName+serviceName)
	}
	return
}
