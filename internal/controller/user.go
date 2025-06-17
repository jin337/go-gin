package controller

import (
	"go-gin/internal/app/database"
	"go-gin/internal/service"
	"go-gin/internal/utils"

	"github.com/gin-gonic/gin"
)

type UserController struct {
}

// 新增
func (u *UserController) Create(c *gin.Context) {
	res, err := service.CreateUser(c, database.DB)
	if err != nil {
		utils.GetResponseJson(c, utils.FAIL, err.Error(), nil)
	} else {
		utils.GetResponseJson(c, utils.SUCCESS, "新增成功", res)
	}
}

// 查询
func (u *UserController) GetUser(c *gin.Context) {
	res, err := service.GetUser(c, database.DB)
	if err != nil {
		utils.GetResponseJson(c, utils.FAIL, err.Error(), nil)
	} else {

		utils.GetResponseJson(c, utils.SUCCESS, "查询成功", res)
	}
}

// 更新
func (u *UserController) UpdateUser(c *gin.Context) {
	res, err := service.UpdateUser(c, database.DB)
	if err != nil {
		utils.GetResponseJson(c, utils.FAIL, err.Error(), nil)
	} else {
		utils.GetResponseJson(c, utils.SUCCESS, "更新成功", res)
	}
}

// 删除
func (u *UserController) DeleteUser(c *gin.Context) {
	res, err := service.DeleteUser(c, database.DB)
	if err != nil {
		utils.GetResponseJson(c, utils.FAIL, err.Error(), nil)
	} else {
		utils.GetResponseJson(c, utils.SUCCESS, "删除成功", res)
	}
}
