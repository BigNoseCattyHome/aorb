package main

import (
	"context"
	"github.com/BigNoseCattyHome/aorb/backend/go-services/api-gateway/middleware"
	auth2 "github.com/BigNoseCattyHome/aorb/backend/go-services/auth/web"
	comment2 "github.com/BigNoseCattyHome/aorb/backend/go-services/comment/web"
	user2 "github.com/BigNoseCattyHome/aorb/backend/go-services/user/web"
	"github.com/BigNoseCattyHome/aorb/backend/utils/constants/config"
	"github.com/BigNoseCattyHome/aorb/backend/utils/extra/tracing"
	"github.com/BigNoseCattyHome/aorb/backend/utils/logging"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginprometheus "github.com/zsais/go-gin-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"net/http"
)

func main() {
	tp, err := tracing.SetTraceProvider(config.WebServerName)
	if err != nil {
		logging.Logger.WithFields(logrus.Fields{
			"err": err,
		}).Panicf("Error to set the trace")
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			logging.Logger.WithFields(logrus.Fields{
				"err": err,
			}).Errorf("Error to set the trace")
		}
	}()

	g := gin.Default()
	p := ginprometheus.NewPrometheus("aorb-WebGateway")
	p.Use(g)
	g.Use(gzip.Gzip(gzip.DefaultCompression))
	g.Use(otelgin.Middleware(config.WebServerName))
	g.Use(middleware.TokenAuthMiddleware())

	// prometheus，检测有关内存分配、垃圾回收效率、并发级别的洞察，并帮助识别性能瓶颈或潜在的内存泄漏
	//http.Handle("/metrics", promhttp.Handler())
	//go func() {
	//	http.ListenAndServe(":30000", nil)
	//}()

	// 注册路由
	g.GET("/ping", func(c *gin.Context) {
		//c.Redirect(http.StatusPermanentRedirect, "/login")
		c.JSON(http.StatusOK, "success")
	})
	rootPath := g.Group("/aorb")
	user := rootPath.Group("/user")
	{
		user.GET("/", user2.UserHandler)
		user.POST("/login/", auth2.LoginHandler)
		user.POST("/register/", auth2.RegisterHandler)
	}

	//feed := rootPath.Group("/feed")
	//{
	//	feed.GET("/", feed2.ListVideosByRecommendHandle)
	//}

	comment := rootPath.Group("/comment")
	{
		comment.POST("/action/", comment2.ActionCommentHandler)
		comment.GET("/list/", comment2.ListCommentHandler)
		comment.GET("/count/", comment2.CountCommentHandler)
	}

	// run
	if err := g.Run(config.WebServerAddr); err != nil {
		panic("Can not run aorb Gateway, binding port: " + config.WebServerAddr)
	}
}
