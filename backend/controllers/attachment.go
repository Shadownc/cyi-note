package controllers

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	
	"cyi-note/backend/models"
	"cyi-note/backend/utils"
)

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
	
	// 生成文件存储路径
	uploadDir := filepath.Join("uploads", strconv.FormatUint(uint64(userID.(uint)), 10))
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		utils.ServerErrorResponse(c, "创建上传目录失败")
		return
	}
	
	// 生成唯一的文件名
	fileExt := filepath.Ext(file.Filename)
	fileName := uuid.New().String() + fileExt
	filePath := filepath.Join(uploadDir, fileName)
	
	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		utils.ServerErrorResponse(c, "保存文件失败")
		return
	}
	
	// 创建附件记录
	attachment := models.Attachment{
		NoteID:    uint(noteID),
		Filename:  file.Filename,
		Filepath:  filePath,
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
	
	// 设置文件URL
	attachment.FileURL = "/api/attachments/" + attachment.Filepath
	
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
	
	// 获取当前用户ID
	userID, _ := c.Get("userID")
	
	// 检查附件所属的笔记是否属于当前用户
	note, err := models.GetNoteByID(attachment.NoteID)
	if err != nil || note.UserID != userID.(uint) {
		utils.ForbiddenResponse(c, "无权访问此附件")
		return
	}
	
	// 文件不存在则返回错误
	if _, err := os.Stat(attachment.Filepath); os.IsNotExist(err) {
		utils.NotFoundResponse(c, "文件不存在")
		return
	}
	
	// 提供文件下载
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", attachment.Filename))
	c.Header("Content-Type", attachment.Filetype)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
	c.Header("Cache-Control", "must-revalidate")
	c.Header("Pragma", "public")
	
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