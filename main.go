package main

import (
	"go-gin/internal/app/config"
	"go-gin/internal/app/database"
	"go-gin/internal/router"
	"log"
)

func main() {
	// 初始化配置
	cfg, err := config.SetupConfig()
	if err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	// 初始化数据库
	database.SetupDB(cfg)

	// 设置路由
	r := router.SetupRoutes()

	// 启动服务
	if err := r.Run(cfg.Service.Port); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}
