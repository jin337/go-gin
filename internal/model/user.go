package model

import (
	"errors"
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
	IsActive      bool           `json:"is_active" gorm:"default:true"`
	LoginAttempts int            `json:"login_attempts" gorm:"default:0" `
	LastLoginAt   *time.Time     `json:"last_login_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt ` json:"delete_at" gorm:"index"`
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

// 密码加密钩子
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
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

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if tx.Statement.Changed("Password") {
		u.PassWord = hashPassword(u.PassWord)
	}
	return nil
}
