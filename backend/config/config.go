package config

import (
	"fmt"
	"os"
	"strconv"
	
	"github.com/joho/godotenv"
)

// Config 应用程序配置
type Config struct {
	// 服务器配置
	ServerHost string
	ServerPort int
	ClientURL  string
	
	// 数据库配置
	Database DatabaseConfig
	
	// JWT配置
	JWTSecret string
	JWTExpiry int // 过期时间（小时）
	
	// 管理员配置
	AdminUsername string
	AdminPassword string
	AdminEmail    string
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type     string // mysql, postgres, sqlite
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// DSN 返回数据库连接字符串
func (db *DatabaseConfig) DSN() string {
	switch db.Type {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			db.Username, db.Password, db.Host, db.Port, db.DBName)
	case "postgres":
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			db.Host, db.Port, db.Username, db.Password, db.DBName, db.SSLMode)
	case "sqlite":
		return db.DBName
	default:
		return ""
	}
}

// LoadConfig 从环境变量或.env文件加载配置
func LoadConfig() (*Config, error) {
	// 尝试加载.env文件（如果存在）
	_ = godotenv.Load()
	
	// 服务器配置
	serverHost := getEnv("SERVER_HOST", "0.0.0.0")
	serverPort, _ := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	clientURL := getEnv("CLIENT_URL", "http://localhost:3000")
	
	// 数据库配置
	dbType := getEnv("DB_TYPE", "mysql")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "3306"))
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASS", "")
	dbName := getEnv("DB_NAME", "cyi_note")
	dbSSLMode := getEnv("DB_SSL_MODE", "disable")
	
	// JWT配置
	jwtSecret := getEnv("JWT_SECRET", "your-secret-key")
	jwtExpiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY", "24"))
	
	// 管理员配置
	adminUsername := getEnv("ADMIN_USERNAME", "admin")
	adminPassword := getEnv("ADMIN_PASSWORD", "admin123")
	adminEmail := getEnv("ADMIN_EMAIL", "admin@example.com")
	
	return &Config{
		ServerHost: serverHost,
		ServerPort: serverPort,
		ClientURL:  clientURL,
		
		Database: DatabaseConfig{
			Type:     dbType,
			Host:     dbHost,
			Port:     dbPort,
			Username: dbUser,
			Password: dbPass,
			DBName:   dbName,
			SSLMode:  dbSSLMode,
		},
		
		JWTSecret: jwtSecret,
		JWTExpiry: jwtExpiry,
		
		AdminUsername: adminUsername,
		AdminPassword: adminPassword,
		AdminEmail:    adminEmail,
	}, nil
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
} 