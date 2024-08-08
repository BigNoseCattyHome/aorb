package handlers

import (
	"context"

	"github.com/BigNoseCattyHome/aorb/backend/go-services/user/services"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/user"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/strings"
	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
)

// 使用logging库，添加字段日志 UserRpcServerName
var log = logging.LogService(config.UserRpcServerName)

// UserServiceImpl 实现了 user.UserServiceServer 接口
type UserServiceImpl struct {
	user.UnimplementedUserServiceServer // 嵌入未实现的服务器结构体以保证向前兼容性
}

// New 方法用于初始化 UserServiceImpl，但当前实现为空
func (a UserServiceImpl) New() {
	// 初始化逻辑可以在这里添加
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
