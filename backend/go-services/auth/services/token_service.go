package services

import (
	"errors"
	"os"
	"time"

	"github.com/BigNoseCattyHome/aorb/backend/go-services/auth/models"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
)

// 开发的时候需要在本地环境变量中设置AORB_SECRET_KEY，用于生成JWT令牌
var jwtKey = []byte(os.Getenv("AORB_SECRET_KEY")) // 从环境变量中获取JWT密钥，用于生成JWT令牌

// Claims 结构体，用于存储JWT声明
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateAccessToken 生成JWT令牌
func GenerateAccessToken(user *models.User) (string, error) {
	// 创建声明对象
	expirationTime := time.Now().Add(1 * time.Hour) // 设置令牌过期时间为1小时后
	claims := &Claims{
		UserID:   user.ID,
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
		return "", err
	}

	return tokenString, nil
}

// VerifyAccessToken 验证令牌是否有效，返回声明
func VerifyAccessToken(tokenString string) (*Claims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		log.Error("Failed to parse token: ", err)
		return nil, err
	}

	// 获取声明
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		log.Error("Failed to get claims from token")
		return nil, errors.New("invalid token")
	}

	// 检查令牌是否被撤销
	if CheckTokenRevoked(claims.UserID, tokenString) {
		log.Error("Token has been revoked")
		return nil, errors.New("revoked token")
	}

	return claims, nil
}

// RefreshAccessToken 刷新令牌
func RefreshAccessToken(refreshTokenString string) (string, error) {
	// 验证刷新令牌
	claims, err := VerifyRefreshToken(refreshTokenString)
	if err != nil {
		log.Error("Failed to verify refresh token: ", err)
		return "", err
	}

	// 生成新的访问令牌
	newAccessTokenString, err := GenerateAccessToken(&models.User{ID: claims.UserID})
	if err != nil {
		log.Error("Failed to generate new access token: ", err)
		return "", err
	}

	return newAccessTokenString, nil
}

// VerifyRefreshToken 验证刷新令牌
func VerifyRefreshToken(refreshTokenString string) (*Claims, error) {
	// 解析刷新令牌
	token, err := jwt.ParseWithClaims(refreshTokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		log.Error("Failed to parse refresh token: ", err)
		return nil, err
	}

	// 获取声明
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		log.Error("Failed to get claims from refresh token")
		return nil, errors.New("invalid refresh token")
	}

	// 检查刷新令牌是否被撤销
	if CheckTokenRevoked(claims.UserID, refreshTokenString) {
		log.Error("Refresh token has been revoked")
		return nil, errors.New("revoked refresh token")
	}

	return claims, nil
}

// GenerateRefreshToken 生成新的刷新令牌
func GenerateRefreshToken(user *models.User) (string, error) {
	// 创建声明对象
	expirationTime := time.Now().Add(24 * time.Hour) // 设置刷新令牌过期时间为24小时后
	claims := &Claims{
		UserID:   user.ID,
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
		log.Error("Failed to sign refresh token: ", err)
		return "", err
	}

	return tokenString, nil
}

// CheckTokenRevoked 检查令牌是否被撤销
func CheckTokenRevoked(userID, tokenString string) bool {
	// 这里应该有一个数据库查询或其他逻辑来检查令牌是否被撤销
	// 示例：假设有一个函数 IsTokenRevoked(userID, tokenString) 来检查
	return false
}

// RevokeToken 使令牌失效
func RevokeRefreshToken(userID, tokenString string) error {
	// 这里应该有一个数据库更新或其他逻辑来使令牌失效
	// 示例：假设有一个数据库更新来使令牌失效
	return nil
}
