package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/BigNoseCattyHome/aorb/backend/utils/storage/database"

	"github.com/BigNoseCattyHome/aorb/backend/go-services/auth/models"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var log = logging.LogService(config.AuthRpcServerName) // 使用logging库，添加日志字段为微服务的名字

// 注册
func RegisterUser(newuser models.User) error {
	log.Infof("Attempting to register user: %s", newuser.Username)

	// 注册用户的逻辑
	isExistUser, err := getUserByUsername(newuser.Username)
	if err != nil {
		log.Errorf("while checking existing user: %v", err)
		return fmt.Errorf("while checking existing user: %w", err)
	}
	if isExistUser {
		log.Warnf("User already exists: %s", newuser.Username)
		return errors.New("用户名已存在")
	}

	// 生成新的ObjectID
	newuser.ID = primitive.NewObjectID().Hex()

	// 初始化其他字段
	newuser.Coins = 0
	newuser.Blacklist = []string{}
	newuser.CoinsRecord = []models.CoinRecord{}
	newuser.Followed = []string{}
	newuser.Follower = []string{}
	newuser.QuestionsAsk = []string{}
	newuser.QuestionsAsw = []string{}
	newuser.QuestionsCollect = []string{}
	newuser.CreateAt = time.Now()
	newuser.UpdateAt = time.Now()
	newuser.DeleteAt = time.Time{}

	// 保存用户到数据库
	if err := storeUser(newuser); err != nil {
		log.Errorf("Failed to store user: %v", err)
		return errors.New("注册失败")
	}

	log.Infof("User registered successfully: %s", newuser.Username)

	return nil
}

// 将用户保存到数据库
func storeUser(user models.User) error {
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
func AuthenticateUser(user models.LoginRequest) (string, int64, string, models.SimpleUser, error) {
	// 检查用户是否存在
	storedUser, err := getUserbyID(user.Username)
	if err != nil {
		log.Error("Failed to get user from database: ", err)
		return "", 0, "", models.SimpleUser{}, errors.New("failed to get user from database")
	}

	// 检查用户名对应的密码是否正确
	if user.Password != storedUser.Password {
		log.Error("Invalid password")
		return "", 0, "", models.SimpleUser{}, errors.New("invalid password")
	}

	// 生成JWT令牌
	tokenString, exp_token, err := GenerateAccessToken(storedUser)
	if err != nil {
		log.Error("Failed to generate JWT token: ", err)
		return "", 0, "", models.SimpleUser{}, errors.New("failed to generate JWT token")
	}

	// 生成刷新令牌
	fresh_token, err := GenerateRefreshToken(storedUser)
	if err != nil {
		log.Error("Failed to generate refresh token: ", err)
		return "", 0, "", models.SimpleUser{}, errors.New("failed to generate refresh token")
	}

	// 全部顺利执行，返回用户的基本信息
	simple_user := models.SimpleUser{
		ID:        storedUser.ID,
		Nickname:  storedUser.Nickname,
		Avatar:    storedUser.Avatar,
		Ipaddress: storedUser.Ipaddress,
	}

	return tokenString, exp_token, fresh_token, simple_user, nil
}

// 从数据库获取用户
func getUserbyID(id string) (models.User, error) {
	user := models.User{}

	// 将字符串转换为 ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Failed to convert string to ObjectID: ", err)
		return models.User{}, err
	}

	// 使用 ObjectID 进行查询
	result := database.MongoDbClient.Database("aorb").Collection("users").FindOne(context.TODO(), bson.M{"_id": objectID})

	// 解码结果到 user 结构体
	err = result.Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("No user found with ID: ", id)
		} else {
			log.Println("Failed to decode result: ", err)
		}
		return models.User{}, err
	}

	return user, nil
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
