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
	
	DB = db
	log.Printf("数据库连接成功: %s", cfg.Type)
	return nil
}

// 获取数据库实例
func GetDB() *gorm.DB {
	return DB
} 