package controller

import (
	"go-gin/internal/app/database"
	"go-gin/internal/service"
	"go-gin/internal/utils"

	"github.com/gin-gonic/gin"
)

type CommonController struct{}

// 处理 404 请求
func (c *CommonController) NoRoute(ctx *gin.Context) {
	utils.GetResponseJson(ctx, utils.NOT_FOUND, "Not Found", nil)
}

// 登录
func (c *CommonController) Login(ctx *gin.Context) {
	res, err := service.Login(ctx, database.DB)
	if err != nil {
		utils.GetResponseJson(ctx, utils.FAIL, err.Error(), nil)
	} else {

		utils.GetResponseJson(ctx, utils.SUCCESS, "登录成功", res)
	}
}
