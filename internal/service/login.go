package service

import (
	"errors"
	"go-gin/internal/model"
	"go-gin/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 请求参数
type LoginReq struct {
	LoginName string `json:"login_name" validate:"required"`
	PassWord  string `json:"pass_word" validate:"required"`
}

// 响应参数
type UserRes struct {
	ID          uint       `json:"id"`
	LoginName   string     `json:"login_name" `
	UserName    string     `json:"user_name"`
	Phone       string     `json:"phone"`
	IsActive    *int       `json:"is_active"`
	Gender      *int       `json:"gender"`
	Status      int        `json:"status"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Token       string     `json:"token"`
}

// 登录
func Login(ctx *gin.Context, DB *gorm.DB) (interface{}, error) {
	// 解析校验参数
	var req LoginReq
	if err := utils.ValidatorJSON(ctx, &req); err != nil {
		return nil, err
	}
	var (
		account model.Account
		user    model.User
	)
	err := DB.Transaction(func(tx *gorm.DB) error {
		// 查询用户是否存在
		if err := tx.Model(&account).Where("login_name = ?", req.LoginName).First(&account).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("账号不存在")
			}
			return err
		}
		// 验证密码是否匹配
		if err := bcrypt.CompareHashAndPassword([]byte(account.PassWord), []byte(req.PassWord)); err != nil {
			return errors.New("密码错误")
		}
		// 获取账号关联用户
		if err := tx.Model(&user).Preload("Account").Where("id = ?", account.UserID).First(&user).Error; err != nil {
			return err
		}
		// 更新最后登录时间
		now := time.Now()
		if err := tx.Model(&account).Update("last_login_at", &now).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, utils.TranslateDBError(err) // 转换错误提示内容
	}

	// 生成token
	token, err := utils.GenerateToken(user.ID, 1*60) // 生成token有效期:1小时
	if err != nil {
		return nil, err
	}

	return UserRes{
		ID:          user.ID,
		LoginName:   user.Account.LoginName,
		UserName:    user.UserName,
		Phone:       user.Phone,
		Gender:      user.Gender,
		Status:      user.Status,
		LastLoginAt: user.Account.LastLoginAt,
		IsActive:    &user.Account.IsActive,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Token:       token,
	}, nil
}

// 退出
func Logout(ctx *gin.Context, DB *gorm.DB) (interface{}, error) {
	authHeader := ctx.GetHeader("Token")
	if authHeader == "" {
		return nil, errors.New("需要授权Token")
	}

	blacklistedToken := model.BlacklistedToken{
		Token:     authHeader,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	if err := DB.Create(&blacklistedToken).Error; err != nil {
		return nil, utils.TranslateDBError(err) // 转换错误提示内容
	}
	return nil, nil
}
