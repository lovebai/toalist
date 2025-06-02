package conf

import (
	"log/slog"
	"os"

	"gopkg.in/ini.v1"
)

var GlobalConfig = &ConfigType{}

type ConfigType struct {
	Base   BaseConfig   `ini:"base"`
	Alist  AlistConfig  `ini:"alist"`
	Upload UploadConfig `ini:"upload"`
	Page   PageConfig   `ini:"page"`
}

type BaseConfig struct {
	Mode    string `ini:"mode"`
	Host    string `ini:"host"`
	Port    string `ini:"port"`
	SiteUrl string `ini:"site_url"`
}

type AlistConfig struct {
	URL      string `ini:"alist_url"`
	APIURL   string `ini:"alist_api_url"`
	Username string `ini:"alist_username"`
	Password string `ini:"alist_password"`
	Path     string `ini:"alist_path"`
}

type UploadConfig struct {
	Method string `ini:"upload_method"`
}

type PageConfig struct {
	CustomEnable bool   `ini:"custom_enable"`
	CustomCSS    string `ini:"custom_css"`
	CustomJS     string `ini:"custom_js"`
}

// 默认配置
var defaultConfig = ConfigType{
	Base: BaseConfig{
		Mode: "development",
		Host: "127.0.0.1",
		Port: "5245",
	},
	Upload: UploadConfig{
		Method: "form",
	},
	Page: PageConfig{
		CustomEnable: false,
	},
}

func InitConfig() {
	cfg, err := ini.Load("conf.ini")
	if err != nil {
		slog.Error("Failed to read configuration file: %v", err)
		os.Exit(1)
	}

	// 设置默认值
	*GlobalConfig = defaultConfig

	// 映射配置到结构体
	if err := cfg.MapTo(GlobalConfig); err != nil {
		slog.Error("Failed to map configuration: %v", err)
		os.Exit(1)
	}

}
