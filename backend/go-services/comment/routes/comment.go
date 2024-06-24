package routes

import (
	comment2 "github.com/BigNoseCattyHome/aorb/backend/go-services/comment/web"
	"github.com/gin-gonic/gin"
)

func RegisterCommentRoutes(g *gin.Engine) *gin.RouterGroup {
	commentGroup := g.Group("/comment")
	{
		commentGroup.POST("/action/", comment2.ActionCommentHandler)
		commentGroup.GET("/list/", comment2.ListCommentHandler)
		commentGroup.GET("/count/", comment2.CountCommentHandler)
	}
	return commentGroup
}
