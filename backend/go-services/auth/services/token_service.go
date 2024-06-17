package services

import (
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
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

// 生成JWT令牌
func generateJWTToken(user *models.User) (string, error) {
	// 创建声明对象
	expirationTime := time.Now().Add(1 * time.Hour) // 设置令牌过期时间为1小时后
	claims := &Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			// 在声明中设置过期时间
			ExpiresAt: expirationTime.Unix(),
			// 可以设置其他标准声明，如Issuer, Subject等
			// Issuer: "your_app_name",
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

// 验证令牌是否有效，返回声明
func VerifyJWTToken(tokenString string) (*Claims, error) {
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
		return nil, err
	}

	return claims, nil
}

// 由其他服务调用
// 刷新令牌，调用generateJWTToken函数生成新的令牌
func RefreshJWTToken(tokenString string) (string, error) {
	claims, err := VerifyJWTToken(tokenString)
	if err != nil {
		log.Error("Failed to verify token: ", err)
		return "", err
	}

	// 生成新的JWT令牌
	newTokenString, err := generateJWTToken(&models.User{ID: claims.UserID})
	if err != nil {
		log.Error("Failed to generate new token: ", err)
		return "", err
	}

	return newTokenString, nil
}
