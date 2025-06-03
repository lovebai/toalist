package main

import (
	"embed"
	"flag"
	"fmt"
	"os"

	"github.com/lovebai/toalist/conf"
	"github.com/lovebai/toalist/router"
	"github.com/lovebai/toalist/utils"
)

//go:embed views/*
var views embed.FS

func main() {
	md5 := flag.String("md5", "", "生成local模式下的管理员密码，生成后自行填入到conf.ini配置文件中")
	aes := flag.String("aes", "", "对alist的用户密码进行加密，生成后自行填入到conf.ini配置文件中")
	flag.Parse()
	if *md5 != "" {
		en := utils.MD5Encrypt(*md5)
		fmt.Println("生成的密码为：", en)
		os.Exit(0)
	}
	if *aes != "" {
		en, _ := utils.AESEncrypt(*aes)
		fmt.Println("加密后的字符为：", en)
		os.Exit(0)
	}

	// 初始化配置
	conf.InitConfig()

	if conf.GlobalConfig.Upload.Method == "stream" || conf.GlobalConfig.Upload.Method == "form" {
		utils.RefreshToken()
	}

	// 初始化路由和模板
	router.InitRouter(views)

}
