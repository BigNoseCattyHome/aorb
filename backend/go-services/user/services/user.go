package main

import (
	"context"
	userModels "github.com/BigNoseCattyHome/aorb/backend/go-services/user/models"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/user"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/strings"
	"github.com/BigNoseCattyHome/aorb/backend/utils/extra/tracing"
	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/cached"
	"github.com/sirupsen/logrus"
)

type UserServiceImpl struct {
	user.UserServiceServer
}

func (a UserServiceImpl) New() {

}

func (a UserServiceImpl) GetUserInfo(ctx context.Context, request *user.UserRequest) (resp *user.UserResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "GetUserInfo")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("UserService.GetUserInfo").WithContext(ctx)

	var userModel userModels.User
	userModel.ID = request.UserId
	ok, err := cached.ScanGet(ctx, "UserInfo", &userModel, "users")
	if err != nil {
		resp = &user.UserResponse{
			StatusCode: strings.UserServiceInnerErrorCode,
			StatusMsg:  strings.UserServiceInnerError,
		}
		return
	}

	if !ok {
		resp = &user.UserResponse{
			StatusCode: strings.UserNotExistedCode,
			StatusMsg:  strings.UserNotExisted,
			User:       nil,
		}
		logger.WithFields(logrus.Fields{
			"user": request.UserId,
		}).Infof("Do not exist")
		return
	}

	resp = &user.UserResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		User: &user.User{
			Avatar:           userModel.Avatar,
			Blacklist:        nil,
			Coins:            nil,
			CoinsRecord:      nil,
			Followed:         nil,
			Follower:         nil,
			Id:               request.UserId,
			Ipaddress:        nil,
			Nickname:         userModel.Nickname,
			QuestionsAsk:     nil,
			QuestionsAsw:     nil,
			QuestionsCollect: nil,
			Username:         userModel.Username,
		},
	}

	// 这里后期relation功能做好之后需要添加相应内容

	return
}

func (a UserServiceImpl) GetUserExistInformation(ctx context.Context, request *user.UserExistRequest) (resp *user.UserExistResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "GetUserExisted")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("UserService.GetUserExisted").WithContext(ctx)

	var userModel userModels.User
	userModel.ID = request.UserId
	ok, err := cached.ScanGet(ctx, "UserExisted", &userModel, "users")

	if err != nil {
		logger.WithFields(logrus.Fields{
			"err": err,
		}).Errorf("Error when selecting user info")
		logging.SetSpanError(span, err)
		resp = &user.UserExistResponse{
			StatusCode: strings.UserServiceInnerErrorCode,
			StatusMsg:  strings.UserServiceInnerError,
			Existed:    false,
		}
		return
	}

	if !ok {
		resp = &user.UserExistResponse{
			StatusCode: strings.ServiceOKCode,
			StatusMsg:  strings.ServiceOK,
			Existed:    false,
		}
		logger.WithFields(logrus.Fields{
			"user": request.UserId,
		}).Infof("User do not exist")
		return
	}

	resp = &user.UserExistResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		Existed:    true,
	}
	return
}
