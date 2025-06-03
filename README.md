# Toalist

Toalist 是一个基于 Go 语言开发的 Web 应用程序，用于管理和操作 Alist 文件系统。该项目使用 Gin 框架构建，提供了简单易用的 Web 界面来管理文件。

## 功能特性

- 支持文件上传（支持流式上传和表单上传两种模式）
- 集成 Alist 文件系统
- 安全的密码加密机制（支持 MD5 和 AES 加密）
- 基于 JWT 的身份认证
- 响应式 Web 界面
- 配置文件支持（使用 INI 格式）

## 系统要求

- Go 1.24.0 或更高版本
- 支持的操作系统：Windows、Linux、macOS

## 安装说明

1. 克隆项目
```bash
git clone https://github.com/lovebai/toalist.git
cd toalist
```

2. 安装依赖
```bash
go mod download
```

3. 配置
- 复制 `conf.ini` 文件并根据需要修改配置
- 使用以下命令生成管理员密码：
```bash
go run main.go -md5 "你的密码"
```
- 使用以下命令加密 Alist 用户密码：
```bash
go run main.go -aes "需要加密的密码"
```

4. 编译运行
```bash
go build
./toalist
```

## 项目结构

```
toalist/
├── conf/           # 配置相关代码
├── controller/     # 控制器
├── router/         # 路由定义
├── utils/          # 工具函数
├── views/          # 前端模板
├── i/              # 静态资源
├── conf.ini        # 配置文件
├── main.go         # 主程序入口
├── go.mod          # Go 模块定义
└── go.sum          # 依赖版本锁定
```

## 配置说明

配置文件 `conf.ini` 采用 INI 格式，包含以下主要配置项：

### 基础配置 [base]
```ini
[base]
mode = release        # 运行模式：release（生产）或 debug（调试）
host = 0.0.0.0       # 服务器监听地址，0.0.0.0 表示监听所有网络接口
port = 5245          # 服务器监听端口
url = http://127.0.0.1:5245  # 服务器访问地址
```

### Alist 配置 [alist]
```ini
[alist]
alist_url = http://10.10.11.143:5244      # Alist 服务器地址
alist_api_url = http://10.10.11.143:5244  # Alist API 地址
alist_username = admin                    # Alist 用户名
alist_password = [加密后的密码]            # Alist 密码（使用 AES 加密）
alist_path = /test                        # Alist 存储路径
```

### 上传配置 [upload]
```ini
[upload]
upload_method = form                      # 上传方式：form（表单上传）或 stream（流式上传）
allow_types = jpg,jpeg,png,...            # 允许上传的文件类型，用逗号分隔
max_file_size = 50                        # 最大文件大小（单位：MB）
keep_original_name = false                # 是否保持原始文件名
local_upload_path = /i                    # 本地文件上传路径
```

### 登录配置 [login]
```ini
[login]
username = admin                          # 管理员用户名
password = [MD5加密后的密码]               # 管理员密码（使用 MD5 加密）
admin_page = /manager                     # 管理页面路径
```

### 配置说明

1. **密码加密**
   - 管理员密码使用 MD5 加密，可通过命令行生成：
     ```bash
     go run main.go -md5 "你的密码"
     ```
   - Alist 密码使用 AES 加密，可通过命令行生成：
     ```bash
     go run main.go -aes "需要加密的密码"
     ```

2. **上传方式**
   - `form`：使用传统的表单上传，适合小文件
   - `stream`：使用流式上传，适合大文件

3. **文件类型限制**
   - 默认支持的文件类型包括：jpg、jpeg、png、gif、pdf、doc、docx、xls、xlsx、ppt、pptx、txt、zip、rar、7z
   - 可以根据需要在 `allow_types` 中修改

4. **安全建议**
   - 生产环境中建议修改默认端口
   - 使用强密码并定期更换
   - 限制 `allow_types` 只包含必要的文件类型
   - 根据实际需求设置合适的 `max_file_size`

## 使用说明

1. 启动服务后，访问 `http://localhost:端口号` 进入 Web 界面
2. 使用配置的管理员账号登录
3. 通过 Web 界面上传和管理文件

## 开发说明

- 使用 `go mod` 管理依赖
- 遵循 Go 标准项目布局
- 使用 Gin 框架处理 HTTP 请求
- 使用 JWT 进行身份认证

## 许可证

[MIT License](LICENSE)

## 贡献指南

欢迎提交 Issue 和 Pull Request 来帮助改进项目。

## 联系方式

如有问题或建议，请通过 GitHub Issues 提交。
