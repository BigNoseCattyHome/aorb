package controllers

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-gonic/gin"
)

// 实现反向代理，将请求代理到目标地址
// 只把path后面的部分作为api传递，所以在微服务中的路由只用写path部分
func ProxyRequest(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		remote, err := url.Parse(target)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target URL"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		originalDirector := proxy.Director

		proxy.Director = func(req *http.Request) {
			originalDirector(req)
			req.URL.Path = c.Param("path")
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}
