package main

import (
	"github.com/lovebai/toalist/conf"
	"github.com/lovebai/toalist/router"
	"github.com/lovebai/toalist/utils"
)

func main() {
	// 初始化配置
	conf.InitConfig()

	if conf.GlobalConfig.Upload.Method == "stream" || conf.GlobalConfig.Upload.Method == "form" {
		utils.RefreshToken()
	}

	// 初始化路由
	router.InitRouter()

}
