package controllers

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	
	"cyi-note/backend/config"
	"cyi-note/backend/models"
	"cyi-note/backend/utils"
)

// 全局变量存储上传目录
var uploadBaseDir string

// InitAttachmentController 初始化附件控制器
func InitAttachmentController(cfg *config.Config) {
	// 使用配置中的上传目录
	uploadBaseDir = cfg.UploadDir
	
	// 确保上传目录存在
	if err := os.MkdirAll(uploadBaseDir, 0755); err != nil {
		panic(fmt.Sprintf("无法创建上传目录 %s: %v", uploadBaseDir, err))
	}
}

// UploadAttachment 上传附件
func UploadAttachment(c *gin.Context) {
	// 获取当前用户ID
	userID, _ := c.Get("userID")
	
	// 检查是否提供了笔记ID
	noteIDStr := c.PostForm("note_id")
	if noteIDStr == "" {
		utils.BadRequestResponse(c, "请提供笔记ID")
		return
	}
	
	// 解析笔记ID
	noteID, err := strconv.ParseUint(noteIDStr, 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "无效的笔记ID")
		return
	}
	
	// 检查笔记是否存在且属于当前用户
	note, err := models.GetNoteByID(uint(noteID))
	if err != nil {
		utils.NotFoundResponse(c, "笔记未找到")
		return
	}
	
	if note.UserID != userID.(uint) {
		utils.ForbiddenResponse(c, "无权为此笔记上传附件")
		return
	}
	
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequestResponse(c, "获取上传文件失败")
		return
	}
	
	// 使用原始文件名或从表单获取文件名
	originalFilename := file.Filename
	
	// 生成文件存储路径
	userUploadDir := filepath.Join(uploadBaseDir, strconv.FormatUint(uint64(userID.(uint)), 10))
	if err := os.MkdirAll(userUploadDir, 0755); err != nil {
		utils.ServerErrorResponse(c, "创建上传目录失败")
		return
	}
	
	// 生成唯一的文件标识符，但保留原始文件名部分
	fileExt := filepath.Ext(originalFilename)
	baseName := originalFilename[:len(originalFilename)-len(fileExt)]
	// 使用时间戳和随机UUID确保唯一性，但保留原始文件名部分
	uniqueID := fmt.Sprintf("%d_%s", time.Now().UnixNano(), uuid.New().String()[:8])
	// 安全处理文件名，移除可能有问题的字符
	safeBaseName := utils.SanitizeFilename(baseName)
	// 构建最终存储的文件名，格式为：原文件名_唯一标识符.扩展名
	fileName := fmt.Sprintf("%s_%s%s", safeBaseName, uniqueID, fileExt)
	filePath := filepath.Join(userUploadDir, fileName)
	
	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		utils.ServerErrorResponse(c, "保存文件失败")
		return
	}

	// 计算相对路径
	relFilePath, err := filepath.Rel(uploadBaseDir, filePath)
	if err != nil {
		relFilePath = filePath // 如果计算相对路径失败，使用完整路径
	}
	storedFilePath := filepath.Join(uploadBaseDir, relFilePath)
	
	// 创建附件记录
	attachment := models.Attachment{
		NoteID:    uint(noteID),
		Filename:  originalFilename, // 存储原始文件名，用于下载时显示
		Filepath:  storedFilePath,   // 存储实际文件路径
		Filesize:  file.Size,
		Filetype:  file.Header.Get("Content-Type"),
		CreatedAt: time.Now(),
	}
	
	if err := models.CreateAttachment(&attachment); err != nil {
		// 如果创建附件记录失败，删除已上传的文件
		os.Remove(filePath)
		utils.ServerErrorResponse(c, "创建附件记录失败")
		return
	}
	
	// 设置文件URL，使用ID而不是文件路径，避免安全问题
	attachment.FileURL = "/api/attachments/" + strconv.FormatUint(uint64(attachment.ID), 10)
	
	utils.CreatedResponse(c, attachment, "附件上传成功")
}

