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

// 404
func NullPage(c *gin.Context) {
	html := `<html><head><title>404 Not Found</title></head><body><center><h1>404 Not Found</h1></center>
<hr><center>ToAlist For Golang / <a href="/" style="text-decoration: none;color: #03A9F4;">Go To Home</a></center>
</body></html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}
