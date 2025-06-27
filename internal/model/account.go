package model

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	UserID uint `json:"user_id" gorm:"comment:用户ID;default:NULL"`

	ID          uint       `json:"id" gorm:"primarykey;comment:唯一标识符"`
	LoginName   string     `json:"login_name" validate:"max=20" gorm:"unique;comment:账号/登录名"`
	UserName    string     `json:"user_name" validate:"required,min=3,max=20" gorm:"comment:昵称"`
	PassWord    string     `json:"pass_word" validate:"required,min=6" gorm:"comment:密码"`
	Phone       string     `json:"phone" validate:"required,len=11,numeric" gorm:"unique;comment:手机号"`
	IsActive    int        `json:"is_active" validate:"required,oneof=0 1" gorm:"comment:状态 (0: 禁用, 1: 启用)"`
	LastLoginAt *time.Time `json:"last_login_at" gorm:"comment:最后一次登录时间"`

	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;comment:软删除时间"`
}

// 响应字段
type AccountRes struct {
	ID          uint       `json:"id"`
	UserID      uint       `json:"user_id"`
	LoginName   string     `json:"login_name"`
	UserName    string     `json:"user_name"`
	Phone       string     `json:"phone"`
	IsActive    int        `json:"is_active"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// 请求字段-新增
type AccountReq struct {
	UserName string `json:"user_name" validate:"required,min=3,max=20"`
	PassWord string `json:"pass_word" validate:"required,min=6"`
	Phone    string `json:"phone" validate:"required,len=11,numeric" gorm:"unique"`
}

// 请求参数-更新
type AccountUpdateReq struct {
	ID       uint    `json:"id" validate:"required"`
	UserName *string `json:"user_name" validate:"omitempty,min=3,max=20"`
	Phone    *string `json:"phone" validate:"omitempty,len=11,numeric" gorm:"unique"`
	IsActive *int    `json:"is_active" validate:"omitempty,oneof=0 1"`
	UserID   *uint   `json:"user_id" validate:"omitempty"`
}

// 请求参数-查询
type AccountOffsetReq struct {
	UserID   *uint   `json:"user_id" validate:"omitempty"`
	ID       *uint   `json:"id" validate:"omitempty"`
	UserName *string `json:"user_name" validate:"omitempty,min=3,max=20"`
	Phone    *string `json:"phone" validate:"omitempty,len=11,numeric" gorm:"unique"`
	IsActive *int    `json:"is_active" validate:"omitempty,oneof=0 1"`
	Page     int     `json:"page" validate:"required"`
	PageSize int     `json:"page_size" validate:"required"`
}
