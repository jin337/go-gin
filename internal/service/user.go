package service

import (
	"errors"
	"fmt"
	"go-gin/internal/model"
	"go-gin/internal/utils"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DB.AutoMigrate(&model.User{})
// 新增
func CreateUser(ctx *gin.Context, DB *gorm.DB) (interface{}, error) {
	// 解析校验参数
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
	// 创建
	if err := DB.Create(body).Error; err != nil {
		return nil, utils.TranslateDBError(err) // 转换错误提示内容
	}

	return nil, nil
}

// 查询
func GetUser(ctx *gin.Context, DB *gorm.DB) (interface{}, error) {
	// 定义切片用于存储查询结果
	var list []model.User

	// 执行查询,Omit排除字段
	result := DB.Find(&list)
	if result.Error != nil {
		return nil, utils.TranslateDBError(result.Error) // 转换错误提示内容
	}
	// 处理空结果集
	if len(list) == 0 {
		return []model.User{}, nil
	}
	// 转换结果集
	var res []model.UserRes
	for _, item := range list {
		res = append(res, model.UserRes{
			ID:        item.ID,
			UserName:  item.UserName,
			Phone:     item.Phone,
			Gender:    item.Gender,
			Status:    item.Status,
			LoginName: item.LoginName,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}
	return res, nil
}

// 修改
func UpdateUser(ctx *gin.Context, DB *gorm.DB) (interface{}, error) {
	// 解析校验参数
	var req model.UserUpdateReq
	err := utils.ValidatorJSON(ctx, &req)
	if err != nil {
		return nil, err
	}
	// 检查用户是否存在
	var item model.User
	if err := DB.First(&item, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, utils.TranslateDBError(err) // 转换错误提示内容
	}
	// 更新
	updateData := map[string]interface{}{
		"updated_at": time.Now(),
	}
	if req.UserName != nil {
		updateData["user_name"] = *req.UserName
	}
	if req.Phone != nil {
		updateData["phone"] = *req.Phone
	}
	if req.Gender != nil {
		updateData["gender"] = *req.Gender
	}
	if req.Status != nil {
		updateData["status"] = *req.Status
	}

	result := DB.Model(&item).Updates(updateData)
	if result.Error != nil {
		return nil, utils.TranslateDBError(result.Error) // 转换错误提示内容
	}
	// 查询更新数据
	if err := DB.First(&item, req.ID).Error; err != nil {
		return nil, utils.TranslateDBError(result.Error) // 转换错误提示内容
	}
	res := model.UserRes{
		ID:        item.ID,
		UserName:  item.UserName,
		Phone:     item.Phone,
		Gender:    item.Gender,
		Status:    item.Status,
		LoginName: item.LoginName,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
	return res, nil
}

// 删除
func DeleteUser(ctx *gin.Context, DB *gorm.DB) (interface{}, error) {
	// 解析校验参数
	var req model.UserUpdateReq
	err := utils.ValidatorJSON(ctx, &req)
	if err != nil {
		return nil, err
	}
	// 检查用户是否存在
	var item model.User
	if err := DB.First(&item, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, utils.TranslateDBError(err) // 转换错误提示内容
	}
	// 删除
	if err := DB.Delete(&item).Error; err != nil {
		return nil, utils.TranslateDBError(err)
	}

	return nil, nil
}
