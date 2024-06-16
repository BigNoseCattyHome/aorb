package routes

import (
	"github.com/BigNoseCattyHome/aorb/backend/go-services/auth/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine) {
	// r.Group方法注册了一个路由组，下面注册了两个子路由
	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register)	// 表示在路由处理到/auth/register时，调用controllers.Register方法
		auth.POST("/login", controllers.Login)			// 表示在路由处理到/auth/login时，调用controllers.Login方法
	}
}
