package controller

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lovebai/toalist/utils"
)

// 显示隐藏图片的
func HideImage(c *gin.Context) {
	en := c.Param("hide")
	en = strings.TrimPrefix(en, "/")
	path, _ := utils.AESDecrypt(en)

	if !utils.IsImageFile(path) {
		c.JSON(400, gin.H{
			"code": 400,
			"msg":  "当前文件不存在",
		})
		return
	}

	c.File(path)
}
