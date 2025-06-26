package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/go-sql-driver/mysql"
)

// 预编译正则表达式，避免每次调用都重新编译
var (
	keyRegex  *regexp.Regexp
	regexOnce sync.Once
)

func init() {
	// 延迟初始化正则表达式
	regexOnce.Do(func() {
		keyRegex = regexp.MustCompile(`^[^.]+\.[^_]+_[^_]+_(.+)$`)
	})
}

func TranslateDBError(err error) error {
	var mysqlErr *mysql.MySQLError
	if !errors.As(err, &mysqlErr) {
		return err // 非 MySQL 错误，直接返回
	}

	// 类型判断
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
		// 避免每次都创建新字符串
		return fmt.Errorf("数据库操作失败: %s", mysqlErr.Message)
	}
}

// 减少字符串分配
var stringPool = sync.Pool{
	New: func() interface{} {
		return new(strings.Builder)
	},
}

// 获取字段名
func extractFieldFromKey(errorMsg string) string {
	// 快速检查是否包含关键词，避免不必要的处理
	if !strings.Contains(errorMsg, "for key") {
		return "该"
	}

	// 使用 strings.Cut 代替 Split 提高性能
	_, after, found := strings.Cut(errorMsg, "for key")
	if !found {
		return "该"
	}

	// 从池中获取 Builder 减少内存分配
	builder := stringPool.Get().(*strings.Builder)
	defer func() {
		builder.Reset()
		stringPool.Put(builder)
	}()

	// 直接写入 Builder 避免中间字符串
	for _, r := range after {
		if r == '\'' || r == ' ' {
			continue
		}
		builder.WriteRune(r)
	}
	keyPart := builder.String()

	// 表名.前缀_表名_字段名 → 字段名
	if strings.Contains(keyPart, ".") {
		matches := keyRegex.FindStringSubmatch(keyPart)
		if len(matches) > 1 {
			return matches[1]
		}
	} else {
		// 直接返回字段名
		return keyPart
	}

	return "该"
}
