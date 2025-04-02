package controllers

import (
	"strconv"
	
	"github.com/gin-gonic/gin"
	
	"cyi-note/backend/models"
	"cyi-note/backend/utils"
)

// 笔记请求
type NoteRequest struct {
	Title    string   `json:"title" binding:"required"`
	Content  string   `json:"content"`
	Tags     []string `json:"tags"` // 标签名称列表
	IsPublic bool     `json:"is_public"` // 是否公开笔记
}

// CreateNote 创建笔记
func CreateNote(c *gin.Context) {
	var req NoteRequest
	
	// 验证请求
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "无效的请求参数")
		return
	}
	
	// 获取当前用户ID
	userID, _ := c.Get("userID")
	
	// 创建笔记
	note := models.Note{
		UserID:   userID.(uint),
		Title:    req.Title,
		Content:  req.Content,
		IsPublic: req.IsPublic, // 设置是否公开
	}
	
	// 保存笔记
	if err := models.CreateNote(&note); err != nil {
		utils.ServerErrorResponse(c, "创建笔记失败")
		return
	}
	
	// 处理标签
	if len(req.Tags) > 0 {
		for _, tagName := range req.Tags {
			tag, err := models.GetOrCreateTag(tagName)
			if err != nil {
				utils.ServerErrorResponse(c, "处理标签失败")
				return
			}
			
			// 添加标签到笔记
			if err := models.AddTagToNote(note.ID, tag.ID); err != nil {
				utils.ServerErrorResponse(c, "添加标签失败")
				return
			}
		}
	}
	
	// 查询带有标签的笔记
	createdNote, err := models.GetNoteByID(note.ID)
	if err != nil {
		utils.ServerErrorResponse(c, "获取创建的笔记失败")
		return
	}
	
	utils.CreatedResponse(c, createdNote, "笔记创建成功")
}

// GetNote 获取笔记详情
func GetNote(c *gin.Context) {
	// 获取笔记ID
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "无效的笔记ID")
		return
	}
	
	// 获取笔记
	note, err := models.GetNoteByID(uint(noteID))
	if err != nil {
		utils.NotFoundResponse(c, "笔记未找到")
		return
	}
	
	// 检查笔记所有权
	userID, _ := c.Get("userID")
	if note.UserID != userID.(uint) {
		utils.ForbiddenResponse(c, "无权访问此笔记")
		return
	}
	
	utils.OkResponse(c, note, "获取笔记成功")
}

// GetNotes 获取用户的所有笔记
func GetNotes(c *gin.Context) {
	// 获取当前用户ID
	userID, _ := c.Get("userID")
	
	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	
	// 标签过滤
	tagID := c.Query("tag")
	
	var notes []models.Note
	var total int64
	var err error
	
	if tagID != "" {
		// 如果提供了标签ID，按标签过滤
		tagIDUint, _ := strconv.ParseUint(tagID, 10, 64)
		notes, total, err = models.SearchNotesByTag(userID.(uint), uint(tagIDUint), page, pageSize)
	} else {
		// 否则获取所有笔记
		notes, total, err = models.GetNotesByUserID(userID.(uint), page, pageSize)
	}
	
	if err != nil {
		utils.ServerErrorResponse(c, "获取笔记失败")
		return
	}
	
	utils.OkResponse(c, gin.H{
		"notes": notes,
		"total": total,
		"page":  page,
		"size":  pageSize,
	}, "获取笔记列表成功")
}

// UpdateNote 更新笔记
func UpdateNote(c *gin.Context) {
	var req NoteRequest
	
	// 验证请求
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "无效的请求参数")
		return
	}
	
	// 获取笔记ID
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "无效的笔记ID")
		return
	}
	
	// 获取当前用户ID
	userID, _ := c.Get("userID")
	
	// 获取笔记
	note, err := models.GetNoteByID(uint(noteID))
	if err != nil {
		utils.NotFoundResponse(c, "笔记未找到")
		return
	}
	
	// 检查笔记所有权
	if note.UserID != userID.(uint) {
		utils.ForbiddenResponse(c, "无权修改此笔记")
		return
	}
	
	// 更新笔记
	note.Title = req.Title
	note.Content = req.Content
	note.IsPublic = req.IsPublic
	
	if err := models.UpdateNote(note); err != nil {
		utils.ServerErrorResponse(c, "更新笔记失败")
		return
	}
	
	// 处理标签
	// 先清除所有标签关联
	if len(note.Tags) > 0 {
		for _, tag := range note.Tags {
			if err := models.RemoveTagFromNote(note.ID, tag.ID); err != nil {
				utils.ServerErrorResponse(c, "更新标签失败")
				return
			}
		}
	}
	
	// 添加新标签
	if len(req.Tags) > 0 {
		for _, tagName := range req.Tags {
			tag, err := models.GetOrCreateTag(tagName)
			if err != nil {
				utils.ServerErrorResponse(c, "处理标签失败")
				return
			}
			
			// 添加标签到笔记
			if err := models.AddTagToNote(note.ID, tag.ID); err != nil {
				utils.ServerErrorResponse(c, "添加标签失败")
				return
			}
		}
	}
	
	// 获取更新后的笔记
	updatedNote, err := models.GetNoteByID(note.ID)
	if err != nil {
		utils.ServerErrorResponse(c, "获取更新后的笔记失败")
		return
	}
	
	utils.OkResponse(c, updatedNote, "笔记更新成功")
}

// DeleteNote 删除笔记
func DeleteNote(c *gin.Context) {
	// 获取笔记ID
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "无效的笔记ID")
		return
	}
	
	// 获取当前用户ID
	userID, _ := c.Get("userID")
	
	// 获取笔记
	note, err := models.GetNoteByID(uint(noteID))
	if err != nil {
		utils.NotFoundResponse(c, "笔记未找到")
		return
	}
	
	// 检查笔记所有权
	if note.UserID != userID.(uint) {
		utils.ForbiddenResponse(c, "无权删除此笔记")
		return
	}
	
	// 删除笔记
	if err := models.DeleteNote(uint(noteID)); err != nil {
		utils.ServerErrorResponse(c, "删除笔记失败")
		return
	}
	
	utils.OkResponse(c, nil, "笔记已删除")
}

// SearchNotes 搜索笔记
func SearchNotes(c *gin.Context) {
	// 获取当前用户ID
	userID, _ := c.Get("userID")
	
	// 获取搜索关键词
	keyword := c.Query("keyword")
	if keyword == "" {
		utils.BadRequestResponse(c, "请提供搜索关键词")
		return
	}
	
	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	
	// 搜索笔记
	notes, total, err := models.SearchNotes(userID.(uint), keyword, page, pageSize)
	if err != nil {
		utils.ServerErrorResponse(c, "搜索笔记失败")
		return
	}
	
	utils.OkResponse(c, gin.H{
		"notes": notes,
		"total": total,
		"page":  page,
		"size":  pageSize,
	}, "搜索笔记成功")
}

// GetPublicNotes 获取所有公开的笔记
func GetPublicNotes(c *gin.Context) {
	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	
	// 获取公开笔记
	notes, total, err := models.GetPublicNotes(page, pageSize)
	if err != nil {
		utils.ServerErrorResponse(c, "获取公开笔记失败")
		return
	}
	
	utils.OkResponse(c, gin.H{
		"notes": notes,
		"total": total,
		"page":  page,
		"size":  pageSize,
	}, "获取公开笔记成功")
} 