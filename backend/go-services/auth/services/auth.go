package main

import (
	"context"
	"fmt"
	userModels "github.com/BigNoseCattyHome/aorb/backend/go-services/user/models"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/auth"
	user2 "github.com/BigNoseCattyHome/aorb/backend/rpc/user"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/strings"
	"github.com/BigNoseCattyHome/aorb/backend/utils/extra/tracing"
	grpc2 "github.com/BigNoseCattyHome/aorb/backend/utils/grpc"
	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/cached"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/database"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/redis"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/willf/bloom"
	"go.mongodb.org/mongo-driver/bson"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

var userClient user2.UserServiceClient
var BloomFilter *bloom.BloomFilter

type AuthServiceImpl struct {
	auth.AuthServiceServer
}

func (a AuthServiceImpl) New() {
	userRpcConn := grpc2.Connect(config.UserRpcServerName)
	userClient = user2.NewUserServiceClient(userRpcConn)
}

func (a AuthServiceImpl) Authenticate(ctx context.Context, request *auth.AuthenticateRequest) (resp *auth.AuthenticateResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "AuthenticateService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("AuthService.Authenticate").WithContext(ctx)

	userId, ok, err := hasToken(ctx, request.Token)

	if err != nil {
		resp = &auth.AuthenticateResponse{
			StatusCode: strings.AuthServiceInnerErrorCode,
			StatusMsg:  strings.AuthServiceInnerError,
		}
		return
	}

	if !ok {
		resp = &auth.AuthenticateResponse{
			StatusCode: strings.UserNotExistedCode,
			StatusMsg:  strings.UserNotExisted,
		}
		return
	}

	id, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err":   err,
			"token": request.Token,
		}).Warnf("AuthService Authenticate Action failed to response when parsering uint")
		logging.SetSpanError(span, err)

		resp = &auth.AuthenticateResponse{
			StatusCode: strings.AuthServiceInnerErrorCode,
			StatusMsg:  strings.AuthServiceInnerError,
		}
		return
	}

	resp = &auth.AuthenticateResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		UserId:     uint32(id),
	}
	return
}

func (a AuthServiceImpl) Register(ctx context.Context, request *auth.RegisterRequest) (resp *auth.RegisterResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "RegisterService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("AuthService.Register").WithContext(ctx)

	resp = &auth.RegisterResponse{}
	var user userModels.User
	collection := database.MongoDbClient.Database("aorb").Collection("users")
	cursor := collection.FindOne(ctx, bson.M{"username": request.Username})
	result := cursor.Decode(&user)
	if result != nil {
		resp = &auth.RegisterResponse{
			StatusCode: strings.AuthUserExistedCode,
			StatusMsg:  strings.AuthUserExisted,
		}
		return
	}

	var hashedPassword string
	if hashedPassword, err = hashPassword(ctx, request.Password); err != nil {
		logger.WithFields(logrus.Fields{
			"err":      result.Error,
			"username": request.Username,
		}).Warnf("AuthService Register Action failed to response when hashing password")
		logging.SetSpanError(span, err)

		resp = &auth.RegisterResponse{
			StatusCode: strings.AuthServiceInnerErrorCode,
			StatusMsg:  strings.AuthServiceInnerError,
		}
		return
	}

	user.Password = hashedPassword
	user.CreateAt = time.Now()
	user.UpdateAt = time.Now()
	user.Username = request.Username
	user.Avatar = ""

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err":      result.Error,
			"username": request.Username,
		}).Warnf("AuthService Register Action failed to response when creating user")
		logging.SetSpanError(span, err)

		resp = &auth.RegisterResponse{
			StatusCode: strings.AuthServiceInnerErrorCode,
			StatusMsg:  strings.AuthServiceInnerError,
		}
		return
	}

	resp.Token, err = getToken(ctx, user.ID)

	if err != nil {
		resp = &auth.RegisterResponse{
			StatusCode: strings.AuthServiceInnerErrorCode,
			StatusMsg:  strings.AuthServiceInnerError,
		}
		return
	}

	logger.WithFields(logrus.Fields{
		"username": request.Username,
	}).Infof("User register success!")

	resp.UserId = user.ID
	resp.StatusCode = strings.ServiceOKCode
	resp.StatusMsg = strings.ServiceOK

	BloomFilter.AddString(user.Username)
	logger.WithFields(logrus.Fields{
		"username": user.Username,
	}).Infof("Publishing user name to redis channel")
	err = redis.Client.Publish(ctx, config.BloomRedisChannel, user.Username).Err()
	if err != nil {
		logger.WithFields(logrus.Fields{
			"err":      err,
			"username": user.Username,
		}).Errorf("Publishing user name to redis channel happens error")
		logging.SetSpanError(span, err)
	}

	return
}

