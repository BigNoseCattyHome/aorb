package consul

// 用于与 Consul 服务发现和配置共享系统进行交互

import (
	"fmt"
	"strconv"

	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
	"github.com/google/uuid"
	capi "github.com/hashicorp/consul/api"
	"github.com/sirupsen/logrus"
)

var consulClient *capi.Client

// 初始化 Consul 客户端
// 从配置中读取 Consul 服务器的地址，并创建一个 Consul 客户端实例
func init() {
	cfg := capi.DefaultConfig()
	cfg.Address = config.Conf.Consul.Addr // 从配置中获取 Consul 服务器地址
	client, err := capi.NewClient(cfg)
	if err != nil {
		logging.Logger.Panicf("Error connecting to Consul: %v", err)
		return
	}
	consulClient = client
}

// RegisterConsul 将服务注册到 Consul
// 参数 name 为服务的名称，port 为服务的端口
func RegisterConsul(name string, port string) error {
	// 解析端口号（去掉端口号前的冒号）
	parsedPort, err := strconv.Atoi(port[1:])
	if err != nil {
		return fmt.Errorf("invalid port format: %v", err)
	}

	// 生成服务名称，包含匿名名称前缀
	fullServiceName := config.Conf.Consul.AnonymityName + name

	// 创建服务注册信息
	reg := &capi.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s", fullServiceName, uuid.New().String()[:5]), // 服务ID，使用UUID避免冲突
		Name:    fullServiceName,                                                // 服务名称
		Address: config.Conf.Pod.PodIp,                                          // 服务的IP地址
		Port:    parsedPort,                                                     // 服务的端口
		Check: &capi.AgentServiceCheck{
			Interval:                       "5s",                                                    // 健康检查间隔
			Timeout:                        "5s",                                                    // 健康检查超时时间
			GRPC:                           fmt.Sprintf("%s:%d", config.Conf.Pod.PodIp, parsedPort), // gRPC健康检查地址
			GRPCUseTLS:                     false,                                                   // 不使用TLS
			DeregisterCriticalServiceAfter: "30s",                                                   // 30秒内未响应则注销服务
		},
	}

	// 注册服务到 Consul
	if err := consulClient.Agent().ServiceRegister(reg); err != nil {
		return fmt.Errorf("failed to register service: %v", err)
	}

	logging.Logger.WithFields(logrus.Fields{
		"name": name,
		"port": parsedPort,
	}).Infof("Service registered to Consul successfully")

	return nil
}
