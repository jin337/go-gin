package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

// 全局验证器和翻译器
var (
	validate *validator.Validate
	trans    ut.Translator
)

type ValidationError struct {
	Errors []string
}

func init() {
	// 初始化验证器
	validate = validator.New()

	// 设置中文翻译器
	zh := zh.New()
	uni := ut.New(zh, zh)
	trans, _ = uni.GetTranslator("zh")

	// 注册默认中文翻译
	_ = zh_translations.RegisterDefaultTranslations(validate, trans)

	// 注册自定义字段名翻译（使用json标签作为字段名）
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Validator 验证结构体并返回中文错误信息
func Validator(data interface{}) error {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}
	// 类型断言获取验证错误
	validationErrors, ok := err.(validator.ValidationErrors)

	if !ok {
		return err
	}

	// 收集所有错误信息
	var errMsgs []string
	for _, e := range validationErrors {
		errMsgs = append(errMsgs, e.Translate(trans))
	}

	return &ValidationError{
		Errors: errMsgs,
	}
}

// 合并所有错误信息为一个字符串
func (e *ValidationError) Error() string {
	return strings.Join(e.Errors, "; ")
}

// GetErrors 获取错误列表
func (e *ValidationError) GetErrors() []string {
	return e.Errors
}

// 获取json,并验证类型
func ValidatorJSON(ctx *gin.Context, data interface{}) error {
	if err := ctx.ShouldBindJSON(&data); err != nil {
		if _, ok := err.(validator.ValidationErrors); !ok {
			errMsg := err.Error()
			if strings.Contains(errMsg, "unmarshal") {
				parts := strings.Split(errMsg, "field ")
				if len(parts) > 1 {
					fieldPart := strings.Split(parts[1], " ")[0]
					fieldName := strings.TrimPrefix(fieldPart, strings.Split(fieldPart, ".")[0]+".")
					return fmt.Errorf("字段 %s 类型错误", fieldName)
				}
			} else if strings.Contains(errMsg, "EOF") {
				return errors.New("请求体不能为空")
			} else {
				return errors.New("无效的请求数据")
			}
		}
		return err
	}
	fmt.Printf("data: %+v\n", data)
	// // 忽略 nil 指针字段
	OmitNilFields(validate, data)
	if err := Validator(data); err != nil {
		return err
	}
	return nil
}

// 动态忽略所有为 nil 的指针字段
func OmitNilFields(v *validator.Validate, data interface{}) {
	val := reflect.ValueOf(data).Elem()

	for i := 0; i < val.NumField(); i++ {
		structField := val.Type().Field(i)
		field := val.Field(i)

		// 如果字段是指针类型且为 nil，则跳过验证
		if field.Kind() == reflect.Ptr && field.IsNil() {
			v.RegisterValidation(structField.Name, func(fl validator.FieldLevel) bool {
				return true // 跳过校验
			}, true)
		}
	}
}
