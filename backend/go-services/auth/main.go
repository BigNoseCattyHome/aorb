package main

// 在main.go中启动路由
import (
	"os"

	"github.com/BigNoseCattyHome/aorb/backend/go-services/api-gateway/conf"
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
	// 初始化config配置文件
	if err := conf.LoadConfig(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	r := gin.Default()   // 这里是gin框架的方法，用来创建一个gin的实例
	routes.AuthRoutes(r) // 这里调用我们自己的routes.go中的AuthRoutes方法，传入r参数
	r.Run(":8080")       // 让gin框架监听8080端口
}
