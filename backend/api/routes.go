package api

import (
	"github.com/gin-gonic/gin"
	
	"cyi-note/backend/controllers"
	"cyi-note/backend/middleware"
)

// SetupRoutes 设置API路由
func SetupRoutes(r *gin.Engine) {
	// 静态文件服务
	r.Static("/uploads", "./uploads")
	
	// 创建API路由组
	api := r.Group("/api")
	
	// 认证相关路由
	auth := api.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.GET("/user", middleware.AuthRequired(), controllers.GetCurrentUser)
		auth.PUT("/user", middleware.AuthRequired(), controllers.UpdateUser)
	}
	
	// 笔记相关路由
	notes := api.Group("/notes", middleware.AuthRequired())
	{
		notes.GET("", controllers.GetNotes)
		notes.POST("", controllers.CreateNote)
		notes.GET("/:id", controllers.GetNote)
		notes.PUT("/:id", controllers.UpdateNote)
		notes.DELETE("/:id", controllers.DeleteNote)
		notes.GET("/search", controllers.SearchNotes)
	}
	
	// 标签相关路由
	tags := api.Group("/tags", middleware.AuthRequired())
	{
		tags.GET("", controllers.GetTags)
		tags.POST("", controllers.CreateTag)
		tags.DELETE("/:id", controllers.DeleteTag)
		
		// 标签操作
		tags.POST("/:id/notes/:noteId", controllers.AddTagToNote)
		tags.DELETE("/:id/notes/:noteId", controllers.RemoveTagFromNote)
	}
	
	// 附件相关路由
	attachments := api.Group("/attachments", middleware.AuthRequired())
	{
		attachments.POST("", controllers.UploadAttachment)
		attachments.GET("/:id", controllers.GetAttachment)
		attachments.DELETE("/:id", controllers.DeleteAttachment)
	}
	
	// AI相关路由
	ai := api.Group("/ai", middleware.AuthRequired())
	{
		ai.POST("/tags", controllers.GenerateTags)
		ai.POST("/summary", controllers.GenerateSummary)
	}
} 