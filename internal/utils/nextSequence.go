package utils

import (
	"fmt"

	"gorm.io/gorm"
)

func NextSequence(db *gorm.DB, sequenceName string) (uint64, error) {
	var nextVal uint64

	// 在事务中执行
	err := db.Transaction(func(tx *gorm.DB) error {
		// 锁定行（防止并发冲突）
		var seq struct {
			CurrentValue uint64
			Increment    int
		}
		if err := tx.Raw("SELECT current_value, increment FROM sys_sequence WHERE name = ? FOR UPDATE", sequenceName).Scan(&seq).Error; err != nil {
			return err
		}

		// 计算下一个值
		nextVal = seq.CurrentValue + uint64(seq.Increment)

		// 更新序列值
		if err := tx.Exec("UPDATE sys_sequence SET current_value = ? WHERE name = ?", nextVal, sequenceName).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("获取序列失败: %w", err)
	}

	return nextVal, nil
}
