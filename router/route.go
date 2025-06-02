package router

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lovebai/toalist/conf"
	"github.com/lovebai/toalist/controller"
)

// 文件管理相关的路由
func FileManagerRoutes(r *gin.Engine) {
	managerPage := r.Group(conf.GlobalConfig.Login.AdminPage)
	{
		// 页面路由
		managerPage.GET("/login", controller.LoginPage)
		managerPage.GET("/file-manager", controller.FileManagerPage)
		managerPage.POST("/login", controller.Login)
	}

	// 文件管理API路由组
	fileManager := r.Group("/file")
	fileManager.Use(controller.AuthMiddleware())
	{
		fileManager.POST("/mkdir", controller.CreateDirectory)
		fileManager.POST("/upload", controller.UploadFile)
		fileManager.DELETE("/delete", controller.DeleteFileOrDir)
		fileManager.GET("/list", controller.ListDirectory)
	}
}

// 首页路由，不分组了
func IndexPage(router *gin.Engine) {
	router.GET("/", controller.Index)
	router.GET(conf.GlobalConfig.Upload.LocalUploadPath, NullPage)
	router.POST("/api/upload", controller.UploadForm)
}

// 404
func NullPage(c *gin.Context) {
	html := `<html><head><title>404 Not Found</title></head><body><center><h1>404 Not Found</h1></center>
<hr><center>ToAlist For Golang / <a href="/" style="text-decoration: none;color: #03A9F4;">Go To Home</a></center>
</body></html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

// 初始化路由
func InitRouter() {
	gin.SetMode(conf.GlobalConfig.Base.Mode)
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})
	router.Delims("[[", "]]")
	router.LoadHTMLGlob("views/**")

	if conf.GlobalConfig.Upload.Method == "local" {
		router.Static(conf.GlobalConfig.Upload.LocalUploadPath, "./"+conf.GlobalConfig.Upload.LocalUploadPath)
		controller.InitLoginInfo()
		controller.InitManager()
		FileManagerRoutes(router)
		slog.Info("已选择本地(local)模式，上传文件路径为：" + conf.GlobalConfig.Upload.LocalUploadPath)
	}

	router.NoRoute(NullPage)
	IndexPage(router)

	server := conf.GlobalConfig.Base.Host + ":" + conf.GlobalConfig.Base.Port
	slog.Info("ToAlist 服务已启动在：" + server)
	router.Run(server)
}
