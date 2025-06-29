package database

import (
	"go-gin/internal/app/config"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 全局变量，用于存储数据库连接
var DB *gorm.DB

// 初始化数据库
func SetupDB(cfg *config.Config) error {
	var err error
	DB, err = gorm.Open(mysql.Open(cfg.Database.Link), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
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
