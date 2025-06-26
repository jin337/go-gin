package main

import (
	"go-gin/internal/app/config"
	"go-gin/internal/app/database"
	"go-gin/internal/router"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	cfg, err := config.SetupConfig()
	if err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	// 初始化数据库
	if err := database.SetupDB(cfg); err != nil {
		log.Fatalf("数据库配置失败: %v", err)
	}
	// 设置Gin模式:DebugMode/ReleaseMode
	gin.SetMode(gin.ReleaseMode)

	// 初始化服务
	if err := router.SetupRoutes(cfg); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}
