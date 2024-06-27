package controllers

import (
	"net/http"

	"github.com/BigNoseCattyHome/aorb/backend/go-services/auth/models"
	"github.com/BigNoseCattyHome/aorb/backend/go-services/auth/services"
	"github.com/gin-gonic/gin"
)

// controller就是用来接收我们的API请求的，然后调用service层的方法，最后返回结果给客户端

// 注册
func Register(c *gin.Context) {
	var user models.User

	// 把user使用ShouldBindJSON方法解析请求的json数据
	if err := c.ShouldBindJSON(&user); err != nil {
		// 这里是一些错误处理
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Invalid request"})
		return
	}

	// 调用service层的RegisterUser方法，把user传进去
	// controller要做的就是把user解析出来，然后把user传递到Service中做具体的任务
	if err := services.RegisterUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to register user"})
		return
	}

	// 上面是把user传递给RegisterUser方法
	// 返回注册成功的消息
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User registered successfully",
	})
}

// 登录
// 这个在文档中写好了，就先写这个把
func Login(c *gin.Context) {
	var request models.RequestLogin // 根据文档中的请求数据结构，定义了一个RequestLogin的模型

	// 这里是gin框架的方法，用来解析请求的json数据，然后把解析后的数据放到user变量中
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// 这里Login方法中也是和上面的Register方法一样
	// 把user从请求的JSON中解析出来，传递给service中的方法做具体的操作
	token, exp_token, refresh_token, simple_user, err := services.AuthenticateUser(&request) // 这里就是把user作为参数传递给了AuthenticateUser方法
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// 定义了好返回的数据结构models/response.go之后，这里就进行具体的处理
	// 现在我们就是要把处理之后的User中的特定字段，放到Response中，然后返回给客户端
	// 具体赋值，把处理之后的user中的字段赋值给res中的字段
	res := models.ResponseLogin{
		Message:      "User login successfully",
		Success:      true,
		Token:        token, //这个就是经常念叨的JWT，前面生成好了这里直接赋值过来，生成token的代码在service中
		ExpiresIn:    exp_token,
		RefreshToken: refresh_token,
		User: models.SimpleUser{
			Avatar:    simple_user.Avatar,
			ID:        simple_user.ID,
			Ipaddress: simple_user.Ipaddress,
			Nickname:  simple_user.Nickname,
		},
	}
	c.JSON(http.StatusOK, res) // 然后就把这个res发出去就好啦

}

func Verify(c *gin.Context) {
	// 从请求头中获取token
	token := c.Request.Header.Get("Authorization")

	// 如果token为空，返回错误
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"valid": false, "error": "token is missing"})
		return
	}

	// 验证token是否有效
	claim, err := services.VerifyAccessToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"valid": false, "error": "invalid token"})
		return
	}

	res := models.ResponseVerify{
		Valid:    true,
		UserID:   claim.UserID,
		Username: claim.Username,
		Exp:      claim.ExpiresAt,
	}

	c.JSON(http.StatusOK, res)
}

// 登出
// 在客户端删除token，服务端吊销token
func Logout(c *gin.Context) {
	// 从请求头中获取token
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"valid": false, "error": "token is missing"})
		return
	}

	// 从请求头中获取refresh token
	refreshToken := c.Request.Header.Get("Refresh")
	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"valid": false, "error": "refresh token is missing"})
		return
	}

	// 验证访问令牌，确保合法用户的操作
	claim, err := services.VerifyAccessToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"valid": false, "error": "invalid token"})
		return
	}

	// 撤销刷新令牌
	err = services.RevokeRefreshToken(claim.UserID, refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"valid": false, "error": "failed to revoke refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true, "message": "tokens revoked"})
}

// 刷新token
func Refresh(c *gin.Context) {
	// 从请求头中获取refresh token
	refreshToken := c.Request.Header.Get("Refresh")
	if refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"valid": false, "error": "refresh token is missing"})
		return
	}

	// 调用服务层的RefreshAccessToken方法，传递refresh token
	newToken, exp_token, err := services.RefreshAccessToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"valid": false, "error": "failed to refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": true, "token": newToken, "expires_in": exp_token})
}
