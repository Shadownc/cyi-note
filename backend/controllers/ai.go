package controllers

import (
	"strconv"
	
	"github.com/gin-gonic/gin"
	
	"cyi-note/backend/models"
	"cyi-note/backend/utils"
)

// AI标签生成请求
type GenerateTagsRequest struct {
	Content string `json:"content" binding:"required"`
}

// AI摘要生成请求
type GenerateSummaryRequest struct {
	Content string `json:"content" binding:"required"`
}

// 生成标签建议请求
type GenerateTagSuggestionsRequest struct {
	Content string `json:"content" binding:"required"`
}

// GenerateTags 生成标签
func GenerateTags(c *gin.Context) {
	var req GenerateTagsRequest
	
	// 验证请求
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "无效的请求参数")
		return
	}
	
	// 使用AI服务提取关键词
	tags, err := utils.ExtractKeywords(req.Content)
	if err != nil {
		utils.ServerErrorResponse(c, "生成标签失败")
		return
	}
	
	// 如果内容太短或无法提取关键词，返回空标签
	if len(tags) == 0 {
		utils.OkResponse(c, gin.H{
			"tags": []string{},
		}, "无法从内容中提取标签")
		return
	}
	
	utils.OkResponse(c, gin.H{
		"tags": tags,
	}, "标签生成成功")
}

// GenerateTagSuggestions 生成标签建议
func GenerateTagSuggestions(c *gin.Context) {
	var req GenerateTagSuggestionsRequest
	
	// 验证请求
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "无效的请求参数")
		return
	}
	
	// 生成标签建议
	tags, err := utils.GenerateTagSuggestions(req.Content)
	if err != nil {
		utils.ServerErrorResponse(c, "生成标签建议失败")
		return
	}
	
	utils.OkResponse(c, gin.H{
		"tags": tags,
	}, "标签建议生成成功")
}

// GenerateSummary 生成内容摘要
func GenerateSummary(c *gin.Context) {
	var req GenerateSummaryRequest
	
	// 验证请求
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "无效的请求参数")
		return
	}
	
	// 检查是否提供了笔记ID
	noteIDStr := c.Query("note_id")
	if noteIDStr == "" {
		// 如果没有提供笔记ID，只返回生成的摘要
		summary, err := utils.GenerateSummary(req.Content)
		if err != nil {
			utils.ServerErrorResponse(c, "生成摘要失败")
			return
		}
		
		utils.OkResponse(c, gin.H{
			"summary": summary,
		}, "摘要生成成功")
		return
	}
	
	// 如果提供了笔记ID，更新笔记的摘要
	noteID, err := strconv.ParseUint(noteIDStr, 10, 64)
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
	
	// 生成摘要
	summary, err := utils.GenerateSummary(req.Content)
	if err != nil {
		utils.ServerErrorResponse(c, "生成摘要失败")
		return
	}
	
	// 更新笔记摘要
	note.Summary = summary
	if err := models.UpdateNote(note); err != nil {
		utils.ServerErrorResponse(c, "更新笔记摘要失败")
		return
	}
	
	utils.OkResponse(c, gin.H{
		"summary": summary,
	}, "摘要生成并更新成功")
} 