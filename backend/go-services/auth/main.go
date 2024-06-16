package main

// 在main.go中启动路由
import (
	"os"

	"github.com/BigNoseCattyHome/aorb/backend/go-services/auth/conf"
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
	// 读取配置api的前缀和版本
	apiPrefix := conf.AppConfig.GetString("api.prefix")
	apiVersion := conf.AppConfig.GetString("api.version")

	router := gin.Default() // 这里是gin框架的方法，用来创建一个gin的实例

	// 设置api前缀，返回一个新的RouterGroup
	apiGroup := router.Group(apiPrefix + "/" + apiVersion)

	// 注册路由，RegisterRoutes函数在routes/gateway_routes.go中定义
	routes.AuthRoutes(apiGroup)

	router.Run(":8080") // 让gin框架监听8080端口
}
