package router

import (
	"go-gin/internal/controller"
	"go-gin/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes() *gin.Engine {
	router := gin.Default()

	router.SetTrustedProxies([]string{"127.0.0.1"}) // 信任ip
	Routes(router)
	return router
}

// SetupRoutes 设置路由
func Routes(r *gin.Engine) {
	CommonController := new(controller.CommonController)

	// 未匹配到任何路由
	r.NoRoute(CommonController.NoRoute)

	// 定义路由组：v1版
	v1 := r.Group("/api/v1")
	{
		common := v1.Group("/common")
		{
			common.GET("/login", CommonController.Login)
		}

		// 监控权限
		auth := v1.Group("/")
		auth.Use(middleware.AuthMiddleware()) // 使用校验身份中间件
		{
			UserController := new(controller.UserController)

			auth.GET("/user", UserController.GetUser)
			auth.POST("/user/create", UserController.Create)
			auth.POST("/user/update", UserController.UpdateUser)
			auth.POST("/user/delete", UserController.DeleteUser)
		}
	}
}
