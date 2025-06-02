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

// UploadForm 处理前端表单上传，并转发到文件服务器
func UploadForm(c *gin.Context) {
	// 获取所有上传的文件
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

		// 创建multipart表单
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		// 创建文件表单字段
		part, err := writer.CreateFormFile("file", file.Filename)
		if err != nil {
			src.Close()
			errors = append(errors, file.Filename+": "+err.Error())
			continue
		}

		// 复制文件内容到表单
		_, err = io.Copy(part, src)
		src.Close()
		if err != nil {
			errors = append(errors, file.Filename+": "+err.Error())
			continue
		}

		// 关闭writer
		err = writer.Close()
		if err != nil {
			errors = append(errors, file.Filename+": "+err.Error())
			continue
		}

		// 构建文件服务器URL
		uploadURL := config.Alist.APIURL + "/api/fs/form"

		// 创建转发请求
		req, err := http.NewRequest("PUT", uploadURL, body)
		if err != nil {
			errors = append(errors, file.Filename+": "+err.Error())
			continue
		}

		filepath := config.Alist.Path + "/" + file.Filename

		// 设置请求头
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("File-Path", filepath)
		req.Header.Set("Authorization", utils.GetToken())

		// 发送请求到文件服务器
		client := &http.Client{}
		resp, err := client.Do(req)
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

// UploadStream 处理前端流式上传，并转发到文件服务器
func UploadStream(c *gin.Context) {
	// 获取文件名
	filename := c.Query("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少文件名参数",
		})
		return
	}

	// 读取请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "读取请求体失败",
			"error":   err.Error(),
		})
		return
	}

	// 构建文件服务器URL
	uploadURL := config.Alist.APIURL + "/api/fs/put"

	// 创建转发请求
	req, err := http.NewRequest("PUT", uploadURL, bytes.NewReader(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建请求失败",
			"error":   err.Error(),
		})
		return
	}

	filepath := config.Alist.Path + "/" + filename

	// 设置请求头
	req.Header.Set("Content-Type", c.GetHeader("Content-Type"))
	req.Header.Set("File-Path", filepath)
	req.Header.Set("Authorization", utils.GetToken())

	// 发送请求到文件服务器
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "发送请求到文件服务器失败",
			"error":   err.Error(),
		})
		return
	}
	defer resp.Body.Close()

	// 读取响应
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "读取文件服务器响应失败",
			"error":   err.Error(),
		})
		return
	}

	url, err := utils.GetFsLink(filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "读取文件服务器响应失败",
			"error":   err,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "文件上传成功",
		"url":     url,
		"name":    filename,
	})

}
