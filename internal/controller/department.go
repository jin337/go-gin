package controller

import (
	"go-gin/internal/app/database"
	"go-gin/internal/service"
	"go-gin/internal/utils"

	"github.com/gin-gonic/gin"
)

type DepartmentController struct{}

// 新增
func (c *DepartmentController) CreateDepartment(ctx *gin.Context) {
	res, err := service.CreateDepartment(ctx, database.DB)
	if err != nil {
		utils.GetResponseJson(ctx, utils.FAIL, err.Error(), nil)
	} else {
		utils.GetResponseJson(ctx, utils.SUCCESS, "新增成功", res)
	}
}

// 查询
func (c *DepartmentController) GetDepartmentList(ctx *gin.Context) {
	res, err := service.GetDepartmentList(ctx, database.DB)
	if err != nil {
		utils.GetResponseJson(ctx, utils.FAIL, err.Error(), nil)
	} else {

		utils.GetResponseJson(ctx, utils.SUCCESS, "查询成功", res)
	}
}

// 更新
func (c *DepartmentController) UpdateDepartment(ctx *gin.Context) {
	res, err := service.UpdateDepartment(ctx, database.DB)
	if err != nil {
		utils.GetResponseJson(ctx, utils.FAIL, err.Error(), nil)
	} else {
		utils.GetResponseJson(ctx, utils.SUCCESS, "更新成功", res)
	}
}

// 删除
func (c *DepartmentController) DeleteDepartment(ctx *gin.Context) {
	res, err := service.DeleteDepartment(ctx, database.DB)
	if err != nil {
		utils.GetResponseJson(ctx, utils.FAIL, err.Error(), nil)
	} else {
		utils.GetResponseJson(ctx, utils.SUCCESS, "删除成功", res)
	}
}
