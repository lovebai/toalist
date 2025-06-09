package router

import (
	"embed"
	"html/template"
	"log/slog"
	"net"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/lovebai/toalist/conf"
	"github.com/lovebai/toalist/controller"
)

// 检查端口是否可用
func isPortAvailable(host string, port string) bool {
	ln, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		return false
	}
	ln.Close()
	return true
}

// 文件管理相关的路由
func FileManagerRoutes(r *gin.Engine) {
	managerPage := r.Group(conf.GlobalConfig.Login.AdminPage)
	{
		// 页面路由
		managerPage.GET("/login", controller.LoginPage)
		managerPage.POST("/login", controller.Login)
		managerPage.GET("/", controller.FileManagerPage)
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
	router.GET(conf.GlobalConfig.Upload.LocalUploadPath, controller.NullPage)
	router.POST("/api/upload", controller.UploadForm)
	router.GET("/hide/*hide", controller.HideImage)
}

// 代理路由
func ProxyRoutes(r *gin.Engine) {
	// 创建代理路由组
	proxyGroup := r.Group("/dw")
	{
		proxyGroup.Any("/*path", controller.Proxy)
	}
}

// 初始化路由和模板
func InitRouter(views embed.FS) {
	gin.SetMode(conf.GlobalConfig.Base.Mode)
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"})
	router.Delims("[[", "]]")

	tmpl := template.Must(template.New("").ParseFS(views, "views/*.html"))

	if conf.GlobalConfig.Base.Mode == "debug" {
		router.LoadHTMLGlob("views/**")
	} else {
		router.SetHTMLTemplate(tmpl)
	}

	if conf.GlobalConfig.Upload.Method == "local" {
		router.Static(conf.GlobalConfig.Upload.LocalUploadPath, "./"+conf.GlobalConfig.Upload.LocalUploadPath)
		controller.InitLoginInfo()
		controller.InitManager()
		FileManagerRoutes(router)
		slog.Info("已选择本地(local)模式，上传文件路径为：" + conf.GlobalConfig.Upload.LocalUploadPath)
	}

	// 根据配置是否启用代理
	if conf.GlobalConfig.Alist.IsProxy && conf.GlobalConfig.Upload.Method != "local" {
		ProxyRoutes(router)
		slog.Info("已启用对Alist的代理，代理地址为：/dw")
	}

	router.NoRoute(controller.NullPage)
	IndexPage(router)

	server := conf.GlobalConfig.Base.Host + ":" + conf.GlobalConfig.Base.Port

	if !isPortAvailable(conf.GlobalConfig.Base.Host, conf.GlobalConfig.Base.Port) {
		slog.Error("当前端口已被占用，请更换端口号后重启服务")
		os.Exit(1)
	}

	slog.Info("ToAlist 服务已启动在：http://" + server)
	router.Run(server)
}
