package service

import (
	"go-gin/internal/model"
	"go-gin/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 新增
func CreateUser(ctx *gin.Context, DB *gorm.DB) (interface{}, error) {
	var req model.UserReq
	err := utils.ValidatorJSON(ctx, &req)
	if err != nil {
		return nil, err
	}
	body := &model.User{
		LoginName: "111",
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
