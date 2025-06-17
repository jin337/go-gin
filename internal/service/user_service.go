package service

import (
	"errors"
	"go-gin/internal/model"
	"go-gin/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 新增
func CreateUser(c *gin.Context, DB *gorm.DB) (interface{}, error) {
	// 获取参数
	var body model.User
	if err := c.ShouldBind(&body); err != nil {
		return nil, errors.New("参数错误")
	}

	user := &model.User{
		UserName: body.UserName,
		PassWord: body.PassWord,
		Email:    body.Email,
		Phone:    body.Phone,
	}

	// 判断是否存在
	searchKey := map[string]interface{}{
		"user_name": body.UserName,
	}
	if err := utils.CheckUniqueFields(DB, &model.User{}, searchKey); err != nil {
		return nil, err
	}

	// DB.AutoMigrate(&model.User{}) // 在表不存在时自动创建数据库表
	result := DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	} else {
		return map[string]interface{}{
			"id":        user.ID,
			"user_name": user.UserName,
			"email":     user.Email,
			"phone":     user.Phone,
			"is_active": user.IsActive,
		}, nil
	}
}

// 查询
func GetUser(c *gin.Context, DB *gorm.DB) (interface{}, error) {
	var body model.User
	name := c.Query("user_name")
	result := DB.Where("user_name = ?", name).First(&body)
	if result.Error != nil {
		return nil, errors.New("用户不存在")
	} else {
		return body, nil
	}
}

// 更新
func UpdateUser(c *gin.Context, DB *gorm.DB) (interface{}, error) {
	// 获取参数
	var body model.User
	if err := c.ShouldBind(&body); err != nil {
		return nil, errors.New("参数错误")
	}

	// 判断是否存在
	searchKey := map[string]interface{}{
		"user_name": body.UserName,
	}
	if err := utils.CheckUniqueFields(DB, &model.User{}, searchKey); err != nil {
		return nil, err
	}

	result := DB.Model(&body).Where("id=?", body.ID).Updates(body)
	if result.Error != nil {
		return nil, errors.New("更新失败")
	} else {
		return body, nil
	}
}

// 删除
func DeleteUser(c *gin.Context, DB *gorm.DB) (interface{}, error) {
	return nil, errors.New("删除失败")
}
