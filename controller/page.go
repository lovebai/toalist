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
	c.HTML(http.StatusOK, "index", data)
}
