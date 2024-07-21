package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

var (
	// connPool 用于存储 gRPC 连接的池
	connPool = make(map[string]*grpc.ClientConn)
	// connMu 用于保护 connPool 的读写锁
	connMu sync.RWMutex
)

// 使用logging库，添加字段日志Gateway
var log = logging.LogService(config.WebServerName)

// GatewayServer 是一个空的结构体，用于实现 gRPC 服务
type GatewayServer struct{}

// getServiceConn 获取指定服务的 gRPC 连接
func getServiceConn(consulClient *api.Client, serviceName string) (*grpc.ClientConn, error) {
	// 读锁保护 connPool 的读操作
	connMu.RLock()
	conn, exists := connPool[serviceName]
	connMu.RUnlock()
	if exists {
		return conn, nil
	}

	// 写锁保护 connPool 的写操作
	connMu.Lock()
	defer connMu.Unlock()

	// 再次检查连接是否已经存在，防止重复创建
	if conn, exists = connPool[serviceName]; exists {
		return conn, nil
	}

	// 使用 Consul 查找服务
	services, _, err := consulClient.Health().Service(serviceName, "", true, nil)
	if err != nil {
		log.Printf("service discovery failed: %v", err)
		return nil, fmt.Errorf("service discovery failed: %v", err)
	}
	if len(services) == 0 {
		log.Printf("no healthy instances found for service: %s", serviceName)
		return nil, fmt.Errorf("no healthy instances found for service: %s", serviceName)
	}

	// 选择第一个健康的服务实例
	service := services[0].Service

	// 创建 gRPC 连接
	conn, err = grpc.NewClient(
		fmt.Sprintf("%s:%d", service.Address, service.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		return nil, err
	}

	// 将连接存入连接池
	connPool[serviceName] = conn
	return conn, nil
}

// forwardInterceptor 是一个 gRPC 拦截器，用于转发请求到相应的服务
func forwardInterceptor(consulClient *api.Client) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 从方法名中提取服务名
		serviceName := extractServiceName(info.FullMethod)

		// 获取服务的 gRPC 连接
		conn, err := getServiceConn(consulClient, serviceName)
		if err != nil {
			return nil, err
		}

		// 创建 outgoing 上下文
		outgoingCtx := metadata.NewOutgoingContext(ctx, metadata.MD{})

		// 调用远程服务
		var header, trailer metadata.MD
		var resp interface{}
		err = conn.Invoke(outgoingCtx, info.FullMethod, req, &resp, grpc.Header(&header), grpc.Trailer(&trailer))
		if err != nil {
			return nil, fmt.Errorf("RPC failed: %v", err)
		}

		return resp, nil
	}
}

// extractServiceName 从 gRPC 方法名中提取服务名
func extractServiceName(fullMethod string) string {
	parts := strings.Split(fullMethod, "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}

func main() {
	// 创建 Consul 客户端
	consulConfig := api.DefaultConfig()
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
	}

	// 监听指定端口
	lis, err := net.Listen("tcp", ":37000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 创建 gRPC 服务器，并注册拦截器
	s := grpc.NewServer(
		grpc.UnaryInterceptor(forwardInterceptor(consulClient)),
	)

	// 注册反射服务
	reflection.Register(s)

	// 启动 gRPC 服务器
	go func() {
		log.Printf("starting gRPC server on %s", lis.Addr().String())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// 监听系统信号，用于优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 优雅关闭 gRPC 服务器
	s.GracefulStop()

	// 关闭所有 gRPC 连接
	for _, conn := range connPool {
		conn.Close()
	}

	log.Println("Server exiting")
}
