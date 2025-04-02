package models

import (
	"log"
	"time"
	
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	
	"cyi-note/backend/config"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB(cfg config.DatabaseConfig) error {
	var db *gorm.DB
	var err error
	
	// 配置GORM日志
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}
	
	// 根据配置选择数据库类型
	switch cfg.Type {
	case "mysql":
		db, err = gorm.Open(mysql.Open(cfg.DSN()), gormConfig)
	case "postgres":
		db, err = gorm.Open(postgres.Open(cfg.DSN()), gormConfig)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(cfg.DSN()), gormConfig)
	default:
		log.Printf("不支持的数据库类型: %s, 默认使用MySQL", cfg.Type)
		db, err = gorm.Open(mysql.Open(cfg.DSN()), gormConfig)
	}
	
	if err != nil {
		return err
	}
	
	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	
	// 设置连接池参数
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	
	// 禁用外键约束检查
	if cfg.Type == "mysql" {
		db.Exec("SET FOREIGN_KEY_CHECKS = 0")
		log.Println("暂时禁用外键检查，以便进行迁移")
	}
	
	// 强制删除attachments表并重新创建
	log.Println("正在强制同步数据库表结构...")
	if err := db.Migrator().DropTable(&Attachment{}); err != nil {
		log.Printf("警告：删除附件表失败: %v", err)
	} else {
		log.Println("成功删除附件表，将重新创建")
	}
	
	// 自动迁移数据库表结构
	err = db.AutoMigrate(
		&User{},
		&Note{},
		&Tag{},
		&Attachment{},
	)
	
	if err != nil {
		return err
	}
	
	// 恢复外键约束
	if cfg.Type == "mysql" {
		db.Exec("SET FOREIGN_KEY_CHECKS = 1")
		log.Println("重新启用外键检查")
	}
	
	DB = db
	log.Printf("数据库连接成功: %s", cfg.Type)
	return nil
}

// 获取数据库实例
func GetDB() *gorm.DB {
	return DB
} 