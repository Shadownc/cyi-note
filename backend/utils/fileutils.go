package utils

import (
	"regexp"
	"strings"
)

// SanitizeFilename 清理文件名，移除不安全或不合法的字符
func SanitizeFilename(filename string) string {
	// 移除路径分隔符和控制字符
	re := regexp.MustCompile(`[\\/:*?"<>|]`)
	sanitized := re.ReplaceAllString(filename, "_")
	
	// 移除多余的空格
	sanitized = strings.TrimSpace(sanitized)
	
	// 限制文件名长度，避免过长的文件名
	if len(sanitized) > 100 {
		sanitized = sanitized[:100]
	}
	
	// 如果清理后为空，使用默认名称
	if sanitized == "" {
		sanitized = "file"
	}
	
	return sanitized
}

// GetFileExtension 获取文件扩展名（带点）
func GetFileExtension(filename string) string {
	idx := strings.LastIndex(filename, ".")
	if idx == -1 {
		return ""
	}
	return filename[idx:]
}

// IsImageFile 检查文件是否为图像文件
func IsImageFile(contentType string) bool {
	return strings.HasPrefix(contentType, "image/")
}

// IsDocumentFile 检查文件是否为文档文件
func IsDocumentFile(contentType string) bool {
	documentTypes := []string{
		"application/pdf",
		"application/msword",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.ms-excel",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/vnd.ms-powerpoint",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"text/plain",
		"text/csv",
	}
	
	for _, docType := range documentTypes {
		if contentType == docType {
			return true
		}
	}
	
	return false
} 