func (a AuthServiceImpl) Login(ctx context.Context, request *auth.LoginRequest) (resp *auth.LoginResponse, err error) {
	ctx, span := tracing.Tracer.Start(ctx, "LoginService")
	defer span.End()
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("AuthService.Login").WithContext(ctx)
	logger.WithFields(logrus.Fields{
		"username": request.Username,
	}).Debugf("User try to log in.")

	if !BloomFilter.TestString(request.Username) {
		resp = &auth.LoginResponse{
			StatusCode: strings.UnableToQueryUserErrorCode,
			StatusMsg:  strings.UnableToQueryUserError,
		}

		logger.WithFields(logrus.Fields{
			"username": request.Username,
		}).Infof("The user is blocked by Bloom Filter")
		return
	}

	resp = &auth.LoginResponse{}
	user := userModels.User{
		Username: request.Username,
	}

	ok, err := isUserVerifiedInRedis(ctx, request.Username, request.Password)
	if err != nil {
		resp = &auth.LoginResponse{
			StatusCode: strings.AuthServiceInnerErrorCode,
			StatusMsg:  strings.AuthServiceInnerError,
		}
		logging.SetSpanError(span, err)
		return
	}

	if !ok {
		collection := database.MongoDbClient.Database("aorb").Collection("users")
		singleResult := collection.FindOne(ctx, bson.M{"username": request.Username})
		var result userModels.User
		err = singleResult.Decode(&result)

		if err != nil {
			logger.WithFields(logrus.Fields{
				"err":      err,
				"username": request.Username,
			}).Warnf("AuthService Login Action failed to response with inner err.")
			logging.SetSpanError(span, err)

			resp = &auth.LoginResponse{
				StatusCode: strings.AuthServiceInnerErrorCode,
				StatusMsg:  strings.AuthServiceInnerError,
			}
			logging.SetSpanError(span, err)
			return
		}

		if singleResult == nil {
			resp = &auth.LoginResponse{
				StatusCode: strings.UserNotExistedCode,
				StatusMsg:  strings.UserNotExisted,
			}
			return
		}

		if !checkPasswordHash(ctx, request.Password, user.Password) {
			resp = &auth.LoginResponse{
				StatusCode: strings.AuthUserLoginFailedCode,
				StatusMsg:  strings.AuthUserLoginFailed,
			}
			return
		}

		hashed, errs := hashPassword(ctx, request.Password)
		if errs != nil {
			logger.WithFields(logrus.Fields{
				"err":      errs,
				"username": request.Username,
			}).Warnf("AuthService Login Action failed to response with inner err.")
			logging.SetSpanError(span, errs)

			resp = &auth.LoginResponse{
				StatusCode: strings.AuthServiceInnerErrorCode,
				StatusMsg:  strings.AuthServiceInnerError,
			}
			logging.SetSpanError(span, err)
			return
		}

		if err = setUserInfoToRedis(ctx, user.Username, hashed); err != nil {
			resp = &auth.LoginResponse{
				StatusCode: strings.AuthServiceInnerErrorCode,
				StatusMsg:  strings.AuthServiceInnerError,
			}
			logging.SetSpanError(span, err)
			return
		}
		cached.Write(ctx, fmt.Sprintf("UserId%s", request.Username), strconv.Itoa(int(user.ID)), true)
	} else {
		id, _, err := cached.Get(ctx, fmt.Sprintf("UserId%s", request.Username))
		if err != nil {
			resp = &auth.LoginResponse{
				StatusCode: strings.AuthServiceInnerErrorCode,
				StatusMsg:  strings.AuthServiceInnerError,
			}
			logging.SetSpanError(span, err)
			return nil, err
		}
		uintId, _ := strconv.ParseUint(id, 10, 32)
		user.ID = uint32(uintId)
	}

	token, err := getToken(ctx, user.ID)
	if err != nil {
		resp = &auth.LoginResponse{
			StatusCode: strings.AuthServiceInnerErrorCode,
			StatusMsg:  strings.AuthServiceInnerError,
		}
		return
	}

	logger.WithFields(logrus.Fields{
		"token":  token,
		"userId": user.ID,
	}).Debugf("User log in sucess !")
	resp = &auth.LoginResponse{
		StatusCode: strings.ServiceOKCode,
		StatusMsg:  strings.ServiceOK,
		UserId:     user.ID,
		Token:      token,
	}
	return
}

func setUserInfoToRedis(ctx context.Context, username string, password string) error {
	_, ok, err := cached.Get(ctx, "UserLog"+username)
	if err != nil {
		return err
	}
	if ok {
		cached.TagDelete(ctx, "UserLog"+username)
	}
	cached.Write(ctx, "UserLog"+username, password, true)
	return nil
}

func isUserVerifiedInRedis(ctx context.Context, username string, password string) (bool, error) {
	pass, ok, err := cached.Get(ctx, "UserLog"+username)
	if err != nil {
		return false, err
	}

	if !ok {
		return false, nil
	}

	if checkPasswordHash(ctx, password, pass) {
		return true, nil
	}

	return false, nil
}

func checkPasswordHash(ctx context.Context, password string, hash string) bool {
	_, span := tracing.Tracer.Start(ctx, "PasswordHashChecked")
	defer span.End()
	logging.SetSpanWithHostname(span)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func getToken(ctx context.Context, userId uint32) (string, error) {
	span := trace.SpanFromContext(ctx)
	logging.SetSpanWithHostname(span)
	logger := logging.LogService("AuthService.Login").WithContext(ctx)
	logger.WithFields(logrus.Fields{
		"userId": userId,
	}).Debugf("Select for user token")
	return cached.GetWithFunc(ctx, "U2T"+strconv.FormatUint(uint64(userId), 10),
		func(ctx context.Context, key string) (string, error) {
			span := trace.SpanFromContext(ctx)
			token := uuid.New().String()
			span.SetAttributes(attribute.String("token", token))
			cached.Write(ctx, "T2U"+token, strconv.FormatUint(uint64(userId), 10), true)
			return token, nil
		})
}

func hashPassword(ctx context.Context, password string) (string, error) {
	_, span := tracing.Tracer.Start(ctx, "PasswordHash")
	defer span.End()
	logging.SetSpanWithHostname(span)
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func hasToken(ctx context.Context, token string) (string, bool, error) {
	return cached.Get(ctx, "T2U"+token)
}
