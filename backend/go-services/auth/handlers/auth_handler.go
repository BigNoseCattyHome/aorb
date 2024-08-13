package handlers

import (
	"context"
	"os"

	"github.com/BigNoseCattyHome/aorb/backend/go-services/auth/conf"
	"github.com/BigNoseCattyHome/aorb/backend/go-services/auth/services"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/auth"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/user"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/strings"
	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"
)

var log = logging.LogService(config.AuthRpcServerName) // 使用logging库，添加字段日志AuthRpcServer
var conn *amqp.Connection                              // amqp.Connection用于连接RabbitMQ服务器
var channel *amqp.Channel                              // amqp.Channel用于与RabbitMQ服务器通信

// exitOnError 如果err不为nil，则panic
func exitOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// CloseMQConn 关闭RabbitMQ连接
func CloseMQConn() {
	if err := conn.Close(); err != nil {
		panic(err)
	}
	if err := channel.Close(); err != nil {
		panic(err)
	}
}

// AuthServiceImpl AuthService服务实现
type AuthServiceImpl struct {
	auth.AuthServiceServer
}

// 初始化
func (a AuthServiceImpl) Init() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	exitOnError(err)
	channel, err = conn.Channel()
	exitOnError(err)
}

// 创建AuthService服务实例
func (a AuthServiceImpl) New() {

}

// 登录
func (a AuthServiceImpl) Login(ctx context.Context, request *auth.LoginRequest) (*auth.LoginResponse, error) {

	// 解析请求
	login_request := auth.LoginRequest{
		Username:  request.Username,
		Password:  request.Password,
		DeviceId:  request.DeviceId,
		Nonce:     request.Nonce,
		Timestamp: request.Timestamp,
		Ipaddress: request.Ipaddress,
	}

	// 调用服务
	token, exp_token, refresh_token, simple_user, err := services.AuthenticateUser(ctx, &login_request)
	if err != nil {
		// 当出现预期中的错误（用户帐号密码错误）的时候，返回错误信息
		if err.Error() == "invalid password" || err.Error() == "failed to get user from database" {
			return &auth.LoginResponse{
				StatusCode: strings.AuthUserLoginFailedCode,
				StatusMsg:  strings.AuthUserLoginFailed,
			}, nil
		}

		// 出现预期外的错误，返回错误信息
		return nil, status.Errorf(codes.Unauthenticated, "login failed: %v", err)
	}
	log.Debug("simple_user: ", simple_user)

	// 返回响应
	loginResponse := &auth.LoginResponse{
		StatusCode:   strings.ServiceOKCode,
		StatusMsg:    strings.ServiceOK,
		Token:        token,
		TokenType:    "Bearer",
		ExpiresAt:    exp_token,
		RefreshToken: refresh_token,
		SimpleUser: &auth.SimpleUser{
			Avatar:    simple_user.Avatar,
			Username:  simple_user.Username,
			Ipaddress: simple_user.Ipaddress,
			Nickname:  simple_user.Nickname,
		},
	}
	return loginResponse, nil
}

// Verify 验证
func (a AuthServiceImpl) Verify(ctx context.Context, request *auth.VerifyRequest) (*auth.VerifyResponse, error) {

	// 调用服务
	claims, err := services.VerifyAccessToken(request.Token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token verification failed: %v", err)
	}

	// 返回响应
	verifyResponse := &auth.VerifyResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		UserId:     claims.UserId,
		Username:   claims.Username,
		ExpiresAt:  claims.ExpiresAt,
	}
	return verifyResponse, nil
}

// Refresh 刷新
func (a AuthServiceImpl) Refresh(ctx context.Context, request *auth.RefreshRequest) (*auth.RefreshResponse, error) {

	// 调用服务
	newToken, exp_token, err := services.RefreshAccessToken(request.RefreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "refresh token failed: %v", err)
	}

	// 返回响应
	refreshResponse := &auth.RefreshResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		Token:      newToken,
		ExpiresAt:  exp_token,
	}
	return refreshResponse, nil
}

// Logout 登出
func (a AuthServiceImpl) Logout(ctx context.Context, request *auth.LogoutRequest) (*auth.LogoutResponse, error) {
	// 解析参数
	accessToken := request.AccessToken
	refreshToken := request.RefreshToken

	// 调用服务
	// 验证访问令牌，确保合法用户的操作
	claim, err := services.VerifyAccessToken(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid access token: %v", err)
	}
	// 撤销刷新令牌
	err = services.RevokeRefreshToken(claim.UserId, refreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "revoke refresh token failed: %v", err)
	}

	// 返回响应
	logoutResponse := &auth.LogoutResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
	}
	return logoutResponse, nil

}

// Register 注册
func (a AuthServiceImpl) Register(ctx context.Context, request *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	log.Infof("Received Register request: %v", request)

	// 创建一个带有超时的新上下文
	// ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	// defer cancel()

	// 解析参数
	user := user.User{
		Username:  request.Username,
		Password:  &request.Password,
		Nickname:  request.Nickname,
		Ipaddress: &request.Ipaddress,
		Gender:    request.Gender,
		Avatar:    conf.DefaultUserAvatar, // 使用默认头像,之后会异步生成头像更换
	}

	// 检查用户是否已存在
	isExist, err := services.CheckUsernameExists(user.Username)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to check username: %v", err)
	}
	if isExist {
		return &auth.RegisterResponse{
			StatusCode: strings.AuthUserExistedCode,
			StatusMsg:  strings.AuthUserExisted,
		}, nil
	}

	// 注册用户,不涉及pictures的操作，头像采用默认头像
	err = services.RegisterUser(&user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "register failed: %v", err)
	}

	// 异步生成头像，创建pictures记录，更新user的avatar字段
	go func() {
		imageURL := "https://api.multiavatar.com/" + request.Username + ".png"
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		smmsToken := os.Getenv("SMMS_TOKEN")
		multiavatarToken := os.Getenv("MULTIAVATAR_KEY")

		avatarUrl, err := services.GenerateAvatar(context.TODO(), imageURL, multiavatarToken, smmsToken, "avatar_"+request.Username+".png", user.Id)
		if err != nil {
			log.Errorf("Failed to generate avatar: %v", err)
			return
		}

		// 更新user的avatar字段
		err = services.UpdateUserAvatar(user.Username, *avatarUrl)
		if err != nil {
			log.Errorf("Failed to update user avatar: %v", err)
		}
	}()

	// 返回响应
	return &auth.RegisterResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
	}, nil
}
