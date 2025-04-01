package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	
	"cyi-note/backend/models"
	"cyi-note/backend/utils"
)

// 用户列表请求参数
type UserListParams struct {
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"pageSize" binding:"min=1,max=100"`
	Keyword  string `form:"keyword"`
}

// CreateUserRequest 管理员创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
	Role     string `json:"role" binding:"required,oneof=admin user"`
}

// UpdateUserRoleRequest 更新用户角色请求
type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=admin user"`
}

// GetUsers 获取所有用户（分页）
func GetUsers(c *gin.Context) {
	var params UserListParams
	
	// 设置默认值
	params.Page = 1
	params.PageSize = 10
	
	// 绑定参数
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.BadRequestResponse(c, "无效的请求参数")
		return
	}
	
	// 获取用户列表
	users, total, err := models.GetUsers(params.Page, params.PageSize, params.Keyword)
	if err != nil {
		utils.ServerErrorResponse(c, "获取用户列表失败")
		return
	}
	
	// 返回数据
	utils.OkResponse(c, gin.H{
		"users": users,
		"pagination": gin.H{
			"page":      params.Page,
			"pageSize":  params.PageSize,
			"total":     total,
			"totalPage": (total + int64(params.PageSize) - 1) / int64(params.PageSize),
		},
	}, "获取用户列表成功")
}

// GetUser 获取用户详情
func GetUser(c *gin.Context) {
	// 获取用户ID
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "无效的用户ID")
		return
	}
	
	// 获取用户
	user, err := models.GetUserByID(uint(userID))
	if err != nil {
		utils.NotFoundResponse(c, "用户不存在")
		return
	}
	
	// 获取用户笔记总数
	noteCount, err := models.GetUserNoteCount(uint(userID))
	if err != nil {
		utils.ServerErrorResponse(c, "获取用户笔记数量失败")
		return
	}
	
	// 返回数据
	utils.OkResponse(c, gin.H{
		"user":      user,
		"noteCount": noteCount,
	}, "获取用户详情成功")
}

// CreateUser 管理员创建用户
func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	
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
		Role:     req.Role,
	}
	
	if err := models.CreateUser(&user); err != nil {
		utils.ServerErrorResponse(c, "创建用户失败")
		return
	}
	
	// 返回用户信息
	utils.CreatedResponse(c, user, "创建用户成功")
}

// UpdateUserRole 管理员更新用户角色
func UpdateUserRole(c *gin.Context) {
	// 获取用户ID
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "无效的用户ID")
		return
	}
	
	var req UpdateUserRoleRequest
	
	// 验证请求
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequestResponse(c, "无效的请求参数")
		return
	}
	
	// 获取用户
	user, err := models.GetUserByID(uint(userID))
	if err != nil {
		utils.NotFoundResponse(c, "用户不存在")
		return
	}
	
	// 更新角色
	user.Role = req.Role
	
	if err := models.UpdateUser(user); err != nil {
		utils.ServerErrorResponse(c, "更新用户角色失败")
		return
	}
	
	// 返回更新后的用户信息
	utils.OkResponse(c, user, "更新用户角色成功")
}

// DeleteUser 管理员删除用户
func DeleteUserByAdmin(c *gin.Context) {
	// 获取用户ID
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequestResponse(c, "无效的用户ID")
		return
	}
	
	// 获取当前管理员ID
	adminID, _ := c.Get("userID")
	
	// 管理员不能删除自己
	if uint(userID) == adminID.(uint) {
		utils.BadRequestResponse(c, "管理员不能删除自己的账号")
		return
	}
	
	// 获取要删除的用户以确保存在
	_, err = models.GetUserByID(uint(userID))
	if err != nil {
		utils.NotFoundResponse(c, "用户不存在")
		return
	}
	
	// 删除用户
	if err := models.DeleteUser(uint(userID)); err != nil {
		utils.ServerErrorResponse(c, "删除用户失败")
		return
	}
	
	// 返回成功消息
	utils.OkResponse(c, nil, "用户已删除")
} 