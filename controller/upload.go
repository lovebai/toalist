package controller

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lovebai/toalist/conf"
	"github.com/lovebai/toalist/utils"
)

var config = conf.GlobalConfig

// 处理前端上传，转发到文件服务器
func UploadForm(c *gin.Context) {
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
		filepath := config.Alist.Path + "/" + file.Filename

		// 根据配置选择上传方式
		if config.Upload.Method == "stream" {
			// 流式上传
			uploadURL := config.Alist.APIURL + "/api/fs/put"
			req, err := http.NewRequest("PUT", uploadURL, bytes.NewReader(fileContent))
			if err != nil {
				errors = append(errors, file.Filename+": "+err.Error())
				continue
			}

			req.Header.Set("Content-Type", "application/octet-stream")
			req.Header.Set("File-Path", filepath)
			req.Header.Set("Authorization", utils.GetToken())

			client := &http.Client{}
			resp, _ = client.Do(req)
		} else {
			// 表单上传（默认方式）
			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			part, err := writer.CreateFormFile("file", file.Filename)
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
		}

		if err != nil {
			errors = append(errors, file.Filename+": "+err.Error())
			continue
		}

		// 读取响应
		_, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			errors = append(errors, file.Filename+": "+err.Error())
			continue
		}

		url, err := utils.GetFsLink(filepath)
		if err != nil {
			errors = append(errors, file.Filename+": "+err.Error())
			continue
		}

		results = append(results, FileResult{
			URL:  url,
			Name: file.Filename,
		})
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
