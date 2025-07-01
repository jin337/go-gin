package database

import (
	"go-gin/internal/app/config"
	"go-gin/internal/app/logger"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局变量，用于存储数据库连接
var DB *gorm.DB

// 初始化数据库
func SetupDB(cfg *config.Config) error {
	var err error

	// 自定义mysql日志
	sqlLog, err := logger.SetMySqlLogger()
	if err != nil {
		log.Fatalf("创建日志器失败：%v", err)
		return err
	}

	DB, err = gorm.Open(mysql.Open(cfg.Database.Link), &gorm.Config{
		Logger: sqlLog,
	})
	if err != nil {
		log.Fatalf("数据库初始化失败：%v", err)
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("数据库连接失败：%v", err)
		return err
	}
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdle)                       // 设置最大空闲连接数
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpen)                       // 设置最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.MaxLifeTime)) // 设置连接可复用的最大时间

	log.Println("数据库连接成功")
	return nil
}
