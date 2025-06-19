package controller

import (
	"go-gin/internal/app/database"
	"go-gin/internal/service"
	"go-gin/internal/utils"

	"github.com/gin-gonic/gin"
)

type AccountController struct{}

// 新增
func (c *AccountController) CreateAccount(ctx *gin.Context) {
	res, err := service.CreateAccount(ctx, database.DB)
	if err != nil {
		utils.GetResponseJson(ctx, utils.FAIL, err.Error(), nil)
	} else {
		utils.GetResponseJson(ctx, utils.SUCCESS, "新增成功", res)
	}
}

// 查询
func (c *AccountController) GetAccount(ctx *gin.Context) {
	res, err := service.GetAccount(ctx, database.DB)
	if err != nil {
		utils.GetResponseJson(ctx, utils.FAIL, err.Error(), nil)
	} else {

		utils.GetResponseJson(ctx, utils.SUCCESS, "查询成功", res)
	}
}

// 更新
func (c *AccountController) UpdateAccount(ctx *gin.Context) {
	res, err := service.UpdateAccount(ctx, database.DB)
	if err != nil {
		utils.GetResponseJson(ctx, utils.FAIL, err.Error(), nil)
	} else {
		utils.GetResponseJson(ctx, utils.SUCCESS, "更新成功", res)
	}
}

// 删除
func (c *AccountController) DeleteAccount(ctx *gin.Context) {
	res, err := service.DeleteAccount(ctx, database.DB)
	if err != nil {
		utils.GetResponseJson(ctx, utils.FAIL, err.Error(), nil)
	} else {
		utils.GetResponseJson(ctx, utils.SUCCESS, "删除成功", res)
	}
}
