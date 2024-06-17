package services

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/BigNoseCattyHome/aorb/backend/go-services/auth/conf"
	"github.com/BigNoseCattyHome/aorb/backend/go-services/auth/models"
	log "github.com/sirupsen/logrus"
)

var client *mongo.Client // MongoDB客户端

// init函数在main函数之前执行
func init() {
	// 初始化AppConfig
	if err := conf.LoadConfig(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 设置客户端选项，从AppConfig中获取MongoDB的URL
	clientOptions := options.Client().ApplyURI(conf.AppConfig.GetString("db.mongodb_url"))

	// 连接到MongoDB
	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
}

// 注册
func RegisterUser(user *models.User) error {
	// 在这里写注册用户的逻辑

	return nil
}

// 将用户保存到数据库
func storeUser(user *models.User) error {
	// 实现此函数以将用户保存到数据库

	return nil
}

// 验证用户（登录）
// 返回(成功): JWT令牌, 用户信息, nil
// 返回(失败): "", 空的SimpleUser, 错误信息
func AuthenticateUser(user *models.RequestLogin) (string, models.SimpleUser, error) {
	// 检查用户是否存在
	storedUser, err := getUser(user.ID)
	if err != nil {
		log.Error("Failed to get user from database: ", err)
		return "", models.SimpleUser{}, errors.New("failed to get user from database")
	}

	// 检查用户名对应的密码是否正确
	if user.Password != storedUser.Password {
		log.Error("Invalid password")
		return "", models.SimpleUser{}, errors.New("invalid password")
	}

	// 生成JWT令牌
	tokenString, err := generateJWTToken(storedUser)
	if err != nil {
		log.Error("Failed to generate JWT token: ", err)
		return "", models.SimpleUser{}, errors.New("failed to generate JWT token")
	}

	// 全部顺利执行，返回用户的基本信息
	simple_user := models.SimpleUser{
		ID:        storedUser.ID,
		Nickname:  storedUser.Nickname,
		Avatar:    storedUser.Avatar,
		Ipaddress: storedUser.Ipaddress,
	}

	return tokenString, simple_user, nil
}

// 从数据库获取用户
func getUser(id string) (*models.User, error) {
	user := models.User{}

	// 将字符串转换为 ObjectId
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Error("Failed to convert string to ObjectID: ", err)
		return &models.User{}, err
	}

	// 使用 ObjectId 进行查询
	result := client.Database("aorb").Collection("users").FindOne(context.TODO(), bson.M{"_id": objectID})

	// 解码结果到 user 结构体
	err = result.Decode(&user)
	if err != nil {
		log.Error("Failed to decode result: ", err)
		return &models.User{}, err
	}

	return &user, nil
}
