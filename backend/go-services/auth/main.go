package main

import (
	"os"

	"github.com/BigNoseCattyHome/aorb/backend/go-services/auth/routes"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func init() {
	// 设置日志级别为 Debug
	log.SetLevel(log.DebugLevel)

	// 设置日志输出格式为 JSON
	log.SetFormatter(&log.JSONFormatter{})

	// 设置日志输出到标准输出
	log.SetOutput(os.Stdout)

}

func main() {

	router := gin.Default() // 这里是gin框架的方法，用来创建一个gin的实例
	routes.AuthRoutes(router)
	router.Run(":8081") // 让gin框架监听8081端口
}
