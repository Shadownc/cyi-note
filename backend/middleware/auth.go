package middleware

import (
	"errors"
	"strings"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	
	"cyi-note/backend/config"
	"cyi-note/backend/models"
	"cyi-note/backend/utils"
)

// JWTClaims 自定义JWT Claims
type JWTClaims struct {
	UserID uint `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// AuthRequired 认证中间件，验证JWT令牌
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Bearer Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c)
			c.Abort()
			return
		}
		
		// 检查是否为Bearer Token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.UnauthorizedResponse(c)
			c.Abort()
			return
		}
		tokenString := parts[1]
		
		// 解析并验证Token
		claims, err := ParseToken(tokenString)
		if err != nil {
			utils.UnauthorizedResponse(c)
			c.Abort()
			return
		}
		
		// 将用户信息存入上下文
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		
		c.Next()
	}
}

// AdminRequired 管理员认证中间件
func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先执行基本认证
		AuthRequired()(c)
		
		// 检查是否已终止请求处理
		if c.IsAborted() {
			return
		}
		
		// 检查角色是否为管理员
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			utils.ForbiddenResponse(c, "需要管理员权限")
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// GenerateToken 生成JWT令牌
func GenerateToken(user *models.User) (string, error) {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		return "", err
	}
	
	// 设置过期时间
	expiresAt := time.Now().Add(time.Hour * time.Duration(cfg.JWTExpiry))
	
	// 创建Claims
	claims := JWTClaims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	
	// 生成Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

// ParseToken 解析JWT令牌
func ParseToken(tokenString string) (*JWTClaims, error) {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}
	
	// 解析Token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.JWTSecret), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	// 验证Token
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, errors.New("无效的令牌")
} 