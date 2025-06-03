package conf

var GlobalConfig = &ConfigType{}

type ConfigType struct {
	Base   BaseConfig   `ini:"base"`
	Alist  AlistConfig  `ini:"alist"`
	Upload UploadConfig `ini:"upload"`
	Login  LoginConfig  `ini:"login"`
}

type BaseConfig struct {
	Mode string `ini:"mode"`
	Host string `ini:"host"`
	Port string `ini:"port"`
	Url  string `ini:"url"`
}

type AlistConfig struct {
	URL      string `ini:"alist_url"`
	APIURL   string `ini:"alist_api_url"`
	Username string `ini:"alist_username"`
	Password string `ini:"alist_password"`
	Path     string `ini:"alist_path"`
}

type UploadConfig struct {
	Method           string `ini:"upload_method"`
	AllowTypes       string `ini:"allow_types"`
	MaxFileSize      int    `ini:"max_file_size"`
	KeepOriginalName bool   `ini:"keep_original_name"`
	LocalUploadPath  string `ini:"local_upload_path"`
}

type LoginConfig struct {
	Username  string `ini:"username"`
	Password  string `ini:"password"`
	AdminPage string `ini:"admin_page"`
}
