package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// V1PollsGet - 获取所有的投票
func V1PollsGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// V1PollsPost - 创建新的投票
func V1PollsPost(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// V1PollsQuestionGet - 获取单个投票详情
func V1PollsQuestionGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// V1PollsQuestionPost - 为问题投票
func V1PollsQuestionPost(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// V1UserBlacklistPost - 拉黑用户
func V1UserBlacklistPost(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// V1UserFollowerPost - 关注用户/取关用户
func V1UserFollowerPost(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// V1UserInfoGet - 获取用户基本信息
func V1UserInfoGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}

// V1UserLoginGet - 登录
func V1UserLoginGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
