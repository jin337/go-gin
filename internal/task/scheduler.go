package task

import (
	"gorm.io/gorm"
)

// 启动所有定时任务
func StartAllTasks(db *gorm.DB) {
	go CleanupBlacklistedToken(db) // 清理token黑名单
}
