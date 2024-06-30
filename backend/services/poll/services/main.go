package main

import (
	"context"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constans/config"
	"github.com/BigNoseCattyHome/aorb/backend/utils/extra/tracing"
	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
	"github.com/sirupsen/logrus"
	"net"
)

func main() {
	tp, err := tracing.SetTraceProvider(config.PollRpcServerName)

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

	log := logging.LogService(config.PollProcessorRpcServiceName)
	lis, err := net.Listen("tcp", config.Conf.Pod.PodIp+config.PollRpcServerAddr)

	if err != nil {
		log.Panicf("Rpc %s listen happens error: %v", config.PollRpcServerName, err)
	}

}
