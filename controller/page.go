package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Data struct {
	Title string
	Desc  string
	Icon  string
}

// 首页
func Index(c *gin.Context) {
	data := Data{
		Title: config.Page.Title,
		Desc:  config.Page.Desc,
		Icon: func() string {
			if len(config.Page.Icon) == 0 {
				return "https://i.obai.cc/i/favicon.png"
			}
			return config.Page.Icon
		}(),
	}
	// c.HTML(http.StatusOK, "index", data) //为了模板兼容性，我改改改改改改
	c.HTML(http.StatusOK, "index.html", data)
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
