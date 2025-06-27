package service

import (
	"errors"
	"fmt"
	"go-gin/internal/model"
	"go-gin/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 新增
func CreateUser(ctx *gin.Context, DB *gorm.DB) (interface{}, error) {
	// 解析校验参数
	var req model.UserReq
	if err := utils.ValidatorJSON(ctx, &req); err != nil {
		return nil, err
	}

	body := &model.User{
		UserName: req.UserName,
		Phone:    req.Phone,
		Gender:   req.Gender,
		Status:   req.Status,
	}
	// 创建员工
	if err := DB.Create(body).Error; err != nil {
		return nil, utils.TranslateDBError(err) // 转换错误提示内容
	}

	// 生成登录账号,暂时使用时间戳
	loginName := fmt.Sprintf("CRM%06d", time.Now().Unix())

	// 明文密码加密为安全的 bcrypt 哈希值
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.PassWord), bcrypt.DefaultCost)
	passWord := string(hashed)
	// 创建员工账号
	account := &model.Account{
		UserID:    body.ID,
		LoginName: loginName,
		UserName:  body.UserName,
		PassWord:  passWord,
		Phone:     body.Phone,
		IsActive:  1,
	}
	if err := DB.Create(account).Error; err != nil {
		return nil, utils.TranslateDBError(err) // 转换错误提示内容
	}

	return nil, nil
}

// 查询
func GetUser(ctx *gin.Context, DB *gorm.DB) (interface{}, error) {
	// list:切片用于存储查询结果
	// total:总数
	var (
		list  []model.User
		total int64
	)

	// 解析校验参数
	var req model.UserOffsetReq
	if err := utils.ValidatorJSON(ctx, &req); err != nil {
		return nil, err
	}
	// 数据
	condition := map[string]interface{}{}
	if req.ID != nil {
		condition["id"] = *req.ID
	}
	if req.UserName != nil {
		condition["user_name"] = *req.UserName
	}
	if req.Phone != nil {
		condition["phone"] = *req.Phone
	}
	if req.Gender != nil {
		condition["gender"] = *req.Gender
	}
	if req.Status != nil {
		condition["status"] = *req.Status
	}

	query := DB.Model(&list).Where(condition).Preload("Account")

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, utils.TranslateDBError(err) // 转换错误提示内容
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(offset).Limit(req.PageSize).Find(&list).Error; err != nil {
		return nil, utils.TranslateDBError(err) // 转换错误提示内容
	}

	// 处理空结果集
	if len(list) == 0 {
		return map[string]interface{}{
			"list":  []model.UserRes{},
			"total": total,
		}, nil
	}

	// 转换结果集
	var res []model.UserRes
	for _, item := range list {
		res = append(res, model.UserRes{
			ID:          item.ID,
			LoginName:   item.Account.LoginName,
			UserName:    item.UserName,
			Phone:       item.Phone,
			Gender:      item.Gender,
			Status:      item.Status,
			LastLoginAt: item.Account.LastLoginAt,
			IsActive:    &item.Account.IsActive,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		})
	}

	return map[string]interface{}{
		"list":  res,
		"total": total,
	}, nil
}

// 更新
func UpdateUser(ctx *gin.Context, DB *gorm.DB) (interface{}, error) {
	// 解析校验参数
	var req model.UserUpdateReq
	if err := utils.ValidatorJSON(ctx, &req); err != nil {
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

	// 数据
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
	// 提交更新
	err := DB.Model(&item).Updates(updateData).Error
	if err != nil {
		return nil, utils.TranslateDBError(err) // 转换错误提示内容
	}

	// 同步账号内容
	accountData := map[string]interface{}{
		"updated_at": time.Now(),
	}
	if req.UserName != nil {
		accountData["user_name"] = *req.UserName
	}
	if req.Phone != nil {
		accountData["phone"] = *req.Phone
	}

	// 处理账号绑定关系
	if req.AccountID != nil {
		var existingAccount model.Account
		// 检查新账号是否已被其他用户绑定
		if err := DB.Where("id = ? AND user_id IS NOT NULL AND user_id != ?", req.AccountID, item.ID).
			First(&existingAccount).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("该账号已被其他用户绑定")
			}
		}

		// 清除当前用户之前的绑定记录
		if err := DB.Model(&model.Account{}).
			Where("user_id = ?", item.ID).
			Update("user_id", nil).Error; err != nil {
			return nil, utils.TranslateDBError(err)
		}

		// 设置新的绑定
		if err := DB.Model(&model.Account{}).
			Where("id = ?", req.AccountID).
			Updates(map[string]interface{}{
				"user_id":    item.ID,
				"updated_at": time.Now(),
			}).Error; err != nil {
			return nil, utils.TranslateDBError(err)
		}
	}
	// 更新内容
	if err := DB.Model(&model.Account{}).Where("user_id=?", item.ID).Updates(accountData).Error; err != nil {
		return nil, utils.TranslateDBError(err)
	}

	// 查询更新数据
	if err := DB.Preload("Account").First(&item, req.ID).Error; err != nil {
		return nil, utils.TranslateDBError(err)
	}
	// 转换结果
	return model.UserRes{
		ID:          item.ID,
		LoginName:   item.Account.LoginName,
		UserName:    item.UserName,
		Phone:       item.Phone,
		Gender:      item.Gender,
		Status:      item.Status,
		LastLoginAt: item.Account.LastLoginAt,
		IsActive:    &item.Account.IsActive,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}, nil
}

// 删除
func DeleteUser(ctx *gin.Context, DB *gorm.DB) (interface{}, error) {
	// 解析校验参数
	var req model.UserUpdateReq
	if err := utils.ValidatorJSON(ctx, &req); err != nil {
		return nil, err
	}

	// 检查用户是否存在
	var item model.User
	if err := DB.Preload("Account").First(&item, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, utils.TranslateDBError(err) // 转换错误提示内容
	}

	// 如果存在关联账户，先禁用(is_active=0)再删除
	if item.Account.ID != 0 {
		// 更新is_active为0
		if err := DB.Model(&item.Account).Update("is_active", 0).Error; err != nil {
			return nil, utils.TranslateDBError(err)
		}

		// 删除账户记录
		if err := DB.Delete(&item.Account).Error; err != nil {
			return nil, utils.TranslateDBError(err)
		}
	}

	// 删除用户
	if err := DB.Delete(&item).Error; err != nil {
		return nil, utils.TranslateDBError(err)
	}

	return nil, nil
}
