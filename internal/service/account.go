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
func CreateAccount(ctx *gin.Context, DB *gorm.DB) (interface{}, error) {
	// 解析校验参数
	var req model.AccountReq
	if err := utils.ValidatorJSON(ctx, &req); err != nil {
		return nil, err
	}

	// 生成登录账号,暂时使用时间戳
	loginName := fmt.Sprintf("CRM%06d", time.Now().Unix())

	// 明文密码加密为安全的 bcrypt 哈希值
	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.PassWord), bcrypt.DefaultCost)
	passWord := string(hashed)
	// 创建账号
	body := &model.Account{
		LoginName: loginName,
		PassWord:  passWord,
		UserName:  req.UserName,
		Phone:     req.Phone,
		IsActive:  1,
	}

	if err := DB.Create(body).Error; err != nil {
		return nil, utils.TranslateDBError(err) // 转换错误提示内容
	}

	return nil, nil
}

// 查询
func GetAccount(ctx *gin.Context, DB *gorm.DB) (interface{}, error) {
	// list:切片用于存储查询结果
	// total:总数
	var (
		list  []model.Account
		total int64
	)

	// 解析校验参数
	var req model.AccountOffsetReq
	if err := utils.ValidatorJSON(ctx, &req); err != nil {
		return nil, err
	}

	// 数据
	condition := map[string]interface{}{}
	if req.ID != nil {
		condition["id"] = *req.ID
	}
	if req.UserID != nil {
		condition["user_id"] = *req.UserID
	}
	if req.UserName != nil {
		condition["user_name"] = *req.UserName
	}
	if req.Phone != nil {
		condition["phone"] = *req.Phone
	}
	if req.IsActive != nil {
		condition["is_active"] = *req.IsActive
	}
	query := DB.Model(&list).Where(condition)

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
			"list":  []model.AccountRes{},
			"total": total,
		}, nil
	}

	// 转换结果集
	var res []model.AccountRes
	for _, item := range list {
		res = append(res, model.AccountRes{
			ID:          item.ID,
			UserID:      item.UserID,
			LoginName:   item.LoginName,
			UserName:    item.UserName,
			Phone:       item.Phone,
			IsActive:    item.IsActive,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
			LastLoginAt: item.LastLoginAt,
		})
	}

	return map[string]interface{}{
		"list":  res,
		"total": total,
	}, nil
}

// 修改
func UpdateAccount(ctx *gin.Context, DB *gorm.DB) (interface{}, error) {
	// 解析校验参数
	var req model.AccountUpdateReq
	if err := utils.ValidatorJSON(ctx, &req); err != nil {
		return nil, err
	}

	var result model.AccountRes
	var item model.Account

	// 使用事务处理核心数据库操作
	err := DB.Transaction(func(tx *gorm.DB) error {
		// 检查用户是否存在
		if err := tx.First(&item, req.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("账号不存在")
			}
			return err
		}

		// 准备更新数据
		updateData := map[string]interface{}{
			"updated_at": time.Now(),
		}
		if req.UserName != nil {
			updateData["user_name"] = *req.UserName
		}
		if req.Phone != nil {
			updateData["phone"] = *req.Phone
		}
		if req.IsActive != nil {
			updateData["is_active"] = *req.IsActive
		}

		// 执行更新
		if err := tx.Model(&item).Updates(updateData).Error; err != nil {
			return err
		}

		// 在事务内重新查询确保数据一致性
		if err := tx.First(&item, req.ID).Error; err != nil {
			return err
		}

		// 准备返回结果
		result = model.AccountRes{
			ID:          item.ID,
			LoginName:   item.LoginName,
			UserName:    item.UserName,
			Phone:       item.Phone,
			IsActive:    item.IsActive,
			LastLoginAt: item.LastLoginAt,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		}

		return nil
	})

	if err != nil {
		return nil, utils.TranslateDBError(err)
	}

	return result, nil
}

// 删除
func DeleteAccount(ctx *gin.Context, DB *gorm.DB) (interface{}, error) {
	// 解析校验参数
	var req model.AccountUpdateReq
	if err := utils.ValidatorJSON(ctx, &req); err != nil {
		return nil, err
	}
	// 检查账号是否存在
	var item model.Account
	if err := DB.First(&item, req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("账号不存在")
		}
		return nil, utils.TranslateDBError(err) // 转换错误提示内容
	}

	// 删除账号
	if err := DB.Delete(&item).Error; err != nil {
		return nil, utils.TranslateDBError(err)
	}

	return nil, nil
}
