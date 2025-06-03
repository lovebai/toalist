package conf

import (
	"log/slog"
	"os"

	"gopkg.in/ini.v1"
)

const ConfigPath = "conf/conf.ini"

func initConf() {
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		dir := "conf"
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				slog.Error("Failed to create config directory", "error", err)
				os.Exit(1)
			}
		}

		file, err := os.Create(ConfigPath)
		if err != nil {
			slog.Error("Failed to create config file", "error", err)
			os.Exit(1)
		}
		defer file.Close()

		// 默认配置内容

		defaultConfig := `
[base]
mode = release
host = 0.0.0.0
port = 5245
url = http://127.0.0.1:5245

[alist]
alist_url = 
alist_api_url = 
alist_username = 
alist_password = 
alist_path = 

[upload]
upload_method = local 
allow_types = jpg,jpeg,png,gif,pdf,doc,docx,xls,xlsx,ppt,pptx,txt,zip,rar,7z
max_file_size = 50 
keep_original_name = true
local_upload_path = /i

[login]
username = admin
password = a9d38dbaadb53aad0c4b14570145c02f  
admin_page = /manager
`
		if _, err := file.WriteString(defaultConfig); err != nil {
			slog.Error("Failed to write default config", "error", err)
			os.Exit(1)
		}
		slog.Info("首次运行，已为您生成默认配置文件", "path", ConfigPath)
		slog.Info("本地模式后台默认管理员账号：admin 密码：admin123 ,请及时修改")
	}
}

func InitConfig() {
	slog.Info("提示：每次修改配置文件后，都需要重启服务。")
	initConf()
	cfg, err := ini.Load(ConfigPath)
	if err != nil {
		slog.Error("Failed to read configuration file", "error", err)
		os.Exit(1)
	}

	// 映射配置到结构体
	if err := cfg.MapTo(GlobalConfig); err != nil {
		slog.Error("Failed to map configuration", "error", err)
		os.Exit(1)
	}

}
