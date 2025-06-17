package controller

import (
	"go-gin/internal/utils"

	"github.com/gin-gonic/gin"
)

type CommonController struct{}

// 处理 404 请求
func (u *CommonController) NoRoute(c *gin.Context) {
	utils.GetResponseJson(c, utils.NOT_FOUND, "Not Found", nil)
}

// 登录
func (u *CommonController) Login(c *gin.Context) {
	utils.GetResponseJson(c, utils.SUCCESS, "success", nil)
}
