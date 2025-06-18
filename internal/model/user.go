package model

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID            uint           `json:"id" gorm:"primarykey"`
	UserName      string         `json:"user_name" gorm:"size:20;uniqueIndex;not null"`
	PassWord      string         `json:"pass_word" gorm:"size:100;not null"`
	Email         string         `json:"email" gorm:"size:50;uniqueIndex;not null"`
	Phone         *string        `json:"phone" gorm:"size:11;uniqueIndex"`
	IsActive      int            `json:"is_active" gorm:"default:1"`
	LoginAttempts int            `json:"login_attempts" gorm:"default:0"`
	LastLoginAt   *time.Time     `json:"last_login_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
type UserDTO struct {
	ID            uint       `json:"id"`
	UserName      string     `json:"user_name" binding:"required"`
	Email         string     `json:"email" binding:"required"`
	Phone         *string    `json:"phone"`
	IsActive      int        `json:"is_active"`
	LoginAttempts int        `json:"login_attempts"`
	LastLoginAt   *time.Time `json:"last_login_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// 表名
func (User) TableName() string {
	return "users"
}

// 密码加密函数
func hashPassword(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)
}

// 创建前
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Println("\033[32m创建前执行\033[0m")
	// 处理逻辑
	if u.PassWord != "" {
		u.PassWord = hashPassword(u.PassWord)
	}
	if u.UserName == "" {
		return errors.New("用户名不能为空")
	}
	if u.PassWord == "" {
		return errors.New("密码不能为空")
	}
	if u.Email == "" {
		return errors.New("邮箱不能为空")
	}
	return nil
}

// 查询前
func (u *User) BeforeFind(tx *gorm.DB) (err error) {
	fmt.Println("\033[32m查询前执行\033[0m")
	// 处理逻辑
	return nil
}

// 更新前
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	fmt.Println("\033[32m更新前执行\033[0m")
	// 处理逻辑
	if tx.Statement.Changed("Password") {
		u.PassWord = hashPassword(u.PassWord)
	}
	return nil
}

// 删除前
func (u *User) BeforeDelete(tx *gorm.DB) error {
	fmt.Println("\033[32m删除前执行\033[0m")
	// 处理逻辑
	return nil
}
