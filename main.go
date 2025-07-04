package main

import (
	"flag"
	"go-gin/internal/app/config"
	"go-gin/internal/app/database"
	"go-gin/internal/app/logger"
	"go-gin/internal/model"
	"go-gin/internal/router"
	"go-gin/internal/task"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 定义命令行参数：参数名，默认值，参数描述
	env := flag.String("env", "dev", "构建环境")
	flag.Parse()
	// 初始化配置
	if err := config.SetupConfig(*env); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	// 创建日志文件
	if err := logger.SetupLog(); err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}

	// 初始化数据库
	db, err := database.SetupDB()
	if err != nil {
		log.Fatalf("数据库配置失败: %v", err)
	}

	// 启动所有定时任务
	task.StartAllTasks(db)

	// 自动创建或更新表和字段
	if config.GetGlobalConfig().Database.MigrateTables {
		log.Println("执行表结构迁移...")
		db.AutoMigrate(
			&model.User{},
			&model.Account{},
			&model.BlacklistedToken{},
		)
	}

	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	// 初始化服务
	if err := router.SetupRoutes(); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}
