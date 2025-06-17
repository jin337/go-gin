package main

import (
	"fmt"
	"go-gin/internal/app/config"
	"go-gin/internal/app/database"
	"go-gin/internal/router"
)

func main() {
	// 初始化配置
	cfg, err := config.SetupConfig()
	if err != nil {
		fmt.Println("初始化配置失败：", err)
	}

	// 初始化数据库
	database.SetupDB(cfg)

	// 设置路由
	r := router.SetupRoutes()

	// 启动服务
	if err := r.Run(cfg.Service.Port); err != nil {
		fmt.Println("启动服务失败：", err)
	}
}
