package routes

import (
	"github.com/BigNoseCattyHome/aorb/backend/go-services/auth/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)
	router.POST("/verify", controllers.Verify)
	router.POST("/logout", controllers.Logout)
	router.POST("/refresh", controllers.Refresh)
}
