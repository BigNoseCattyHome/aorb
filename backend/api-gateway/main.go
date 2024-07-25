package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"sync"
	"syscall"

	"github.com/BigNoseCattyHome/aorb/backend/rpc/auth"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/comment"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/poll"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/recommend"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/user"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/vote"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

var (
	connPool = make(map[string]*grpc.ClientConn) // connPool 用于存储 gRPC 连接的池
	connMu   sync.RWMutex                        // connMu 用于保护 connPool 的读写锁

)

// 使用logging库，添加字段日志Gateway
var log = logging.LogService(config.WebServerName)

// TODO 因为每一次增加新的微服务都需要在GatewayServer中新增一个Unimplemented的服务，并且要在main函数中注册；后期可以使用一个通用的网关服务，通过反射的方式来实现
/*
// DynamicGatewayServer 是一个通用的网关服务器
type DynamicGatewayServer struct{}

// 实现 grpc.UnknownServiceHandler 接口
func (s *DynamicGatewayServer) unknownServiceHandler(srv interface{}, stream grpc.ServerStream) error {
    fullMethodName, ok := grpc.MethodFromServerStream(stream)
    if !ok {
        return fmt.Errorf("failed to get method name from stream")
		}

		// 从方法名中提取服务名
		grpcServiceName, consulServiceName := extractServiceName(fullMethodName)

		// 使用 consulServiceName 获取服务连接
		conn, err := getServiceConn(consulClient, consulServiceName)
		if err != nil {
			return err
			}

			// 转发请求到目标服务
			return grpc.ForwardServerStream(stream, conn, fullMethodName)
			}
*/

// GatewayServer 用于实现 gRPC 服务
type GatewayServer struct {
	auth.UnimplementedAuthServiceServer
	user.UnimplementedUserServiceServer
	comment.UnimplementedCommentServiceServer
	poll.UnimplementedPollServiceServer
	vote.UnimplementedVoteServiceServer
	recommend.UnimplementedRecommendServiceServer
}

// gRPC 服务名到 Consul 服务名的映射
var serviceNameMapping = map[string]string{
	"rpc.auth.AuthService":           config.AuthRpcServerName,
	"rpc.user.UserService":           config.UserRpcServerName,
	"rpc.comment.CommentService":     config.CommentRpcServerName,
	"rpc.vote.VoteService":           config.VoteRpcServerName,
	"rpc.poll.QuestionService":       config.PollRpcServerName,
	"rpc.recommend.RecommendService": config.RecommendRpcServerName,
}

// getServiceConn 获取指定服务的 gRPC 连接
func getServiceConn(consulClient *api.Client, consulServiceName string) (*grpc.ClientConn, error) {
	// 读锁保护 connPool 的读操作
	connMu.RLock()
	conn, exists := connPool[consulServiceName]
	connMu.RUnlock()
	if exists {
		return conn, nil
	}

	// 写锁保护 connPool 的写操作
	connMu.Lock()
	defer connMu.Unlock()

	// 再次检查连接是否已经存在，防止重复创建
	if conn, exists = connPool[consulServiceName]; exists {
		return conn, nil
	}

	// 使用 Consul 查找服务
	log.Infof("Querying Consul for service: %s", consulServiceName)
	services, _, err := consulClient.Health().Service(consulServiceName, "", true, nil)
	// services, _, err := consulClient.Health().Service(serviceName, "", true, nil)
	if err != nil {
		log.Errorf("service discovery failed: %v", err)
		return nil, fmt.Errorf("service discovery failed: %v", err)
	}
	if len(services) == 0 {
		log.Errorf("no healthy instances found for service: %s", consulServiceName)
		return nil, fmt.Errorf("no healthy instances found for service: %s", consulServiceName)
	}

	// 选择第一个健康的服务实例
	service := services[0].Service
	log.Infof("Found healthy instance for service %s: %s:%d", consulServiceName, service.Address, service.Port)

	// 创建 gRPC 连接
	conn, err = grpc.NewClient(
		fmt.Sprintf("%s:%d", service.Address, service.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		log.Errorf("failed to create gRPC connection: %v", err)
		return nil, err
	}

	// 将连接存入连接池
	connPool[consulServiceName] = conn
	return conn, nil
}

// forwardInterceptor 是一个 gRPC 拦截器，用于转发请求到相应的服务
func forwardInterceptor(consulClient *api.Client) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 从方法名中提取服务名
		log.Debug("Start to forward request to service")
		grpcServiceName, consulServiceName := extractServiceName(info.FullMethod)
		log.Infof("Received request for gRPC service: %s (Consul service: %s), method: %s", grpcServiceName, consulServiceName, info.FullMethod)

		// 使用 consulServiceName 获取服务连接
		conn, err := getServiceConn(consulClient, consulServiceName)
		if err != nil {
			log.Errorf("failed to get service connection: %v", err)
			return nil, err
		}

		// 创建 outgoing 上下文
		outgoingCtx := metadata.NewOutgoingContext(ctx, metadata.MD{})

		// 获取正确的响应类型
		respType := getResponseType(info.FullMethod)
		if respType == nil {
			return nil, fmt.Errorf("unknown method: %s", info.FullMethod)
		}

		// 创建正确类型的响应对象
		resp := reflect.New(respType.Elem()).Interface()

		// 调用远程服务
		var header, trailer metadata.MD
		log.Infof("Forwarding request to service: %s, method: %s", consulServiceName, info.FullMethod)
		err = conn.Invoke(outgoingCtx, info.FullMethod, req, resp, grpc.Header(&header), grpc.Trailer(&trailer))
		if err != nil {
			log.Errorf("RPC failed: %v", err)
			return nil, fmt.Errorf("RPC failed: %v", err)
		}

		log.Infof("Successfully forwarded request to service: %s, method: %s", consulServiceName, info.FullMethod)
		log.Debugf("Response type: %T", resp)

		// 返回接口类型
		return resp, nil
	}
}

