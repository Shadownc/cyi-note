package models

import (
	"errors"
	"log"
	"time"
	
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	
	"cyi-note/backend/config"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:255;uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"size:255;not null" json:"-"` // 密码不返回给前端
	Email     string         `gorm:"size:255;uniqueIndex" json:"email"`
	Role      string         `gorm:"size:50;default:user" json:"role"` // admin, user
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	
	// 关联
	Notes []Note `gorm:"foreignKey:UserID" json:"-"`
}

// BeforeSave 加密用户密码
func (u *User) BeforeSave(tx *gorm.DB) error {
	// 如果密码没有被加密（即没有$开头），则进行加密
	if len(u.Password) > 0 && u.Password[0] != '$' {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// CheckPassword 检查密码是否正确
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// CreateUser 创建用户
func CreateUser(user *User) error {
	return DB.Create(user).Error
}

// GetUserByID 通过ID获取用户
func GetUserByID(id uint) (*User, error) {
	var user User
	result := DB.First(&user, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByUsername 通过用户名获取用户
func GetUserByUsername(username string) (*User, error) {
	var user User
	err := DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

// GetUserByEmail 通过邮箱获取用户
func GetUserByEmail(email string) (*User, error) {
	var user User
	result := DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, result.Error
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func UpdateUser(user *User) error {
	result := DB.Save(user)
	return result.Error
}

// DeleteUser 删除用户
func DeleteUser(id uint) error {
	return DB.Delete(&User{}, id).Error
}

// EnsureAdminExists 确保管理员用户存在
func EnsureAdminExists(cfg *config.Config) error {
	// 检查是否存在管理员账号
	var count int64
	if err := DB.Model(&User{}).Where("role = ?", "admin").Count(&count).Error; err != nil {
		return err
	}
	
	// 如果不存在管理员账号，则创建一个
	if count == 0 {
		log.Println("未检测到管理员账号，创建默认管理员账号...")
		
		admin := User{
			Username: cfg.AdminUsername,
			Password: cfg.AdminPassword, // 保存时会自动加密
			Email:    cfg.AdminEmail,
			Role:     "admin",
		}
		
		if err := CreateUser(&admin); err != nil {
			return err
		}
		
		log.Printf("已创建默认管理员账号: %s", cfg.AdminUsername)
	} else {
		log.Println("检测到管理员账号已存在，无需创建。")
	}
	
	return nil
} 