package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// V1BignosecatGet - 开发者信息
func V1BignosecatGet(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
