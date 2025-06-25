package service

import (
	"fmt"
	"go-gin/internal/model"
	"go-gin/internal/utils"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 新增
func CreateUser(ctx *gin.Context, DB *gorm.DB) (interface{}, error) {
	// DB.AutoMigrate(&model.User{}, &model.Account{})

	var req model.UserReq
	err := utils.ValidatorJSON(ctx, &req)
	if err != nil {
		return nil, err
	}
	// 生成用户登录账号,暂时生成一个随机数
	rand.Seed(time.Now().UnixNano())
	loginName := fmt.Sprintf("CRM%06d", rand.Intn(1000000))

	body := &model.User{
		LoginName: loginName,
		UserName:  req.UserName,
		PassWord:  req.PassWord,
		Phone:     req.Phone,
		Gender:    req.Gender,
		Status:    req.Status,
	}
	// 字段校验
	if err := utils.Validator(body); err != nil {
		return nil, err
	}
	// 创建
	if err := DB.Create(body).Error; err != nil {
		return nil, utils.TranslateDBError(err) // 转换错误提示内容
	}

	return nil, nil
}

// 查询
func GetUser(ctx *gin.Context, DB *gorm.DB) (interface{}, error) {
	return nil, nil
}

// 修改
func UpdateUser(ctx *gin.Context, DB *gorm.DB) (interface{}, error) {
	return nil, nil
}

// 删除
func DeleteUser(ctx *gin.Context, DB *gorm.DB) (interface{}, error) {
	return nil, nil
}
