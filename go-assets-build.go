// +build ignore

package main

import (
	"log"
	"net/http"

	"github.com/shurcooL/vfsgen"
)

func main() {
	// 1. 指定要打包的文件夹路径 (这里是 "assets" 文件夹)
	var fs http.FileSystem = http.Dir("static")

	// 2. 配置生成选项
	err := vfsgen.Generate(fs, vfsgen.Options{
		// 生成的文件所属的包名
		PackageName:  "src",
		// 生成的文件名 (默认为 assets_vfsdata.go)
		Filename:     "src/assets.go",
		// 生成的代码中，存储文件数据的变量名
		VariableName: "Assets",
	})

	if err != nil {
		log.Fatalln(err)
	}
}