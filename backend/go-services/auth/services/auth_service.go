package services

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/BigNoseCattyHome/aorb/backend/go-services/api-gateway/models"

	"github.com/dgrijalva/jwt-go"
)

// 从环境变量中获取JWT密钥
var jwtKey = []byte(os.Getenv("AORB_SECRET_KEY"))

// 声明（Claims）结构体
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// 注册用户函数
func 注册用户(user *models.User) error {
	// 使用MD5哈希用户密码
	hasher := md5.New()
	hasher.Write([]byte(user.Password))
	user.Password = hex.EncodeToString(hasher.Sum(nil))

	// 将用户保存到数据库（假设你有一个函数来执行此操作）
	err := 保存用户到数据库(user)
	if err != nil {
		return err
	}

	return nil
}

// 认证用户函数
func 认证用户(user *models.User) (string, error) {
	// 检查用户凭据
	dbUser, err := 从数据库获取用户(user.Username)
	if err != nil {
		return "", err
	}

	// 使用MD5哈希输入的密码并进行比较
	hasher := md5.New()
	hasher.Write([]byte(user.Password))
	if dbUser.Password != hex.EncodeToString(hasher.Sum(nil)) {
		return "", errors.New("用户名或密码无效")
	}

	// 创建JWT令牌
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// 占位：将用户保存到数据库
func 保存用户到数据库(user *models.User) error {
	// 实现此函数以将用户保存到数据库

	return nil
}

// 占位：从数据库获取用户
func 从数据库获取用户(username string) (*models.User, error) {
	// 实现此函数以从数据库获取用户
	return &models.User{
		Username: username,
		Password: "hashed_password", // 替换为数据库中实际的哈希密码
	}, nil
}
