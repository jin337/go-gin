package service

import (
	"errors"
	"go-gin/internal/model"
	"go-gin/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 统一返回内容
func ToUserDTO(r *model.User) model.UserDTO {
	return model.UserDTO{
		ID:            r.ID,
		UserName:      r.UserName,
		Email:         r.Email,
		Phone:         r.Phone,
		IsActive:      r.IsActive,
		LoginAttempts: r.LoginAttempts,
		LastLoginAt:   r.LastLoginAt,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
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
	// 构建查询
	query := DB.Model(&model.User{})

	// 提取分页与查询条件
	req, err := utils.GetOffsetFields(body)
	if err != nil {
		return nil, errors.New("参数解析失败")
	}

	// 添加查询条件
	for key, value := range req.Conditions {
		query = query.Where(key+" = ?", value)
	}

	// 分页查询
	if req.Page > 0 && req.PageSize > 0 {

		// 分页计算
		offset := (req.Page - 1) * req.PageSize
		result := query.Offset(offset).Limit(req.PageSize).Find(&list)

		if result.Error != nil {
			return nil, result.Error
		}

		for _, v := range list {
			usersDTO = append(usersDTO, ToUserDTO(&v))
		}

		// 获取总数
		if err := query.Count(&total).Error; err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"list":      usersDTO,
			"total":     total,
			"page":      req.Page,
			"page_size": req.PageSize,
		}, nil
	}

	// 查询全部用户
	result := query.Find(&list)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, v := range list {
		usersDTO = append(usersDTO, ToUserDTO(&v))
	}
	return usersDTO, nil
}

// 更新
func UpdateUser(c *gin.Context, DB *gorm.DB) (interface{}, error) {
	var (
		body map[string]interface{}
		user model.User
	)
	// 获取参数
	if err := c.ShouldBind(&body); err != nil {
		return nil, errors.New("参数错误")
	}

	// 判断是否存在
	if err := DB.First(&model.User{}, body["id"]).Error; err != nil {
		return nil, errors.New("数据不存在")
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
		return nil, errors.New("数据不存在")
	}
	return ToUserDTO(&user), nil

}

// 删除
func DeleteUser(c *gin.Context, DB *gorm.DB) (interface{}, error) {
	var body map[string]interface{}
	// 获取参数
	if err := c.ShouldBind(&body); err != nil {
		return nil, errors.New("参数错误")
	}
	// 判断是否存在
	if err := DB.First(&model.User{}, body["id"]).Error; err != nil {
		return nil, errors.New("数据不存在")
	}
	// 删除
	result := DB.Delete(&model.User{}, body["id"])
	if result.Error != nil {
		return nil, errors.New("删除失败")
	}
	return nil, nil
}
