package model

import (
	"time"

	"gorm.io/gorm"
)

// 表结构
type User struct {
	ID        uint           `json:"id" gorm:"primarykey;comment:唯一标识符"`
	UserName  string         `json:"user_name" validate:"required,min=3,max=20" gorm:"comment:昵称"`
	PassWord  string         `json:"pass_word" validate:"required,min=6" gorm:"comment:密码"`
	Phone     string         `json:"phone" validate:"required,len=11,numeric" gorm:"unique;comment:手机号"`
	Gender    *int           `json:"gender" validate:"required,oneof=0 1" gorm:"comment:性别 (0: 女, 1: 男)"`
	Status    int            `json:"status" validate:"required,oneof=1 2 3" gorm:"comment:状态 (1: 正式, 2: 试用期, 3: 离职)"`
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;comment:软删除时间"`

	Account Account `json:"account" gorm:"foreignkey:UserID"`
}

// 响应字段
type UserRes struct {
	ID          uint       `json:"id" gorm:"primarykey"`
	LoginName   string     `json:"login_name" validate:"required,max=20" gorm:"unique"`
	UserName    string     `json:"user_name" validate:"required,min=3,max=20"`
	Phone       string     `json:"phone" validate:"required,len=11,numeric" gorm:"unique"`
	IsActive    *int       `json:"is_active" validate:"required,oneof=0 1"`
	Gender      *int       `json:"gender" validate:"required,oneof=0 1"`
	Status      int        `json:"status" validate:"required,oneof=1 2 3"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// 请求字段-新增
type UserReq struct {
	UserName string `json:"user_name" validate:"required,min=3,max=20"`
	PassWord string `json:"pass_word" validate:"required,min=6"`
	Phone    string `json:"phone" validate:"required,len=11,numeric" gorm:"unique"`
	Gender   *int   `json:"gender" validate:"required,oneof=0 1"`
	Status   int    `json:"status" validate:"required,oneof=1 2 3"`
}

// 请求参数-更新
type UserUpdateReq struct {
	ID        uint    `json:"id" gorm:"primarykey"`
	UserName  *string `json:"user_name" validate:"omitempty,min=3,max=20"`
	Phone     *string `json:"phone" validate:"omitempty,len=11,numeric" gorm:"unique"`
	Gender    *int    `json:"gender" validate:"omitempty,oneof=0 1"`
	Status    *int    `json:"status" validate:"omitempty,oneof=1 2 3"`
	AccountID *int    `json:"account_id" validate:"omitempty"`
}

// 请求参数-查询
type UserOffsetReq struct {
	ID       *uint   `json:"id" validate:"omitempty" gorm:"primarykey"`
	UserName *string `json:"user_name" validate:"omitempty,min=3,max=20"`
	Phone    *string `json:"phone" validate:"omitempty,len=11,numeric" gorm:"unique"`
	Gender   *int    `json:"gender" validate:"omitempty,oneof=0 1"`
	Status   *int    `json:"status" validate:"omitempty,oneof=1 2 3"`
	Page     int     `json:"page" validate:"required"`
	PageSize int     `json:"page_size" validate:"required"`
}
