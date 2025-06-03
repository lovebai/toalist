package utils

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/lovebai/toalist/conf"
)

// 验证文件类型和大小
func ValidateFile(file *multipart.FileHeader) error {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(file.Filename), "."))

	// 验证文件类型
	allowTypes := strings.Split(conf.GlobalConfig.Upload.AllowTypes, ",")
	validType := false
	for _, t := range allowTypes {
		if strings.TrimSpace(t) == ext {
			validType = true
			break
		}
	}
	if !validType {
		return fmt.Errorf("不支持的文件类型：%s，允许的类型：%s", ext, conf.GlobalConfig.Upload.AllowTypes)
	}

	// 验证文件大小
	fileSizeMB := float64(file.Size) / 1024.0 / 1024.0
	maxSizeMB := float64(conf.GlobalConfig.Upload.MaxFileSize)

	if fileSizeMB > maxSizeMB {
		return fmt.Errorf("文件大小超过限制：%.2fMB，最大允许：%dMB",
			fileSizeMB,
			conf.GlobalConfig.Upload.MaxFileSize)
	}

	return nil
}

// 处理文件名，根据配置决定是否添加时间戳
func ProcessFileName(originalName string) string {
	if conf.GlobalConfig.Upload.KeepOriginalName {
		return originalName
	}

	// 获取文件名和扩展名
	ext := filepath.Ext(originalName)
	nameWithoutExt := strings.TrimSuffix(originalName, ext)

	// 添加时间戳
	timestamp := time.Now().Format("20060102150405")

	// 组合新文件名
	return fmt.Sprintf("%s_%s%s", nameWithoutExt, timestamp, ext)
}

// IsImageFile 检查文件是否为图片类型
func IsImageFile(filePath string) bool {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filePath), "."))
	imageExts := []string{"jpg", "jpeg", "png", "gif", "bmp", "webp", "svg", "ico"}

	for _, imgExt := range imageExts {
		if ext == imgExt {
			return true
		}
	}

	return false
}
