package controller

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lovebai/toalist/conf"
	"github.com/lovebai/toalist/utils"
)

var config = conf.GlobalConfig

// 处理前端上传，转发到文件服务器
func UploadForm(c *gin.Context) {
	// 验证 sign 参数
	sign := c.PostForm("sign")
	if sign == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少必要的 sign 参数",
		})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "获取表单数据失败",
			"error":   err.Error(),
		})
		return
	}

	files := form.File["file"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "未找到上传的文件",
		})
		return
	}

	type FileResult struct {
		URL  string `json:"url"`
		Name string `json:"name"`
	}

	var results []FileResult
	var errors []string

	for _, file := range files {
		// 处理文件名
		processedFileName := utils.ProcessFileName(file.Filename)

		// 验证文件类型和大小
		if err := utils.ValidateFile(file); err != nil {
			errors = append(errors, file.Filename+": "+err.Error())
			continue
		}

		// 打开文件
		src, err := file.Open()
		if err != nil {
			errors = append(errors, file.Filename+": "+err.Error())
			continue
		}
		defer src.Close()

		// 读取文件内容
		fileContent, err := io.ReadAll(src)
		if err != nil {
			errors = append(errors, file.Filename+": "+err.Error())
			continue
		}

		var resp *http.Response
		var url string

		switch config.Upload.Method {
		case "stream":
			// 流式上传到文件服务器
			uploadURL := config.Alist.APIURL + "/api/fs/put"
			req, err := http.NewRequest("PUT", uploadURL, bytes.NewReader(fileContent))
			if err != nil {
				errors = append(errors, file.Filename+": "+err.Error())
				continue
			}

			filepath := config.Alist.Path + "/" + processedFileName
			req.Header.Set("Content-Type", "application/octet-stream")
			req.Header.Set("File-Path", filepath)
			req.Header.Set("Authorization", utils.GetToken())

			client := &http.Client{}
			resp, _ = client.Do(req)
			if resp != nil {
				url, _ = utils.GetFsLink(filepath)
			}

		case "form":
			// 表单上传到文件服务器
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			part, err := writer.CreateFormFile("file", processedFileName)
			if err != nil {
				errors = append(errors, file.Filename+": "+err.Error())
				continue
			}

			_, err = io.Copy(part, bytes.NewReader(fileContent))
			if err != nil {
				errors = append(errors, file.Filename+": "+err.Error())
				continue
			}

			err = writer.Close()
			if err != nil {
				errors = append(errors, file.Filename+": "+err.Error())
				continue
			}

			filepath := config.Alist.Path + "/" + processedFileName
			uploadURL := config.Alist.APIURL + "/api/fs/form"
			req, err := http.NewRequest("PUT", uploadURL, body)
			if err != nil {
				errors = append(errors, file.Filename+": "+err.Error())
				continue
			}

			req.Header.Set("Content-Type", writer.FormDataContentType())
			req.Header.Set("File-Path", filepath)
			req.Header.Set("Authorization", utils.GetToken())

			client := &http.Client{}
			resp, _ = client.Do(req)
			if resp != nil {
				url, _ = utils.GetFsLink(filepath)
			}

		case "local":
			// 存储到本地
			now := time.Now()
			year := fmt.Sprintf("%d", now.Year())
			month := fmt.Sprintf("%02d", now.Month())

			// 构建存储路径：i/年/月/文件名
			storePath := filepath.Join(strings.ReplaceAll(config.Upload.LocalUploadPath, "/", ""), year, month)

			// 确保目录存在
			if err := os.MkdirAll(storePath, 0755); err != nil {
				errors = append(errors, file.Filename+": 创建存储目录失败 - "+err.Error())
				continue
			}

			// 构建完整的文件路径
			fullPath := filepath.Join(storePath, processedFileName)

			// 写入文件
			if err := os.WriteFile(fullPath, fileContent, 0644); err != nil {
				errors = append(errors, file.Filename+": 写入文件失败 - "+err.Error())
				continue
			}

			// 构建访问URL
			url = config.Base.Url + "/" + filepath.ToSlash(fullPath)
		default:
			// 未开启上传功能
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "管理员未开启上传功能",
			})
			return
		}

		if url != "" {
			if config.Alist.IsProxy {
				url = strings.ReplaceAll(url, config.Alist.URL, config.Base.Url+"/dw")
			}

			results = append(results, FileResult{
				URL:  url,
				Name: processedFileName,
			})
		} else {
			errors = append(errors, file.Filename+": 上传失败")
		}

		if resp != nil {
			resp.Body.Close()
		}
	}

	// 构建响应
	response := gin.H{
		"code":    200,
		"message": "文件上传完成",
		"files":   results,
	}

	if len(errors) > 0 {
		response["code"] = 207 // 使用207状态码表示部分成功
		response["errors"] = errors
	}

	c.JSON(http.StatusOK, response)
}
