package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lovebai/toalist/conf"
	"github.com/lovebai/toalist/router"
	"github.com/lovebai/toalist/utils"
)

func main() {
	pass := flag.String("pass", "", "生成local模式下的管理员密码，生成后自行填入到conf.ini配置文件中")
	flag.Parse()
	if *pass != "" {
		md5 := utils.MD5Encrypt(*pass)
		fmt.Println("生成的密码为：", md5)
		os.Exit(0)
	}

	// 初始化配置
	conf.InitConfig()

	if conf.GlobalConfig.Upload.Method == "stream" || conf.GlobalConfig.Upload.Method == "form" {
		utils.RefreshToken()
	}

	// 初始化路由
	router.InitRouter()

}
