package controller

import (
	"go-gin/internal/app/database"
	"go-gin/internal/service"
	"go-gin/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct{}

// 新增
func (c *UserController) CreateUser(ctx *gin.Context) {
	res, err := service.CreateUser(ctx, database.DB)
	if err != nil {
		utils.GetResponseJson(ctx, utils.FAIL, err.Error(), nil)
	} else {
		utils.GetResponseJson(ctx, utils.SUCCESS, "新增成功", res)
	}
}

// 查询
func (c *UserController) GetUserList(ctx *gin.Context) {
	res, err := service.GetUser(ctx, database.DB)
	if err != nil {
		utils.GetResponseJson(ctx, utils.FAIL, err.Error(), nil)
	} else {

		utils.GetResponseJson(ctx, utils.SUCCESS, "查询成功", res)
	}
}

// 更新
func (c *UserController) UpdateUser(ctx *gin.Context) {
	res, err := service.UpdateUser(ctx, database.DB)
	if err != nil {
		utils.GetResponseJson(ctx, utils.FAIL, err.Error(), nil)
	} else {
		utils.GetResponseJson(ctx, utils.SUCCESS, "更新成功", res)
	}
}

// 删除
func (c *UserController) DeleteUser(ctx *gin.Context) {
	res, err := service.DeleteUser(ctx, database.DB)
	if err != nil {
		utils.GetResponseJson(ctx, utils.FAIL, err.Error(), nil)
	} else {
		utils.GetResponseJson(ctx, utils.SUCCESS, "删除成功", res)
	}
}
