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
func SetupDB() (*gorm.DB, error) {
	var err error
	Database := config.GetGlobalConfig().Database

	// 自定义mysql日志
	sqlLog, err := logger.SetMySqlLogger()
	if err != nil {
		log.Fatalf("创建日志器失败：%v", err)
		return nil, err
	}
	// 创建数据库连接
	DB, err = gorm.Open(mysql.Open(Database.Link), &gorm.Config{
		Logger: sqlLog,
	})
	if err != nil {
		log.Fatalf("数据库初始化失败：%v", err)
		return nil, err
	}

	// 设置数据库连接参数
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("数据库连接失败：%v", err)
		return nil, err
	}
	sqlDB.SetMaxIdleConns(Database.MaxIdle)                       // 设置最大空闲连接数
	sqlDB.SetMaxOpenConns(Database.MaxOpen)                       // 设置最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Duration(Database.MaxLifeTime)) // 设置连接可复用的最大时间

	log.Printf("数据库连接成功")
	return DB, nil
}
