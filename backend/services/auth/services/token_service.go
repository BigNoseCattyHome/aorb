package services

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/BigNoseCattyHome/aorb/backend/services/auth/models"

	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// 开发的时候需要在本地环境变量中设置AORB_SECRET_KEY，用于生成JWT令牌
var jwtKey = []byte(os.Getenv("AORB_SECRET_KEY")) // 从环境变量中获取JWT密钥，用于生成JWT令牌

// Claims 结构体，用于存储JWT声明
type Claims struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateAccessToken 生成JWT令牌
func GenerateAccessToken(user models.User) (string, int64, error) {
	// 创建声明对象
	expirationTime := time.Now().Add(1 * time.Hour) // 设置令牌过期时间为1小时后
	claims := &Claims{
		UserId:   user.Id,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			// 在声明中设置过期时间
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// 创建令牌对象，指定算法和声明
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名令牌
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Error("Failed to sign token: ", err)
		return "", 0, err
	}

	return tokenString, expirationTime.Unix(), nil
}

// VerifyAccessToken 验证令牌是否有效，返回声明
func VerifyAccessToken(tokenString string) (*Claims, error) {
	// 解析令牌，验证签名的合法性
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Println("Invalid token signature")
			return nil, errors.New("invalid token signature")
		}
		log.Error("Failed to parse token: ", err)
		return nil, err
	}

	// 获取声明
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		log.Error("Failed to get claims from token or token is invalid")
		return nil, errors.New("invalid token")
	}

	// 检查令牌是否过期
	if time.Unix(claims.ExpiresAt, 0).Before(time.Now()) {
		log.Error("Token has expired")
		return nil, errors.New("token expired")
	}

	return claims, nil
}

// RefreshAccessToken 刷新令牌
func RefreshAccessToken(refreshTokenString string) (string, int64, error) {
	// 验证刷新令牌
	claims, err := VerifyRefreshToken(refreshTokenString)
	if err != nil {
		log.Error("Failed to verify refresh token: ", err)
		return "", 0, err
	}

	// 生成新的访问令牌
	newAccessTokenString, exp_token, err := GenerateAccessToken(models.User{Id: claims.UserId, Username: claims.Username})
	if err != nil {
		log.Error("Failed to generate new access token: ", err)
		return "", 0, err
	}

	return newAccessTokenString, exp_token, nil
}

// VerifyRefreshToken 验证刷新令牌
func VerifyRefreshToken(refreshTokenString string) (*Claims, error) {
	// 解析刷新令牌
	token, err := jwt.ParseWithClaims(refreshTokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Println("Invalid refresh token signature")
			return nil, errors.New("invalid refresh token signature")
		}
		log.Println("Failed to parse refresh token: ", err)
		return nil, err
	}

	// 获取声明
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		log.Println("Failed to get claims from refresh token or token is invalid")
		return nil, errors.New("invalid refresh token")
	}

	// 检查刷新令牌是否过期
	if time.Unix(claims.ExpiresAt, 0).Before(time.Now()) {
		log.Println("Refresh token has expired")
		return nil, errors.New("refresh token expired")
	}

	// 检查刷新令牌是否被撤销
	if CheckTokenRevoked(claims.UserId, refreshTokenString) {
		log.Println("Refresh token has been revoked")
		return nil, errors.New("revoked refresh token")
	}

	return claims, nil
}

// GenerateRefreshToken 生成新的刷新令牌
func GenerateRefreshToken(user models.User) (string, error) {
	// 创建声明对象
	expirationTime := time.Now().Add(24 * time.Hour) // 设置刷新令牌过期时间为24小时后
	claims := &Claims{
		UserId:   user.Id,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			// 过期时间设置为24小时后
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// 创建令牌对象，指定算法和声明
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名令牌
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Error("Failed to sign refresh token: ", err)
		return "", err
	}

	// 保存到数据库
	client.Database("aorb").Collection("refresh_tokens").InsertOne(context.Background(), bson.M{"user_id": user.Id, "token": tokenString, "revoked": false, "expires_at": expirationTime.Unix()})

	return tokenString, nil
}

// CheckTokenRevoked 检查令牌是否被撤销
func CheckTokenRevoked(userID, tokenString string) bool {
	collection := client.Database("aorb").Collection("refresh_tokens")
	filter := bson.M{"user_id": userID, "token": tokenString}

	// 查找文档
	var result bson.M
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false // 令牌不存在，视为未撤销
		}
		log.Fatal(err)
	}

	revoked, ok := result["revoked"].(bool)
	if !ok {
		return false // 如果 revoked 字段不存在或类型不匹配，视为未撤销
	}

	return revoked
}

// RevokeToken 使令牌失效
func RevokeRefreshToken(userID, tokenString string) error {
	collection := client.Database("aorb").Collection("refresh_tokens")

	// filter表示要更新的文档，update表示要更新的字段
	filter := bson.M{"user_id": userID, "token": tokenString}
	update := bson.M{"$set": bson.M{"revoked": true}}

	// 更新文档
	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Print("Failed to revoke token: ", err)
		return err
	}

	// 没有找到要更新的文档
	if result.ModifiedCount == 0 {
		return errors.New("token not found")
	}

	return nil
}
