package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	LoginName   string         `json:"login_name" validate:"required,max=20" gorm:"unique"`
	UserName    string         `json:"user_name" validate:"required,min=3,max=20"`
	PassWord    string         `json:"pass_word" validate:"required,min=6"`
	Phone       string         `json:"phone" validate:"required,len=11,numeric" gorm:"unique"`
	Gender      *int           `json:"gender" validate:"required,oneof=0 1"`
	Status      *int           `json:"status" validate:"required,oneof=1 2 3"`
	LastLoginAt *time.Time     `json:"last_login_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type UserReq struct {
	UserName string `json:"user_name" validate:"required,min=3,max=20"`
	PassWord string `json:"pass_word" validate:"required,min=6"`
	Phone    string `json:"phone" validate:"required,len=11,numeric" gorm:"unique"`
	Gender   *int   `json:"gender" validate:"required,oneof=0 1"`
	Status   *int   `json:"status" validate:"required,oneof=1 2 3"`
}
