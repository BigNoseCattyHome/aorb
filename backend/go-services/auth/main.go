package main

// 在main.go中启动路由
import (
	"github.com/BigNoseCattyHome/aorb/backend/go-services/api-gateway/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()   // 这里是gin框架的方法，用来创建一个gin的实例
	routes.AuthRoutes(r) // 这里调用我们自己的routes.go中的AuthRoutes方法，传入r参数
	r.Run(":8080")       // 让gin框架监听8080端口
}
