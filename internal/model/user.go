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
	UserName      string         `json:"user_name" gorm:"size:20;not null"`
	LoginName     *string        `json:"login_name" gorm:"size:20;uniqueIndex"`
	PassWord      string         `json:"pass_word" gorm:"size:100;not null"`
	Phone         string         `json:"phone" gorm:"size:11;uniqueIndex"`
	Gender        int            `json:"gender" gorm:"check:gender in (0,1);default:1"`       // 性别 0:女 1:男
	IsActive      int            `json:"is_active" gorm:"check:is_active in (0,1);default:1"` // 状态 0:禁用 1:正常 2:锁定
	Status        int            `json:"status" gorm:"check:status in (1,2,3);default:1"`     // 类型 1:在职 2:离职 3:试用期
	LoginAttempts int            `json:"login_attempts" gorm:"default:0"`
	LastLoginAt   *time.Time     `json:"last_login_at"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
type UserDTO struct {
	ID            uint       `json:"id"`
	UserName      string     `json:"user_name"`
	LoginName     *string    `json:"login_name"`
	Phone         string     `json:"phone"`
	Gender        int        `json:"gender"`
	IsActive      int        `json:"is_active"`
	Status        int        `json:"status"`
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
	if u.Phone == "" {
		return errors.New("手机号不能为空")
	}
	return nil
}

// 创建后
func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	fmt.Println("\033[32m创建后执行\033[0m")
	// 固定字段+自增6为账号
	loginName := fmt.Sprintf("%s%06d", "CN", u.ID)
	u.LoginName = &loginName
	return nil
}

// 更新前
func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	fmt.Println("\033[32m更新前执行\033[0m")
	// 处理逻辑
	if tx.Statement.Changed("Password") {
		u.PassWord = hashPassword(u.PassWord)
	}
	if tx.Statement.Changed("IsActive") {
		if u.IsActive != 0 && u.IsActive != 1 {
			return errors.New("状态必须为 0 或 1")
		}
	}
	if tx.Statement.Changed("Status") {
		if u.Status != 1 && u.Status != 2 && u.Status != 3 {
			return errors.New("状态必须为 1、2 或 3")
		}
	}
	if tx.Statement.Changed("Gender") {
		if u.Gender != 0 && u.Gender != 1 {
			return errors.New("性别必须为 0 或 1")
		}
	}
	return nil
}

// 删除前
func (u *User) BeforeDelete(tx *gorm.DB) error {
	fmt.Println("\033[32m删除前执行\033[0m")
	// 处理逻辑
	return nil
}
