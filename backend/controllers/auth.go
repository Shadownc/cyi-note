package controllers

import (
	"github.com/gin-gonic/gin"
	
	"cyi-note/backend/middleware"
	"cyi-note/backend/models"
	"cyi-note/backend/utils"
)

// 用户注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

// 用户登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 用户信息更新请求
type UpdateUserRequest struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

// Register 用户注册
func Register(c *gin.Context) {
	var req RegisterRequest
	
	// 验证请求
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "无效的请求参数")
		return
	}
	
	// 检查用户名是否已存在
	_, err := models.GetUserByUsername(req.Username)
	if err == nil {
		utils.BadRequestResponse(c, "用户名已存在")
		return
	}
	
	// 检查邮箱是否已存在
	_, err = models.GetUserByEmail(req.Email)
	if err == nil {
		utils.BadRequestResponse(c, "邮箱已存在")
		return
	}
	
	// 创建用户
	user := models.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Role:     "user", // 默认为普通用户
	}
	
	if err := models.CreateUser(&user); err != nil {
		utils.ServerErrorResponse(c, "创建用户失败")
		return
	}
	
	// 生成JWT令牌
	token, err := middleware.GenerateToken(&user)
	if err != nil {
		utils.ServerErrorResponse(c, "生成令牌失败")
		return
	}
	
	// 返回用户信息和令牌
	utils.CreatedResponse(c, gin.H{
		"user":  user,
		"token": token,
	}, "注册成功")
}

// Login 用户登录
func Login(c *gin.Context) {
	var req LoginRequest
	
	// 验证请求
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "无效的请求参数")
		return
	}
	
	// 查找用户
	user, err := models.GetUserByUsername(req.Username)
	if err != nil {
		utils.UnauthorizedResponse(c)
		return
	}
	
	// 验证密码
	if !user.CheckPassword(req.Password) {
		utils.UnauthorizedResponse(c)
		return
	}
	
	// 生成JWT令牌
	token, err := middleware.GenerateToken(user)
	if err != nil {
		utils.ServerErrorResponse(c, "生成令牌失败")
		return
	}
	
	// 返回用户信息和令牌
	utils.OkResponse(c, gin.H{
		"user":  user,
		"token": token,
	}, "登录成功")
}

// GetCurrentUser 获取当前用户信息
func GetCurrentUser(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.UnauthorizedResponse(c)
		return
	}
	
	// 查找用户
	user, err := models.GetUserByID(userID.(uint))
	if err != nil {
		utils.ServerErrorResponse(c, "获取用户信息失败")
		return
	}
	
	// 返回用户信息
	utils.OkResponse(c, user, "获取用户信息成功")
}

// UpdateUser 更新用户信息
func UpdateUser(c *gin.Context) {
	// 从上下文中获取用户ID
	userID, exists := c.Get("userID")
	if !exists {
		utils.UnauthorizedResponse(c)
		return
	}
	
	// 解析请求
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "无效的请求参数")
		return
	}
	
	// 查找用户
	user, err := models.GetUserByID(userID.(uint))
	if err != nil {
		utils.ServerErrorResponse(c, "获取用户信息失败")
		return
	}
	
	// 标记是否修改了密码
	passwordChanged := false
	
	// 更新基本信息（如果提供）
	if req.Username != "" && req.Username != user.Username {
		// 检查新用户名是否已被使用
		existingUser, err := models.GetUserByUsername(req.Username)
		if err == nil && existingUser.ID != user.ID {
			utils.BadRequestResponse(c, "用户名已存在")
			return
		}
		user.Username = req.Username
	}
	
	if req.Email != "" && req.Email != user.Email {
		// 检查新邮箱是否已被使用
		existingUser, err := models.GetUserByEmail(req.Email)
		if err == nil && existingUser.ID != user.ID {
			utils.BadRequestResponse(c, "邮箱已存在")
			return
		}
		user.Email = req.Email
	}
	
	// 处理密码修改（如果提供）
	if req.CurrentPassword != "" && req.NewPassword != "" {
		// 验证当前密码
		if !user.CheckPassword(req.CurrentPassword) {
			utils.BadRequestResponse(c, "当前密码不正确")
			return
		}
		
		// 设置新密码
		user.Password = req.NewPassword
		passwordChanged = true
	}
	
	// 保存更新
	if err := models.UpdateUser(user); err != nil {
		utils.ServerErrorResponse(c, "更新用户信息失败")
		return
	}
	
	// 返回更新后的用户信息
	message := "用户信息已更新"
	if passwordChanged {
		message = "密码已更新，请重新登录"
	}
	
	utils.OkResponse(c, gin.H{
		"user": user,
		"passwordChanged": passwordChanged,
	}, message)
} 