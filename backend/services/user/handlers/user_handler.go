package handlers

import (
	"context"
	"fmt"

	messagePb "github.com/BigNoseCattyHome/aorb/backend/rpc/message"
	"github.com/BigNoseCattyHome/aorb/backend/utils/extra/tracing"
	grpc2 "github.com/BigNoseCattyHome/aorb/backend/utils/grpc"
	"github.com/BigNoseCattyHome/aorb/backend/utils/rabbitmq"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/database"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/BigNoseCattyHome/aorb/backend/rpc/user"
	"github.com/BigNoseCattyHome/aorb/backend/services/user/services"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/strings"
	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
)

// 使用logging库，添加字段日志 UserRpcServerName
var log = logging.LogService(config.UserRpcServerName)
var messageClient messagePb.MessageServiceClient

// UserServiceImpl 实现了 user.UserServiceServer 接口
type UserServiceImpl struct {
	user.UnimplementedUserServiceServer // 嵌入未实现的服务器结构体以保证向前兼容性
}

var conn *amqp.Connection

var channel *amqp.Channel

func exitOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// New 方法用于初始化 UserServiceImpl，但当前实现为空
func (a UserServiceImpl) New() {
	// 初始化逻辑可以在这里添加
	messageRpcConn := grpc2.Connect(config.MessageRpcServerName)
	messageClient = messagePb.NewMessageServiceClient(messageRpcConn)

	var err error
	conn, err = amqp.Dial(rabbitmq.BuildMqConnAddr())
	exitOnError(err)

	channel, err = conn.Channel()
	exitOnError(err)

	err = channel.ExchangeDeclare(
		strings.EventExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	exitOnError(err)
}

func CloseMQConn() {
	if err := conn.Close(); err != nil {
		panic(err)
	}

	if err := channel.Close(); err != nil {
		panic(err)
	}
}

// GetUserInfo 方法用于获取用户信息
func (a UserServiceImpl) GetUserInfo(ctx context.Context, request *user.UserRequest) (resp *user.UserResponse, err error) {
	// TODO 从缓存中获取用户信息
	/*
		ok, err := cached.ScanGet(ctx, "UserInfo", &userModel, "users")
		if err != nil {
			resp = &user.UserResponse{
				StatusCode: strings.UserServiceInnerErrorCode,
				StatusMsg:  strings.UserServiceInnerError,
			}
			return resp, err
		}
	*/

	username := request.GetUsername() // 获取请求中的 username
	fields := request.GetFields()     // 获取请求中的 fields
	log.Debug("GetUserInfo: ", username, fields)

	res, err := services.GetUserInfo(username, fields) // 调用服务层的 GetUserInfo 方法
	if err != nil {
		log.Error("GetUserInfo error: ", err)
		return nil, err
	}
	return res, nil
}

// CheckUserExists 方法用于检查用户是否存在
func (a UserServiceImpl) CheckUserExists(ctx context.Context, request *user.UserExistRequest) (resp *user.UserExistResponse, err error) {
	username := request.GetUsername()
	log.Debug("CheckUserExists: ", username)
	resp, err = services.CheckUserExists(username)
	if err != nil {
		log.Error("CheckUserExists error: ", err)
		return nil, err
	}
	return resp, nil

}

// IsUserFollowing 方法用于检查用户是否关注了另一个用户
func (a UserServiceImpl) IsUserFollowing(ctx context.Context, request *user.IsUserFollowingRequest) (resp *user.IsUserFollowingResponse, err error) {
	username := request.GetUsername()
	taeget_username := request.GetTargetUsername()
	log.Debug("Searching for: ", username, " following ", taeget_username)

	resp, err = services.IsUserFollowing(username, taeget_username)
	if err != nil {
		log.Error("IsUserFollowing error: ", err)
		return nil, err
	}

	return resp, nil
}

// UpdateUser 方法用于更新用户信息
func (a UserServiceImpl) UpdateUser(ctx context.Context, request *user.UpdateUserRequest) (resp *user.UpdateUserResponse, err error) {
	userId := request.GetUserId()
	log.Debug("Updating user: ", userId)

	// 获取需要更新的字段
	updateFields := make(map[string]interface{})
	if request.Username != nil {
		updateFields["username"] = *request.Username
	}
	if request.Nickname != nil {
		updateFields["nickname"] = *request.Nickname
	}
	if request.Bio != nil {
		updateFields["bio"] = *request.Bio
	}
	if request.Gender != nil {
		updateFields["gender"] = *request.Gender
	}
	if request.BgpicMe != nil {
		updateFields["bgpic_me"] = request.BgpicMe
	}
	if request.BgpicPollcard != nil {
		updateFields["bgpic_pollcard"] = request.BgpicPollcard
	}
	if request.Avatar != nil {
		updateFields["avatar"] = request.Avatar
	}

	// 调用服务层的 UpdateUser 方法
	resp, err = services.UpdateUserInService(ctx, userId, updateFields)
	if err != nil {
		log.Error("Failed to update user: ", err)

		// 如果用户已经存在，返回错误信息
		if err.Error() == "username already exists" {
			return &user.UpdateUserResponse{
				StatusCode: strings.AuthUserExistedCode,
				StatusMsg:  strings.AuthUserExisted,
			}, nil
		}
		return nil, err
	}

	return resp, nil
}

func (a UserServiceImpl) FollowUser(ctx context.Context, request *user.FollowUserRequest) (response *user.FollowUserResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "UserService.FollowUser")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("UserService.FollowUser").WithContext(ctx)
	logger.WithFields(logrus.Fields{
		"username":        request.Username,
		"target_username": request.TargetUsername,
	}).Debugf("Process start")

	// 不能关注自己
	if request.Username == request.TargetUsername {
		logger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Errorf("关注失败, User can't follow self")
		logging.SetSpanError(span, err)
		response = &user.FollowUserResponse{
			StatusCode: strings.UnableToFollowErrorCode,
			StatusMsg:  strings.UnableToFollowError,
		}
		return
	}

	// 先查用户是否存在
	checkUserExists, err := a.CheckUserExists(ctx, &user.UserExistRequest{
		Username: request.Username,
	})
	if err != nil {
		logger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Errorf("关注失败, Error when calling rpc CheckUserExists about username: %s", request.Username)
		logging.SetSpanError(span, err)
		response = &user.FollowUserResponse{
			StatusCode: strings.UnableToFollowErrorCode,
			StatusMsg:  strings.UnableToFollowError,
		}
		return
	}
	if !checkUserExists.Existed {
		logger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Errorf("关注失败, username %s does not exist", request.Username)
		logging.SetSpanError(span, err)
		response = &user.FollowUserResponse{
			StatusCode: strings.UnableToFollowErrorCode,
			StatusMsg:  strings.UnableToFollowError,
		}
		return
	}

	checkTargetUserExists, err := a.CheckUserExists(ctx, &user.UserExistRequest{
		Username: request.TargetUsername,
	})
	if err != nil {
		logger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Errorf("关注失败, Error when calling rpc CheckUserExists about username: %s", request.TargetUsername)
		logging.SetSpanError(span, err)
		response = &user.FollowUserResponse{
			StatusCode: strings.UnableToFollowErrorCode,
			StatusMsg:  strings.UnableToFollowError,
		}
		return
	}
	if !checkTargetUserExists.Existed {
		logger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Errorf("关注失败, username %s does not exist", request.TargetUsername)
		logging.SetSpanError(span, err)
		response = &user.FollowUserResponse{
			StatusCode: strings.UnableToFollowErrorCode,
			StatusMsg:  strings.UnableToFollowError,
		}
		return
	}

	userCollection := database.MongoDbClient.Database("aorb").Collection("users")

	// 查看是否已经关注
	filter4Check := bson.D{{"username", request.Username}}
	cursor := userCollection.FindOne(ctx, filter4Check)
	var result bson.M
	if err = cursor.Decode(&result); err != nil {
		logger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Errorf("关注失败, Error when checking username %s's following status from database", request.Username)
		logging.SetSpanError(span, err)
		response = &user.FollowUserResponse{
			StatusCode: strings.UnableToFollowErrorCode,
			StatusMsg:  strings.UnableToFollowError,
		}
		return
	}
	followedList := result["followed"].(bson.M)["usernames"].(bson.A)
	for _, followedUsername := range followedList {
		if request.TargetUsername == followedUsername {
			logger.WithFields(logrus.Fields{
				"username": request.Username,
			}).Errorf("关注失败, username %s has already followed target_username %s", request.Username, request.TargetUsername)
			logging.SetSpanError(span, err)
			response = &user.FollowUserResponse{
				StatusCode: strings.UnableToFollowErrorCode,
				StatusMsg:  strings.UnableToFollowError,
			}
			return
		}
	}

	// 更新username的followed_list
	filter4Follow := bson.D{{"username", request.Username}}
	update4Follow := bson.D{{"$push", bson.D{{"followed.usernames", request.TargetUsername}}}}
	_, err = userCollection.UpdateOne(ctx, filter4Follow, update4Follow)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Errorf("关注失败, Error when updating username %s's follow user", request.Username)
		logging.SetSpanError(span, err)
		response = &user.FollowUserResponse{
			StatusCode: strings.UnableToFollowErrorCode,
			StatusMsg:  strings.UnableToFollowError,
		}
		return
	}

	// 更新targetUsername的follower_list
	filter4Followed := bson.D{{"username", request.TargetUsername}}
	update4Followed := bson.D{{"$push", bson.D{{"follower.usernames", request.Username}}}}
	_, err = userCollection.UpdateOne(ctx, filter4Followed, update4Followed)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Errorf("关注失败, Error when updating target_username %s's followed user", request.TargetUsername)
		logging.SetSpanError(span, err)
		response = &user.FollowUserResponse{
			StatusCode: strings.UnableToFollowErrorCode,
			StatusMsg:  strings.UnableToFollowError,
		}
		return
	}

	// 向targetUsername发送message
	messageActionResponse, err := messageClient.MessageAction(ctx, &messagePb.MessageActionRequest{
		FromUsername: request.Username,
		ToUsername:   request.TargetUsername,
		ActionType:   messagePb.ActionMessageType_ACTION_MESSAGE_TYPE_ADD,
		MessageType:  messagePb.MessageType_MESSAGE_TYPE_FOLLOW,
		Action: &messagePb.MessageActionRequest_MessageContent{
			MessageContent: fmt.Sprintf("收到了来自%s的关注", request.Username),
		},
	})
	if err != nil && messageActionResponse.StatusCode != 0 {
		logger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Errorf("关注失败, Error when calling rpc MessageAction")
		logging.SetSpanError(span, err)
		response = &user.FollowUserResponse{
			StatusCode: strings.UnableToFollowErrorCode,
			StatusMsg:  strings.UnableToFollowError,
		}
		return
	}

	response = &user.FollowUserResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
	}
	return
}

