package router

import (
	"go-gin/internal/app/config"
	"go-gin/internal/controller"
	"go-gin/internal/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 初始化路由并启动服务
func SetupRoutes() error {
	r := gin.Default()

	r.Use(middleware.LoggerMiddleware())       // 中间件-日志
	r.SetTrustedProxies([]string{"127.0.0.1"}) // 信任ip

	// 注册 API 路由
	RegisterAPIRoutes(r)

	// 启动服务
	Port := config.GetGlobalConfig().Service.Port
	log.Printf("运行端口:%s", Port)
	return r.Run(":" + Port)
}

// SetupRoutes 设置路由
func RegisterAPIRoutes(r *gin.Engine) {
	CommonController := new(controller.CommonController)

	// 未匹配到任何路由
	r.NoRoute(CommonController.NoRoute)

	// 定义路由组：v1版
	v1 := r.Group("/api/v1")
	{
		common := v1.Group("/common")
		{
			common.POST("/login", CommonController.Login)
			common.POST("/logout", CommonController.Logout)
		}

		// 监控权限
		auth := v1.Group("/", middleware.AuthMiddleware())
		{
			UserController := new(controller.UserController)
			auth.POST("/user", UserController.GetUser)
			auth.POST("/user/create", UserController.CreateUser)
			auth.POST("/user/update", UserController.UpdateUser)
			auth.POST("/user/delete", UserController.DeleteUser)

			AccountController := new(controller.AccountController)
			auth.POST("/account", AccountController.GetAccount)
			auth.POST("/account/create", AccountController.CreateAccount)
			auth.POST("/account/update", AccountController.UpdateAccount)
			auth.POST("/account/delete", AccountController.DeleteAccount)
		}
	}
}
