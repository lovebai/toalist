package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 首页
func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index", gin.H{
		"title":           "Toalist",
		"formUploadUrl":   "/api/upload/form",
		"streamUploadUrl": "/api/upload/stream",
	})
}
