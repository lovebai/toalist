package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/lovebai/toalist/conf"
	"github.com/lovebai/toalist/controller"
)

func IndexPage(router *gin.Engine) {
	router.GET("/", controller.Index)
	router.POST("/api/upload", controller.UploadForm)
}

// 中间件
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		slog.Info("Middleware")
		c.Next()
	}
}

// 初始化路由
func InitRouter() {
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})
	router.Delims("[[", "]]")
	router.LoadHTMLGlob("views/**")
	router.Use(Middleware())
	IndexPage(router)

	router.Run(":" + conf.GlobalConfig.Base.Port)
}
