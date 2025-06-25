package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-sql-driver/mysql"
)

func TranslateDBError(err error) error {
	var mysqlErr *mysql.MySQLError
	if !errors.As(err, &mysqlErr) {
		return err // 非 MySQL 错误，直接返回
	}
	switch mysqlErr.Number {
	case 1062:
		field := extractFieldFromKey(mysqlErr.Message)
		return fmt.Errorf("%s字段数据已存在", field)
	case 1064:
		return errors.New("SQL 语法错误，请检查查询语句")
	case 1146:
		return errors.New("表不存在，请检查数据库结构")
	case 1054:
		return errors.New("列不存在，请检查字段名")
	case 1045:
		return errors.New("数据库访问被拒绝，请检查用户名和密码")
	case 1452:
		return errors.New("外键约束失败，关联数据不存在")
	case 1213:
		return errors.New("事务死锁，请稍后重试")
	case 1205:
		return errors.New("锁等待超时，请减少并发操作")
	case 1366:
		field := extractFieldFromKey(mysqlErr.Message)
		return fmt.Errorf("%s数据类型错误，请检查输入值", field)
	case 1264:
		return errors.New("数值超出范围，请调整输入值")
	case 2002:
		return errors.New("无法连接数据库，请检查 MySQL 服务")
	default:
		return errors.New("数据库操作失败: " + mysqlErr.Message)
	}
}

func extractFieldFromKey(errorMsg string) string {
	if !strings.Contains(errorMsg, "for key") {
		return "该"
	}

	keyPart := strings.Split(errorMsg, "for key")[1]
	keyPart = strings.Trim(keyPart, " '") // 移除引号和空格

	switch {
	// 情况1：表名.前缀_表名_字段名 → 字段名
	case strings.Contains(keyPart, "."):
		re := regexp.MustCompile(`^[^.]+\.[^_]+_[^_]+_(.+)$`)
		matches := re.FindStringSubmatch(keyPart)
		if len(matches) > 1 {
			return matches[1]
		}
	// 直接返回字段名
	default:
		return keyPart
	}

	return "该"
}
