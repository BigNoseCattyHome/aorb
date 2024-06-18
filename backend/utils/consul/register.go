package consul

// 用于与Consul服务发现和配置共享系统进行交互

import (
	"fmt"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constans/config"
	"github.com/google/uuid"
	capi "github.com/hashicorp/consul/api"
	logger "github.com/sirupsen/logrus"
	"strconv"
)

var consulClient *capi.Client

func init() {
	cfg := capi.DefaultConfig()
	cfg.Address = fmt.Sprintf(config.Conf.Consul.Host + ":" + config.Conf.Consul.Port)
	if c, err := capi.NewClient(cfg); err != nil {
		consulClient = c
		return
	} else {
		logger.Panicf("Connect Consul happens error: %v", err)
	}
}

func RegisterConsul(name string, port string) error {
	parsedPort, err := strconv.Atoi(port[1:])
	logger.WithFields(logger.Fields{
		"name": name,
		"port": parsedPort,
	}).Infof("Services Register Consul")
	if err != nil {
		return err
	}
	reg := &capi.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s", name, uuid.New().String()[:5]),
		Name:    name,
		Address: config.Conf.PodIP.PodIpAddress,
		Port:    parsedPort,
		Check: &capi.AgentServiceCheck{
			Interval:                       "5s",
			Timeout:                        "5s",
			GRPC:                           fmt.Sprintf("%s:%d", config.Conf.PodIP.PodIpAddress, parsedPort),
			GRPCUseTLS:                     false,
			DeregisterCriticalServiceAfter: "30s",
		},
	}
	if err := consulClient.Agent().ServiceRegister(reg); err != nil {
		return err
	}
	return nil
}
