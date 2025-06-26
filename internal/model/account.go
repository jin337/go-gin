package model

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	ID        uint   `json:"id" gorm:"primarykey;comment:唯一标识符"`
	LoginName string `json:"login_name" validate:"max=20" gorm:"unique;comment:账号/登录名"`
	UserName  string `json:"user_name" validate:"required,min=3,max=20" gorm:"comment:昵称"`
	PassWord  string `json:"pass_word" validate:"required,min=6" gorm:"comment:密码"`
	Phone     string `json:"phone" validate:"required,len=11,numeric" gorm:"unique;comment:手机号"`
	IsActive  int    `json:"is_active" validate:"required,oneof=0 1" gorm:"comment:状态 (0: 禁用, 1: 启用)"`

	UserID uint `json:"user_id" gorm:"comment:用户ID"`
	// 定义关联关系
	User User `gorm:"foreignKey:UserID"`

	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;comment:软删除时间"`
}
