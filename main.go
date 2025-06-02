package main

import (
	"github.com/lovebai/toalist/router"
	"github.com/lovebai/toalist/conf"
	"fmt"
	"github.com/lovebai/toalist/utils"
)

func init() {
	conf.InitConfig()
}

func main() {
	// 初始化配置
	conf.InitConfig()

	utils.RefreshToken()

	// 初始化路由
	router.InitRouter()

	
	fmt.Println(utils.GetToken())

}