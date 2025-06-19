package model

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	ID            uint           `json:"id" gorm:"primarykey"`
	UserID        *uint          `json:"user_id" gorm:"uniqueIndex"`
	UserName      string         `json:"user_name" gorm:"size:20;not null"`
	LoginName     string         `json:"login_name" gorm:"size:20;unique;not null"`
	PassWord      string         `json:"pass_word" gorm:"size:100;not null"`
	IsActive      int            `json:"is_active" gorm:"default:1"`
	LoginAttempts int            `json:"login_attempts" gorm:"default:0"`
	LastLoginAt   *time.Time     `json:"last_login_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// 定义关联关系
	User User `gorm:"foreignKey:UserID"`
}

// 表名
func (Account) TableName() string {
	return "accounts"
}

// 创建前
func (m *Account) BeforeCreate(tx *gorm.DB) (err error) {
	return nil
}

// 创建后
func (m *Account) AfterCreate(tx *gorm.DB) (err error) {
	return nil
}

// 更新前
func (m *Account) BeforeUpdate(tx *gorm.DB) (err error) {
	return nil
}

// 更新后
func (m *Account) AfterUpdate(tx *gorm.DB) (err error) {
	return nil
}

// 删除前
func (m *Account) BeforeDelete(tx *gorm.DB) error {
	return nil
}

// 删除后
func (m *Account) AfterDelete(tx *gorm.DB) (err error) {
	return nil
}
