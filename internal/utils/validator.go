package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

// 全局验证器和翻译器
var (
	validate     *validator.Validate
	trans        ut.Translator
	initOnce     sync.Once
	fieldCache   sync.Map // 字段缓存，提高性能
	jsonTagCache sync.Map // JSON标签缓存
)

type ValidationError struct {
	Errors []string
}

func init() {
	// 延迟初始化，确保只执行一次
	initOnce.Do(initializeValidator)
}

func initializeValidator() {
	// 初始化验证器
	validate = validator.New()

	// 设置中文翻译器
	zhTranslator := zh.New()
	uni := ut.New(zhTranslator, zhTranslator)
	trans, _ = uni.GetTranslator("zh")

	// 注册默认中文翻译
	_ = zh_translations.RegisterDefaultTranslations(validate, trans)

	// 使用缓存优化的字段名翻译
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		// 使用缓存获取JSON标签
		if name, ok := jsonTagCache.Load(fld.Name); ok {
			return name.(string)
		}

		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			jsonTagCache.Store(fld.Name, "")
			return ""
		}

		jsonTagCache.Store(fld.Name, name)
		return name
	})
}

// ValidatorJSON 获取并验证JSON数据
func ValidatorJSON(ctx *gin.Context, data interface{}) error {
	// 读取原始 body
	bodyBytes, _ := io.ReadAll(ctx.Request.Body)
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 重新设置 body

	decoder := json.NewDecoder(bytes.NewReader(bodyBytes))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(data); err != nil {
		return errors.New("包含非法或未知字段")
	}

	// 确保data是指针类型
	if reflect.TypeOf(data).Kind() != reflect.Ptr {
		return errors.New("数据必须是指针")
	}

	// 绑定JSON数据
	if err := ctx.ShouldBindJSON(data); err != nil {
		return handleBindError(err)
	}

	// 验证数据结构
	if err := Validator(data); err != nil {
		return err
	}
	return nil
}

// handleBindError 处理绑定错误
func handleBindError(err error) error {
	errMsg := err.Error()

	switch {
	case strings.Contains(errMsg, "unmarshal"):
		if fieldPart := extractFieldName(errMsg); fieldPart != "" {
			return fmt.Errorf("字段 %s 类型错误", fieldPart)
		}
		return errors.New("请求数据类型错误")
	case strings.Contains(errMsg, "EOF"):
		return errors.New("请求体不能为空")
	case strings.Contains(errMsg, "field"):
		if fieldPart := extractFieldName(errMsg); fieldPart != "" {
			return fmt.Errorf("字段 %s 验证失败", fieldPart)
		}
	}
	return errors.New("无效的请求数据")
}

// extractFieldName 从错误消息中提取字段名
func extractFieldName(errMsg string) string {
	parts := strings.Split(errMsg, "field ")
	if len(parts) < 2 {
		return ""
	}

	fieldPart := strings.Split(parts[1], " ")[0]
	return strings.TrimPrefix(fieldPart, strings.Split(fieldPart, ".")[0]+".")
}

// Validator 验证结构体并返回中文错误信息
func Validator(data interface{}) error {
	// 使用缓存优化反射操作
	val := reflect.ValueOf(data)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 获取类型缓存键
	typeKey := val.Type().String()

	// 检查是否已缓存
	if _, ok := fieldCache.Load(typeKey); !ok {
		cacheFields(val)
	}

	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	// 类型断言获取验证错误
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	// 预分配错误切片
	errMsgs := make([]string, 0, len(validationErrors))
	for _, e := range validationErrors {
		errMsgs = append(errMsgs, e.Translate(trans))
	}

	return &ValidationError{Errors: errMsgs}
}

// cacheFields 缓存结构体字段信息
func cacheFields(val reflect.Value) {
	typeKey := val.Type().String()
	fields := make(map[string]bool, val.NumField())

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fields[field.Name] = true
	}

	fieldCache.Store(typeKey, fields)
}

// 合并所有错误信息为一个字符串
func (e *ValidationError) Error() string {
	return strings.Join(e.Errors, "; ")
}

// GetErrors 获取错误列表
func (e *ValidationError) GetErrors() []string {
	return e.Errors
}
