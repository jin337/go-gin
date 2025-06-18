package utils

import (
	"fmt"

	"gorm.io/gorm"
)

// CheckUniqueFields 检查指定模型中字段值的唯一性。
// 参数:
//
//	db (*gorm.DB): 数据库连接对象。
//	model (interface{}): 需要检查的模型实例。
//	fields (map[string]interface{}): 包含需要检查的字段名和对应值的映射,如：{"field_name": value}
//
// 返回值:
//
//	error: 如果检查过程中发生错误或找到重复的字段值，则返回错误。
func CheckUniqueFields(db *gorm.DB, model interface{}, fields map[string]interface{}) error {
	// 遍历需要检查的字段和它们的值。
	for field, value := range fields {
		var count int64
		// 执行数据库查询，统计符合条件的记录数。
		query := db.Model(model).Where(fmt.Sprintf("%s = ?", field), value).Count(&count)
		// 如果查询过程中发生错误，返回错误信息。
		if query.Error != nil {
			return fmt.Errorf("数据库查询错误：%v", query.Error)
		}
		// 如果记录数大于0，说明字段值不唯一，返回错误信息。
		if count > 0 {
			return fmt.Errorf("字段%s的值%v已存在", field, value)
		}
	}
	// 所有字段值都是唯一的，返回nil表示检查通过。
	return nil
}
