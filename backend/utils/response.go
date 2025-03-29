package utils

import (
	"net/http"
	
	"github.com/gin-gonic/gin"
)

// Response 标准API响应结构
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse 返回成功响应
func SuccessResponse(c *gin.Context, status int, data interface{}, message string) {
	c.JSON(status, Response{
		Success: true,
		Data:    data,
		Message: message,
	})
}

// ErrorResponse 返回错误响应
func ErrorResponse(c *gin.Context, status int, err string) {
	c.JSON(status, Response{
		Success: false,
		Error:   err,
	})
}

// OkResponse 返回200成功响应
func OkResponse(c *gin.Context, data interface{}, message string) {
	SuccessResponse(c, http.StatusOK, data, message)
}

// CreatedResponse 返回201创建成功响应
func CreatedResponse(c *gin.Context, data interface{}, message string) {
	SuccessResponse(c, http.StatusCreated, data, message)
}

// BadRequestResponse 返回400错误响应
func BadRequestResponse(c *gin.Context, err string) {
	ErrorResponse(c, http.StatusBadRequest, err)
}

// UnauthorizedResponse 返回401未授权响应
func UnauthorizedResponse(c *gin.Context) {
	ErrorResponse(c, http.StatusUnauthorized, "未授权访问")
}

// ForbiddenResponse 返回403禁止访问响应
func ForbiddenResponse(c *gin.Context, err string) {
	ErrorResponse(c, http.StatusForbidden, err)
}

// NotFoundResponse 返回404资源未找到响应
func NotFoundResponse(c *gin.Context, err string) {
	ErrorResponse(c, http.StatusNotFound, err)
}

// ServerErrorResponse 返回500服务器错误响应
func ServerErrorResponse(c *gin.Context, err string) {
	ErrorResponse(c, http.StatusInternalServerError, err)
} 