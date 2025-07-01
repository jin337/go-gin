package main

import (
	"go-gin/internal/app/config"
	"go-gin/internal/app/database"
	"go-gin/internal/app/logger"
	"go-gin/internal/router"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 创建日志文件
	if err := logger.SetupLog(); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}

	// 初始化配置
	if err := config.SetupConfig(); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	// 初始化数据库
	if err := database.SetupDB(); err != nil {
		log.Fatalf("数据库配置失败: %v", err)
	}

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 初始化服务
	if err := router.SetupRoutes(); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}
