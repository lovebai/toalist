package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 首页
func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

// 登录页面
func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// 文件管理页面
func FileManagerPage(c *gin.Context) {
	c.HTML(http.StatusOK, "file-manager.html", gin.H{
		"upload_path": config.Upload.LocalUploadPath + "/",
	})
}