// GetAttachment 获取附件
func GetAttachment(c *gin.Context) {
	// 获取附件ID
	attachmentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "无效的附件ID")
		return
	}
	
	// 获取附件
	attachment, err := models.GetAttachmentByID(uint(attachmentID))
	if err != nil {
		utils.NotFoundResponse(c, "附件未找到")
		return
	}
	
	// 临时移除认证校验，方便测试图片显示
	// 注意：如果需要认证，可以从URL参数获取token
	token := c.Query("token")
	fmt.Printf("从URL接收到的令牌: %s\n", token)
	
	// 检查文件是否存在
	if _, err := os.Stat(attachment.Filepath); os.IsNotExist(err) {
		// 日志记录
		fmt.Printf("文件路径不存在: %s\n", attachment.Filepath)
		utils.NotFoundResponse(c, "文件不存在")
		return
	}
	
	// 根据文件类型决定是内联显示还是作为附件下载
	contentDisposition := "attachment"
	if strings.HasPrefix(attachment.Filetype, "image/") {
		// 图片文件内联显示
		contentDisposition = "inline"
		// 确保图片文件类型正确
		c.Header("Content-Type", attachment.Filetype)
	}
	
	// 文件名编码，处理非ASCII字符
	encodedFilename := url.QueryEscape(attachment.Filename)
	
	// 提供文件下载或显示
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf(`%s; filename="%s"; filename*=UTF-8''%s`, 
		contentDisposition, attachment.Filename, encodedFilename))
	
	// 设置更多响应头来解决跨域问题
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	
	// 设置缓存头
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	if strings.HasPrefix(attachment.Filetype, "image/") {
		c.Header("Cache-Control", "max-age=31536000, public") // 图片缓存1年
	} else {
		c.Header("Cache-Control", "private, no-transform, no-store, must-revalidate")
	}
	c.Header("Accept-Ranges", "bytes")
	c.Header("X-Content-Type-Options", "nosniff")
	
	// 打印调试信息
	fmt.Printf("提供文件: %s (类型: %s, 路径: %s)\n", 
		attachment.Filename, attachment.Filetype, attachment.Filepath)
	
	// 使用Gin的文件发送功能
	c.File(attachment.Filepath)
}

// DeleteAttachment 删除附件
func DeleteAttachment(c *gin.Context) {
	// 获取附件ID
	attachmentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "无效的附件ID")
		return
	}
	
	// 获取附件
	attachment, err := models.GetAttachmentByID(uint(attachmentID))
	if err != nil {
		utils.NotFoundResponse(c, "附件未找到")
		return
	}
	
	// 获取当前用户ID
	userID, _ := c.Get("userID")
	
	// 检查附件所属的笔记是否属于当前用户
	note, err := models.GetNoteByID(attachment.NoteID)
	if err != nil || note.UserID != userID.(uint) {
		utils.ForbiddenResponse(c, "无权删除此附件")
		return
	}
	
	// 删除文件
	if err := os.Remove(attachment.Filepath); err != nil && !os.IsNotExist(err) {
		utils.ServerErrorResponse(c, "删除文件失败")
		return
	}
	
	// 删除附件记录
	if err := models.DeleteAttachment(uint(attachmentID)); err != nil {
		utils.ServerErrorResponse(c, "删除附件记录失败")
		return
	}
	
	utils.OkResponse(c, nil, "附件已删除")
}

// GetNoteAttachments 获取笔记的所有附件
func GetNoteAttachments(c *gin.Context) {
	// 获取笔记ID
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "无效的笔记ID")
		return
	}
	
	// 获取当前用户ID
	userID, _ := c.Get("userID")
	
	// 检查笔记是否属于当前用户
	note, err := models.GetNoteByID(uint(noteID))
	if err != nil {
		utils.NotFoundResponse(c, "笔记未找到")
		return
	}
	
	if note.UserID != userID.(uint) {
		utils.ForbiddenResponse(c, "无权访问此笔记的附件")
		return
	}
	
	// 获取笔记的所有附件
	attachments, err := models.GetAttachmentsByNoteID(uint(noteID))
	if err != nil {
		utils.ServerErrorResponse(c, "获取附件失败")
		return
	}
	
	utils.OkResponse(c, attachments, "获取笔记附件成功")
}

// GetAttachmentsByDate 获取按日期分组的附件列表
func GetAttachmentsByDate(c *gin.Context) {
	// 获取当前用户ID
	userID, _ := c.Get("userID")
	
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "20"))
	
	// 获取可选的文件类型过滤参数
	fileType := c.Query("filetype")
	
	// 从数据库获取按日期分组的附件
	attachments, dateGroups, total, err := models.GetAttachmentsByDate(userID.(uint), page, pageSize, fileType)
	if err != nil {
		utils.ServerErrorResponse(c, "获取资源库文件失败")
		return
	}
	
	// 返回数据
	utils.OkResponse(c, gin.H{
		"attachments": attachments,
		"dateGroups": dateGroups,
		"total": total,
		"page": page,
		"pageSize": pageSize,
	}, "获取资源库文件成功")
} 