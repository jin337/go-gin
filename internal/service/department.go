package service

import (
	"errors"
	"go-gin/internal/model"
	"go-gin/internal/utils"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 新增
func CreateDepartment(ctx *gin.Context, DB *gorm.DB) (interface{}, error) {
	// 解析校验参数
	var req model.DepartmentReq
	if err := utils.ValidatorJSON(ctx, &req); err != nil {
		return nil, err
	}

	body := &model.Department{
		Name:     req.Name,
		Pid:      *req.Pid,
		Type:     req.Type,
		HeadUpId: req.HeadUpId,
	}

	if err := DB.Create(body).Error; err != nil {
		return nil, utils.TranslateDBError(err)
	}
	return nil, nil
}

// 查询
func GetDepartmentList(ctx *gin.Context, DB *gorm.DB) (interface{}, error) {
	return nil, nil
}

// 更新
func UpdateDepartment(ctx *gin.Context, DB *gorm.DB) (interface{}, error) {
	// 解析校验参数
	var req model.DepartmentUpdateReq
	if err := utils.ValidatorJSON(ctx, &req); err != nil {
		return nil, err
	}

	var result model.DepartmentRes
	err := DB.Transaction(func(tx *gorm.DB) error {
		// 检查数据是否存在
		var item model.Department
		if err := DB.First(&item, req.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("部门不存在")
			}
			return err
		}
		// 准备更新数据
		updateData := map[string]interface{}{
			"updated_at": time.Now(),
		}
		if req.Name != nil {
			updateData["name"] = *req.Name
		}
		if req.Pid != nil {
			updateData["pid"] = *req.Pid
		}
		if req.Type != nil {
			updateData["type"] = *req.Type
		}
		if req.HeadUpId != nil {
			updateData["head_up_id"] = *req.HeadUpId
		}

		// 执行更新
		if err := DB.Model(&item).Updates(updateData).Error; err != nil {
			return utils.TranslateDBError(err) // 转换错误提示内容
		}

		// 查询更新后数据
		if err := DB.First(&item, req.ID).Error; err != nil {
			return err
		}
		result = model.DepartmentRes{
			ID:        item.ID,
			Name:      item.Name,
			Pid:       item.Pid,
			Type:      item.Type,
			HeadUpId:  item.HeadUpId,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		}
		return nil
	})

	if err != nil {
		return nil, utils.TranslateDBError(err) // 转换错误提示内容
	}

	return result, nil
}

// 删除
func DeleteDepartment(ctx *gin.Context, DB *gorm.DB) (interface{}, error) {
	return nil, nil
}
