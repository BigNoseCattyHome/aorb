package routes

import (
	"github.com/BigNoseCattyHome/aorb/backend/go-services/api-gateway/controllers"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册路由
func RegisterRoutes(router *gin.RouterGroup) {
	// 本地开发的时候都可以把地址设置为localhost
	router.GET("/auth/*path", controllers.ProxyRequest("http://auth-service:8080")) // auth可以暂时先不做
	router.GET("/message/*path", controllers.ProxyRequest("http://message-service:8080")) // 将/api/message/*的请求代理到message-service
	router.GET("/poll/*path", controllers.ProxyRequest("http://poll-service:8080"))
	router.GET("/recommendation/*path", controllers.ProxyRequest("http://recommendation-service:8080"))
	router.GET("/user/*path", controllers.ProxyRequest("http://java-user-service:8080")) // Java服务
}
