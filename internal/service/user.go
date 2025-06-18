package service

import (
	"errors"
	"go-gin/internal/model"
	"go-gin/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ToUserDTO(user *model.User) model.UserDTO {
	return model.UserDTO{
		ID:            user.ID,
		UserName:      user.UserName,
		Email:         user.Email,
		Phone:         user.Phone,
		IsActive:      user.IsActive,
		LoginAttempts: user.LoginAttempts,
		LastLoginAt:   user.LastLoginAt,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
}

// 新增
func CreateUser(c *gin.Context, DB *gorm.DB) (interface{}, error) {
	var body model.User
	// 获取参数
	if err := c.ShouldBind(&body); err != nil {
		return nil, errors.New("参数错误")
	}
	// 在表不存在时自动创建数据库表
	// DB.AutoMigrate(&model.User{})

	// 判断参数是否重复
	searchKey := map[string]interface{}{
		"user_name": body.UserName,
		"phone":     body.Phone,
	}
	if err := utils.CheckUniqueFields(DB, &model.User{}, searchKey); err != nil {
		return nil, err
	}
	// 参数赋值给对象指针
	user := &model.User{
		UserName: body.UserName,
		PassWord: body.PassWord,
		Email:    body.Email,
		Phone:    body.Phone,
	}
	// 插入数据库
	result := DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return ToUserDTO(user), nil
}

// 查询
func GetUser(c *gin.Context, DB *gorm.DB) (interface{}, error) {
	var (
		body     map[string]interface{}
		list     []model.User
		usersDTO []model.UserDTO
		total    int64
	)
	// 获取参数
	if err := c.ShouldBind(&body); err != nil {
		return nil, errors.New("参数错误")
	}

	// 判断是否启用分页
	var page, pageSize int

	if p, ok := body["page"].(int); ok {
		page = p
	} else if pFloat, ok := body["page"].(float64); ok {
		page = int(pFloat)
	}

	if pSize, ok := body["page_size"].(int); ok {
		pageSize = pSize
	} else if pSizeFloat, ok := body["page_size"].(float64); ok {
		pageSize = int(pSizeFloat)
	}

	if page > 0 && pageSize > 0 {
		// 构建查询
		query := DB.Model(&model.User{})
		// 添加查询条件
		conditions := map[string]interface{}{}
		for key, value := range body {
			if key == "page" || key == "page_size" {
				continue
			}
			conditions[key] = value
		}
		for key, value := range conditions {
			query = query.Where(key+" = ?", value)
		}

		// 获取总数
		if err := query.Count(&total).Error; err != nil {
			return nil, err
		}

		// 分页计算
		offset := (page - 1) * pageSize
		result := query.Offset(offset).Limit(pageSize).Find(&list)

		if result.Error != nil {
			return nil, result.Error
		}

		for _, user := range list {
			usersDTO = append(usersDTO, ToUserDTO(&user))
		}
		return map[string]interface{}{
			"list":      usersDTO,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		}, nil
	}

	// 查询全部用户
	result := DB.Model(&model.User{}).Find(&list)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, user := range list {
		usersDTO = append(usersDTO, ToUserDTO(&user))
	}
	return usersDTO, nil
}

// 更新
func UpdateUser(c *gin.Context, DB *gorm.DB) (interface{}, error) {
	// 获取参数
	var (
		body map[string]interface{}
		user model.User
	)
	if err := c.ShouldBind(&body); err != nil {
		return nil, errors.New("参数错误")
	}

	// 判断是否存在
	if err := DB.First(&model.User{}, body["id"]).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 判断参数是否重复
	searchKey := map[string]interface{}{
		"user_name": user.UserName,
		"phone":     user.Phone,
	}
	if err := utils.CheckUniqueFields(DB, &model.User{}, searchKey); err != nil {
		return nil, err
	}

	result := DB.Model(&model.User{}).Where("id = ?", body["id"]).Updates(body)
	if result.Error != nil {
		return nil, errors.New("更新失败")
	}
	if err := DB.Model(&model.User{}).Where("id = ?", body["id"]).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	return ToUserDTO(&user), nil

}

// 删除
func DeleteUser(c *gin.Context, DB *gorm.DB) (interface{}, error) {
	// 获取参数
	var body map[string]interface{}
	if err := c.ShouldBind(&body); err != nil {
		return nil, errors.New("参数错误")
	}
	// 判断是否存在
	if err := DB.First(&model.User{}, body["id"]).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	result := DB.Delete(&model.User{}, body["id"])
	if result.Error != nil {
		return nil, errors.New("删除失败")
	}
	return nil, nil
}
