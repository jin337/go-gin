package task

import (
	"log"
	"time"

	"gorm.io/gorm"
)

// 每天检查是否为每月第一天，并清理 BlacklistedToken 表
func CleanupBlacklistedToken(db *gorm.DB) {
	go func() {
		ticker := time.NewTicker(24 * time.Hour) // 设定触发间隔
		defer ticker.Stop()

		for range ticker.C {
			now := time.Now()
			// 判断是否为每月第一天 每天凌晨0点执行
			if now.Day() == 1 && now.Hour() == 0 {
				log.Println("开始清理BlacklistedToken表...")
				if err := db.Exec("TRUNCATE TABLE blacklisted_tokens").Error; err != nil {
					log.Printf("BlacklistedToken表清理失败: %v\n", err)
				} else {
					log.Println("BlacklistedToken表清理完成")
				}
			}
		}
	}()
}
