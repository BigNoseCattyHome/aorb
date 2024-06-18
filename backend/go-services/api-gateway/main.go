package main

import (
	"log"

	"github.com/BigNoseCattyHome/aorb/backend/go-services/api-gateway/conf"
	"github.com/BigNoseCattyHome/aorb/backend/go-services/api-gateway/routes"
	"github.com/gin-gonic/gin"
)

func init() {
	// 初始化配置文件
	if err := conf.LoadConfig(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
}

func main() {
	// 读取配置api的前缀和版本
	apiPrefix := conf.AppConfig.GetString("api.prefix")
	apiVersion := conf.AppConfig.GetString("api.version")

	router := gin.Default()

	// 设置api前缀，返回一个新的RouterGroup
	apiGroup := router.Group(apiPrefix + "/" + apiVersion)

	// 注册路由，RegisterRoutes函数在routes/gateway_routes.go中定义
	routes.RegisterRoutes(apiGroup)

	// 启动服务器
	router.Run(":8080")
}
