package controller

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

var BaseUploadDir string

// InitManager 初始化文件管理器
func InitManager() {
	BaseUploadDir = strings.ReplaceAll(config.Upload.LocalUploadPath, "/", "")
	// 确保上传目录存在
	if err := os.MkdirAll(BaseUploadDir, 0755); err != nil {
		panic(fmt.Sprintf("创建上传目录失败: %v", err))
	}
}

// FileInfo 文件信息结构体
type FileInfo struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	IsDir   bool   `json:"isDir"`
	Size    int64  `json:"size"`
	ModTime string `json:"modTime"`
}

// CreateDirectory 创建目录
func CreateDirectory(c *gin.Context) {
	dirPath := c.PostForm("path")
	if dirPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "目录路径不能为空"})
		return
	}

	// 清理和规范化路径
	dirPath = filepath.Clean(dirPath)
	if dirPath == "." || dirPath == ".." || strings.HasPrefix(dirPath, "..") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的目录路径"})
		return
	}

	// 确保基础目录存在
	if err := os.MkdirAll(BaseUploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("创建基础目录失败: %v", err)})
		return
	}

	// 确保路径在基础目录下
	fullPath := filepath.Join(BaseUploadDir, dirPath)
	absBasePath, err := filepath.Abs(BaseUploadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("无法解析基础目录路径: %v", err)})
		return
	}

	absFullPath, err := filepath.Abs(fullPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("无法解析目标目录路径: %v", err)})
		return
	}

	// 确保目标路径在基础目录内
	if !strings.HasPrefix(absFullPath, absBasePath) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "目录路径超出允许范围",
			"debug": gin.H{
				"basePath": absBasePath,
				"fullPath": absFullPath,
			},
		})
		return
	}

	// 检查目录是否已存在
	if fileInfo, err := os.Stat(fullPath); err == nil {
		if fileInfo.IsDir() {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "目录已存在",
				"path":  dirPath,
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "指定路径已存在同名文件",
				"path":  dirPath,
			})
		}
		return
	} else if !os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("检查目录状态失败: %v", err),
			"path":  dirPath,
		})
		return
	}

	// 创建目录
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("创建目录失败: %v", err),
			"debug": gin.H{
				"path":     dirPath,
				"fullPath": fullPath,
			},
		})
		return
	}

	// 验证目录是否真的创建成功
	if _, err := os.Stat(fullPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "目录创建后验证失败",
			"debug": gin.H{
				"path":     dirPath,
				"fullPath": fullPath,
				"error":    err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "目录创建成功",
		"path":    dirPath,
		"debug": gin.H{
			"fullPath": fullPath,
		},
	})
}

// UploadFile 上传文件
func UploadFile(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取表单数据失败"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未选择文件"})
		return
	}

	dirPath := c.PostForm("path")
	if dirPath == "" {
		dirPath = "."
	}

	// 确保路径在基础目录下
	fullPath := filepath.Join(BaseUploadDir, dirPath)
	if !strings.HasPrefix(fullPath, BaseUploadDir) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的目录路径"})
		return
	}

	// 确保目标目录存在
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("创建目录失败: %v", err)})
		return
	}

	var uploadedFiles []string
	var failedFiles []string

	for _, file := range files {
		// 保存文件
		dst := filepath.Join(fullPath, file.Filename)
		if err := c.SaveUploadedFile(file, dst); err != nil {
			failedFiles = append(failedFiles, file.Filename)
			continue
		}
		uploadedFiles = append(uploadedFiles, file.Filename)
	}

	if len(failedFiles) > 0 {
		c.JSON(http.StatusPartialContent, gin.H{
			"message": "部分文件上传成功",
			"success": uploadedFiles,
			"failed":  failedFiles,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "所有文件上传成功",
		"files":   uploadedFiles,
	})
}

// DeleteFileOrDir 删除文件或目录
func DeleteFileOrDir(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "路径不能为空"})
		return
	}

	// 确保路径在基础目录下
	fullPath := filepath.Join(BaseUploadDir, path)
	if !strings.HasPrefix(fullPath, BaseUploadDir) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的路径"})
		return
	}

	if err := os.RemoveAll(fullPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("删除失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ListDirectory 列出目录内容
func ListDirectory(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		path = "."
	}

	// 确保路径在基础目录下
	fullPath := filepath.Join(BaseUploadDir, path)
	if !strings.HasPrefix(fullPath, BaseUploadDir) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的目录路径"})
		return
	}

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("读取目录失败: %v", err)})
		return
	}

	var files []FileInfo
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		fileInfo := FileInfo{
			Name:    entry.Name(),
			Path:    filepath.Join(path, entry.Name()),
			IsDir:   entry.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime().Format("2006-01-02 15:04:05"),
		}
		files = append(files, fileInfo)
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}
