package controllers

import (
	"strconv"
	
	"github.com/gin-gonic/gin"
	
	"cyi-note/backend/models"
	"cyi-note/backend/utils"
)

// 标签请求
type TagRequest struct {
	Name string `json:"name" binding:"required"`
}

// GetTags 获取所有标签
func GetTags(c *gin.Context) {
	// 获取当前用户ID
	userID, _ := c.Get("userID")
	
	// 获取带有笔记数量的标签列表
	tagsWithCount, err := models.GetAllTagsWithCount(userID.(uint))
	if err != nil {
		utils.ServerErrorResponse(c, "获取标签失败")
		return
	}
	
	utils.OkResponse(c, tagsWithCount, "获取标签成功")
}

// CreateTag 创建标签
func CreateTag(c *gin.Context) {
	var req TagRequest
	
	// 验证请求
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "无效的请求参数")
		return
	}
	
	// 检查标签是否存在
	existingTag, err := models.GetTagByName(req.Name)
	if err == nil && existingTag != nil {
		// 标签已存在，直接返回
		utils.OkResponse(c, existingTag, "标签已存在")
		return
	}
	
	// 创建新标签
	tag := models.Tag{
		Name: req.Name,
	}
	
	if err := models.CreateTag(&tag); err != nil {
		utils.ServerErrorResponse(c, "创建标签失败")
		return
	}
	
	utils.CreatedResponse(c, tag, "标签创建成功")
}

// DeleteTag 删除标签
func DeleteTag(c *gin.Context) {
	// 获取标签ID
	tagID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "无效的标签ID")
		return
	}
	
	// 检查标签是否存在
	_, err = models.GetTagByID(uint(tagID))
	if err != nil {
		utils.NotFoundResponse(c, "标签未找到")
		return
	}
	
	// 删除标签
	if err := models.DeleteTag(uint(tagID)); err != nil {
		utils.ServerErrorResponse(c, "删除标签失败")
		return
	}
	
	utils.OkResponse(c, nil, "标签已删除")
}

// AddTagToNote 将标签添加到笔记
func AddTagToNote(c *gin.Context) {
	// 获取标签ID和笔记ID
	tagID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "无效的标签ID")
		return
	}
	
	noteID, err := strconv.ParseUint(c.Param("noteId"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "无效的笔记ID")
		return
	}
	
	// 获取当前用户ID
	userID, _ := c.Get("userID")
	
	// 获取笔记，检查笔记所有权
	note, err := models.GetNoteByID(uint(noteID))
	if err != nil {
		utils.NotFoundResponse(c, "笔记未找到")
		return
	}
	
	if note.UserID != userID.(uint) {
		utils.ForbiddenResponse(c, "无权修改此笔记")
		return
	}
	
	// 获取标签
	tag, err := models.GetTagByID(uint(tagID))
	if err != nil {
		utils.NotFoundResponse(c, "标签未找到")
		return
	}
	
	// 添加标签到笔记
	if err := models.AddTagToNote(uint(noteID), uint(tagID)); err != nil {
		utils.ServerErrorResponse(c, "添加标签失败")
		return
	}
	
	utils.OkResponse(c, tag, "标签已添加到笔记")
}

// RemoveTagFromNote 从笔记中移除标签
func RemoveTagFromNote(c *gin.Context) {
	// 获取标签ID和笔记ID
	tagID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "无效的标签ID")
		return
	}
	
	noteID, err := strconv.ParseUint(c.Param("noteId"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "无效的笔记ID")
		return
	}
	
	// 获取当前用户ID
	userID, _ := c.Get("userID")
	
	// 获取笔记，检查笔记所有权
	note, err := models.GetNoteByID(uint(noteID))
	if err != nil {
		utils.NotFoundResponse(c, "笔记未找到")
		return
	}
	
	if note.UserID != userID.(uint) {
		utils.ForbiddenResponse(c, "无权修改此笔记")
		return
	}
	
	// 从笔记中移除标签
	if err := models.RemoveTagFromNote(uint(noteID), uint(tagID)); err != nil {
		utils.ServerErrorResponse(c, "移除标签失败")
		return
	}
	
	utils.OkResponse(c, nil, "标签已从笔记中移除")
} 