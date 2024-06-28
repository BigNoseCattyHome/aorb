// api-gateway
package routes

import (
	"github.com/BigNoseCattyHome/aorb/backend/api-gateway/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.RouterGroup) {

	router.GET("/auth/*path", controllers.ProxyRequest("http://localhost:8081"))
	router.POST("/auth/*path", controllers.ProxyRequest("http://localhost:8081"))
	router.PUT("/auth/*path", controllers.ProxyRequest("http://localhost:8081"))
	router.DELETE("/auth/*path", controllers.ProxyRequest("http://localhost:8081"))

	router.GET("/message/*path", controllers.ProxyRequest("http://localhost:8082"))
	router.POST("/message/*path", controllers.ProxyRequest("http://localhost:8082"))
	router.PUT("/message/*path", controllers.ProxyRequest("http://localhost:8082"))
	router.DELETE("/message/*path", controllers.ProxyRequest("http://localhost:8082"))

	router.GET("/poll/*path", controllers.ProxyRequest("http://localhost:8083"))
	router.POST("/poll/*path", controllers.ProxyRequest("http://localhost:8083"))
	router.PUT("/poll/*path", controllers.ProxyRequest("http://localhost:8083"))
	router.DELETE("/poll/*path", controllers.ProxyRequest("http://localhost:8083"))

	router.GET("/recommendation/*path", controllers.ProxyRequest("http://localhost:8084"))
	router.POST("/recommendation/*path", controllers.ProxyRequest("http://localhost:8084"))
	router.PUT("/recommendation/*path", controllers.ProxyRequest("http://localhost:8084"))
	router.DELETE("/recommendation/*path", controllers.ProxyRequest("http://localhost:8084"))

	router.GET("/user/*path", controllers.ProxyRequest("http://localhost:8085"))
	router.POST("/user/*path", controllers.ProxyRequest("http://localhost:8085"))
	router.PUT("/user/*path", controllers.ProxyRequest("http://localhost:8085"))
	router.DELETE("/user/*path", controllers.ProxyRequest("http://localhost:8085"))
}