// getResponseType 根据方法名返回对应的响应类型
func getResponseType(method string) reflect.Type {
	switch method {
	case "/rpc.auth.AuthService/Register":
		return reflect.TypeOf((*auth.RegisterResponse)(nil))
	case "/rpc.auth.AuthService/Login":
		return reflect.TypeOf((*auth.LoginResponse)(nil))
	case "/rpc.auth.AuthService/Verify":
		return reflect.TypeOf((*auth.VerifyResponse)(nil))
	case "/rpc.auth.AuthService/Refresh":
		return reflect.TypeOf((*auth.RefreshResponse)(nil))
	case "/rpc.auth.AuthService/Logout":
		return reflect.TypeOf((*auth.LogoutResponse)(nil))

	case "/rpc.user.UserService/GetUserInfo":
		return reflect.TypeOf((*user.UserResponse)(nil))
	case "/rpc.user.UserService/GetUserExistInformation":
		return reflect.TypeOf((*user.UserExistResponse)(nil))

	case "/rpc.comment.CommentService/ActionComment":
		return reflect.TypeOf((*comment.ActionCommentResponse)(nil))
	case "/rpc.comment.CommentService/ListComment":
		return reflect.TypeOf((*comment.GetCommentResponse)(nil))
	case "/rpc.comment.CommentService/CountComment":
		return reflect.TypeOf((*comment.CountCommentResponse)(nil))

	case "/rpc.poll.PollService/CreatePoll":
		return reflect.TypeOf((*poll.CreatePollResponse)(nil))
	case "/rpc.poll.PollService/GetPoll":
		return reflect.TypeOf((*poll.GetPollResponse)(nil))
	case "/rpc.poll.PollService/ListPoll":
		return reflect.TypeOf((*poll.ListPollResponse)(nil))

	case "/rpc.vote.VoteService/CreateVote":
		return reflect.TypeOf((*vote.CreateVoteResponse)(nil))
	case "/rpc.vote.VoteService/GetVote":
		return reflect.TypeOf((*vote.GetVoteCountResponse)(nil))

	case "/rpc.recommend.RecommendService/GetRecommendInformation":
		return reflect.TypeOf((*recommend.RecommendResponse)(nil))
	case "/rpc.recommend.RecommendService/RegisterRecommendUser":
		return reflect.TypeOf((*recommend.RecommendRegisterResponse)(nil))
	default:
		return nil
	}
}

// extractServiceName 从 gRPC 方法名中提取服务名
func extractServiceName(fullMethod string) (grpcServiceName, consulServiceName string) {
	parts := strings.Split(fullMethod, "/")
	if len(parts) < 2 {
		log.Errorf("Failed to extract service name from method: %s", fullMethod)
		return "", ""
	}
	grpcServiceName = parts[1]
	consulServiceName, ok := serviceNameMapping[grpcServiceName]
	if !ok {
		log.Errorf("No Consul service name mapping for gRPC service: %s", grpcServiceName)
		return grpcServiceName, grpcServiceName // 默认使用 gRPC 服务名
	}
	log.Infof("Mapped gRPC service %s to Consul service %s", grpcServiceName, consulServiceName)
	return grpcServiceName, consulServiceName
}
func main() {
	// 创建 Consul 客户端
	consulConfig := api.DefaultConfig()
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		log.Fatalf("Failed to create Consul client: %v", err)
	}
	log.Info("Consul client created")

	// 监听指定端口
	lis, err := net.Listen("tcp", ":37000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Info("Gateway is listening on :37000")

	// 创建 gRPC 服务器，并注册拦截器
	s := grpc.NewServer(
		grpc.UnaryInterceptor(forwardInterceptor(consulClient)),
	)

	// 注册服务
	// 注意是先注册服务，再启动gRPC服务器
	auth.RegisterAuthServiceServer(s, &GatewayServer{})
	user.RegisterUserServiceServer(s, &GatewayServer{})
	comment.RegisterCommentServiceServer(s, &GatewayServer{})
	poll.RegisterPollServiceServer(s, &GatewayServer{})
	vote.RegisterVoteServiceServer(s, &GatewayServer{})
	recommend.RegisterRecommendServiceServer(s, &GatewayServer{})

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
