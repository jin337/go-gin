package model

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	ID        uint   `json:"id" gorm:"primarykey"`
	LoginName string `json:"login_name" validate:"max=20"`
	UserName  string `json:"user_name" validate:"required,min=3,max=20"`
	PassWord  string `json:"pass_word" validate:"required,min=6"`
	Phone     string `json:"phone" validate:"required,len=11,numeric"`
	IsActive  int    `json:"is_active" validate:"required,oneof=0 1"`

	UserID uint `json:"user_id"`
	// 定义关联关系
	User User `gorm:"foreignKey:UserID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
