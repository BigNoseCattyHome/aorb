package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/BigNoseCattyHome/aorb/backend/go-services/auth/conf"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/auth"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/user"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/database"
	"google.golang.org/protobuf/types/known/timestamppb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var log = logging.LogService(config.AuthRpcServerName) // 使用logging库，添加日志字段为微服务的名字

// 注册
func RegisterUser(newUser *user.User) error {
	log.Infof("Attempting to register user: %s", newUser.Username)

	// 注册用户的逻辑
	isExistUser, err := getUserByUsername(newUser.Username)
	if err != nil {
		log.Errorf("while checking existing user: %v", err)
		return fmt.Errorf("while checking existing user: %w", err)
	}
	if isExistUser {
		log.Warnf("User already exists: %s", newUser.Username)
		return errors.New("用户名已存在")
	}
	coins := float64(0)
	newUser.Id = primitive.NewObjectID().Hex()
	newUser.Coins = &coins
	newUser.Blacklist = &user.BlackList{}
	newUser.CoinsRecord = &user.CoinRecordList{}
	newUser.Followed = &user.FollowedList{}
	newUser.Follower = &user.FollowerList{}
	newUser.PollAsk = &user.PollAskList{}
	newUser.PollAns = &user.PollAnsList{}
	newUser.PollCollect = &user.PollCollectList{}
	newUser.CreateAt = timestamppb.Now()
	newUser.UpdateAt = timestamppb.Now()
	newUser.DeleteAt = timestamppb.New(time.Time{})
	newUser.BgpicMe = &conf.DefaultUserBgpic
	newUser.BgpicPollcard = &conf.DefaultUserPollcard
	newUser.Bio = &conf.DefaultUserBio

	// 保存用户到数据库
	if err := storeUser(newUser); err != nil {
		log.Errorf("Failed to store user: %v", err)
		return errors.New("注册失败")
	}

	log.Infof("User registered successfully: %s", newUser.Username)

	return nil
}

// 将用户保存到数据库
func storeUser(user *user.User) error {
	collection := database.MongoDbClient.Database("aorb").Collection("users")

	// 将用户信息插入到数据库中
	_, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Error("插入失败", err)
		return err
	}

	return nil
}

// 验证用户密码是否正确，返回 JWT令牌，过期时间，刷新令牌，用户基本信息，错误信息
func AuthenticateUser(user *auth.LoginRequest) (string, int64, string, auth.SimpleUser, error) {
	// 检查用户是否存在
	log.Debug("user: ", user)
	storedUser, err := getUserByID(user.Username)
	log.Info("storedUser: ", storedUser)
	if err != nil {
		log.Error("Failed to get user from database: ", err)
		return "", 0, "", auth.SimpleUser{}, errors.New("failed to get user from database")
	}

	// 检查用户名对应的密码是否正确
	if user.Password != *storedUser.Password {
		log.Error("Invalid password")
		return "", 0, "", auth.SimpleUser{}, errors.New("invalid password")
	}

	// 生成JWT令牌
	tokenString, exp_token, err := GenerateAccessToken(storedUser)
	if err != nil {
		log.Error("Failed to generate JWT token: ", err)
		return "", 0, "", auth.SimpleUser{}, errors.New("failed to generate JWT token")
	}

	// 生成刷新令牌
	fresh_token, err := GenerateRefreshToken(storedUser)
	if err != nil {
		log.Error("Failed to generate refresh token: ", err)
		return "", 0, "", auth.SimpleUser{}, errors.New("failed to generate refresh token")
	}

	// 全部顺利执行，返回用户的基本信息
	simple_user := auth.SimpleUser{
		Username:  storedUser.Username,
		Nickname:  storedUser.Nickname,
		Avatar:    storedUser.Avatar,
		Ipaddress: *storedUser.Ipaddress,
		Gender:    storedUser.Gender,
	}

	return tokenString, exp_token, fresh_token, simple_user, nil
}

// 从数据库获取用户
func getUserByID(userName string) (*user.User, error) {
	res := &user.User{} // 返回指针

	// 使用 ObjectID 进行查询
	result := database.MongoDbClient.Database("aorb").Collection("users").FindOne(context.TODO(), bson.M{"username": userName})

	// 解码结果到 user 结构体
	err := result.Decode(res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("No user found with ID: ", userName)
		} else {
			log.Println("Failed to decode result: ", err)
		}
		return nil, err
	}

	return res, nil
}

// 查询是否username已经存在
// 只有在数据库查询的时候遇到除了mongo.ErrNoDocuments之外的错误才会返回错误
func getUserByUsername(username string) (bool, error) {
	collection := database.MongoDbClient.Database("aorb").Collection("users")

	// 查询用户
	filter := bson.M{"username": username}
	var result bson.M
	err := collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// 没有找到匹配的用户，返回false而不是错误
			return false, nil
		}
		// 其他错误
		log.Fatal(err)
		return false, err
	}

	// 找到匹配的用户
	return true, nil
}

// SmmsResponse represents the response structure from SM.MS API
type SmmsResponse struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		URL string `json:"url"`
	} `json:"data"`
}

// 生成用户头像并上传到 SM.MS 图床
func GenerateAvatar(imageURL, multiavatarToken, smmsToken, fileName string) (*string, error) {
	// 解析原始 URL
	u, err := url.Parse(imageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse image URL: %w", err)
	}

	// 创建查询参数
	query := u.Query()
	query.Set("apikey", multiavatarToken)

	// 将查询参数添加回 URL
	u.RawQuery = query.Encode()

	// 下载头像
	log.Debug("Downloading avatar from: ", u.String())
	avatarResp, err := http.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to downloard the avator: %w", err)
	}
	defer avatarResp.Body.Close()

	// 准备multipart表单
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("smfile", fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	_, err = io.Copy(part, avatarResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}
	writer.Close()

	// 创建上传请求
	req, err := http.NewRequest("POST", "https://sm.ms/api/v2/upload", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", smmsToken)

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var smmsResp SmmsResponse
	err = json.NewDecoder(resp.Body).Decode(&smmsResp)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !smmsResp.Success {
		return nil, fmt.Errorf("failed to upload avatar: %s", smmsResp.Message)
	}

	return &smmsResp.Data.URL, nil
}