func (a UserServiceImpl) UnfollowUser(ctx context.Context, request *user.FollowUserRequest) (response *user.FollowUserResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "UserService.FollowUser")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("UserService.FollowUser").WithContext(ctx)
	logger.WithFields(logrus.Fields{
		"username":        request.Username,
		"target_username": request.TargetUsername,
	}).Debugf("Process start")

	// 不能取关自己
	if request.Username == request.TargetUsername {
		logger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Errorf("取关失败, User can't unfollow self")
		logging.SetSpanError(span, err)
		response = &user.FollowUserResponse{
			StatusCode: strings.UnableToUnFollowErrorCode,
			StatusMsg:  strings.UnableToUnFollowError,
		}
		return
	}

	// 先查用户是否存在
	checkUserExists, err := a.CheckUserExists(ctx, &user.UserExistRequest{
		Username: request.Username,
	})
	if err != nil {
		logger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Errorf("取关失败, Error when calling rpc CheckUserExists about username: %s", request.Username)
		logging.SetSpanError(span, err)
		response = &user.FollowUserResponse{
			StatusCode: strings.UnableToUnFollowErrorCode,
			StatusMsg:  strings.UnableToUnFollowError,
		}
		return
	}
	if !checkUserExists.Existed {
		logger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Errorf("取关失败, username %s does not exist", request.Username)
		logging.SetSpanError(span, err)
		response = &user.FollowUserResponse{
			StatusCode: strings.UnableToUnFollowErrorCode,
			StatusMsg:  strings.UnableToUnFollowError,
		}
		return
	}

	checkTargetUserExists, err := a.CheckUserExists(ctx, &user.UserExistRequest{
		Username: request.TargetUsername,
	})
	if err != nil {
		logger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Errorf("取关失败, Error when calling rpc CheckUserExists about username: %s", request.TargetUsername)
		logging.SetSpanError(span, err)
		response = &user.FollowUserResponse{
			StatusCode: strings.UnableToUnFollowErrorCode,
			StatusMsg:  strings.UnableToUnFollowError,
		}
		return
	}
	if !checkTargetUserExists.Existed {
		logger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Errorf("取关失败, username %s does not exist", request.TargetUsername)
		logging.SetSpanError(span, err)
		response = &user.FollowUserResponse{
			StatusCode: strings.UnableToUnFollowErrorCode,
			StatusMsg:  strings.UnableToUnFollowError,
		}
		return
	}

	userCollection := database.MongoDbClient.Database("aorb").Collection("users")

	// 查看是否已经关注
	filter4Check := bson.D{{"username", request.Username}}
	cursor := userCollection.FindOne(ctx, filter4Check)
	var result bson.M
	if err = cursor.Decode(&result); err != nil {
		logger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Errorf("关注失败, Error when checking username %s's following status from database", request.Username)
		logging.SetSpanError(span, err)
		response = &user.FollowUserResponse{
			StatusCode: strings.UnableToUnFollowErrorCode,
			StatusMsg:  strings.UnableToUnFollowError,
		}
		return
	}
	followedList := result["followed"].(bson.M)["usernames"].(bson.A)
	flag := false
	for _, followedUsername := range followedList {
		if request.TargetUsername == followedUsername {
			flag = true
			break
		}
	}

	if flag == false {
		logger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Errorf("取关失败, username %s hasn't followed target_username %s", request.Username, request.TargetUsername)
		logging.SetSpanError(span, err)
		response = &user.FollowUserResponse{
			StatusCode: strings.UnableToUnFollowErrorCode,
			StatusMsg:  strings.UnableToUnFollowError,
		}
		return
	}

	// 更新username的followed_list
	filter4Follow := bson.D{{"username", request.Username}}
	update4Follow := bson.D{{"$pull", bson.D{{"followed.usernames", request.TargetUsername}}}}
	_, err = userCollection.UpdateOne(ctx, filter4Follow, update4Follow)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Errorf("取关失败, Error when updating username %s's follow user", request.Username)
		logging.SetSpanError(span, err)
		response = &user.FollowUserResponse{
			StatusCode: strings.UnableToUnFollowErrorCode,
			StatusMsg:  strings.UnableToUnFollowError,
		}
		return
	}

	// 更新targetUsername的follower_list
	filter4Followed := bson.D{{"username", request.TargetUsername}}
	update4Followed := bson.D{{"$pull", bson.D{{"follower.usernames", request.Username}}}}
	_, err = userCollection.UpdateOne(ctx, filter4Followed, update4Followed)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Errorf("取关失败, Error when updating target_username %s's followed user", request.TargetUsername)
		logging.SetSpanError(span, err)
		response = &user.FollowUserResponse{
			StatusCode: strings.UnableToUnFollowErrorCode,
			StatusMsg:  strings.UnableToUnFollowError,
		}
		return
	}

	response = &user.FollowUserResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
	}
	return
}
