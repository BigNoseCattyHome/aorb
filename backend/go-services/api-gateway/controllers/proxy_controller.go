package controllers

import (
    "net/http"
    "net/http/httputil"
    "net/url"

    "github.com/gin-gonic/gin"
)

func ProxyRequest(target string) gin.HandlerFunc {
    return func(c *gin.Context) {
        remote, err := url.Parse(target)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid target URL"})
            return
        }
        proxy := httputil.NewSingleHostReverseProxy(remote)
        proxy.ServeHTTP(c.Writer, c.Request)
    }
}
