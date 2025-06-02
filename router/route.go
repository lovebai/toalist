package router

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lovebai/toalist/conf"
	"github.com/lovebai/toalist/controller"
)

func IndexPage(router *gin.Engine) {
	router.GET("/", controller.Index)
	router.GET("/i", NullPage)
	router.POST("/api/upload", controller.UploadForm)
}

// 404
func NullPage(c *gin.Context) {
	html := `<html><head><title>404 Not Found</title></head><body><center><h1>404 Not Found</h1></center>
<hr><center>ToAlist For Golang / <a href="/" style="text-decoration: none;color: #03A9F4;">Go To Home</a></center>
</body></html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
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
	router.Static("/i", "./i")
	router.Use(Middleware())
	router.NoRoute(NullPage)
	IndexPage(router)

	router.Run(":" + conf.GlobalConfig.Base.Port)
}
