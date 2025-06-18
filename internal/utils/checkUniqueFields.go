package utils

import (
	"fmt"

	"gorm.io/gorm"
)

// CheckUniqueFields 检查指定字段是否已存在记录（支持多个字段）
// fields: map[string]interface{}{"field_name": value}
func CheckUniqueFields(db *gorm.DB, model interface{}, fields map[string]interface{}) error {
	for field, value := range fields {
		var count int64
		query := db.Model(model).Where(fmt.Sprintf("%s = ?", field), value).Count(&count)
		if query.Error != nil {
			return fmt.Errorf("数据库查询错误：%v", query.Error)
		}
		if count > 0 {
			return fmt.Errorf("字段%s的值%v已存在", field, value)
		}
	}
	return nil
}
