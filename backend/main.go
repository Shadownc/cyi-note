package main

import (
	"fmt"
	"log"
	
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	
	"cyi-note/backend/config"
	"cyi-note/backend/models"
	"cyi-note/backend/api"
	"cyi-note/backend/controllers"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("无法加载配置: %v", err)
	}
	
	// 初始化数据库
	if err := models.InitDB(cfg.Database); err != nil {
		log.Fatalf("无法初始化数据库: %v", err)
	}
	
	// 确保管理员账号存在
	if err := models.EnsureAdminExists(cfg); err != nil {
		log.Fatalf("检查管理员账号失败: %v", err)
	}
	
	// 初始化附件控制器
	controllers.InitAttachmentController(cfg)
	
	// 创建Gin引擎
	r := gin.Default()
	
	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           86400,
	}))
	
	// 初始化路由
	api.SetupRoutes(r)
	
	// 启动服务器
	serverAddr := fmt.Sprintf("%s:%d", cfg.ServerHost, cfg.ServerPort)
	log.Printf("服务器开始监听: %s\n", serverAddr)
	
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
} 