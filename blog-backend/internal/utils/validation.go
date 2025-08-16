package utils

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

// 获取验证错误信息
func GetValidationError(err error) string {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			switch fe.Tag() {
			case "required":
				return fe.Field() + "是必需的"
			case "email":
				return fe.Field() + "必须是有效的邮箱地址"
			case "min":
				return fe.Field() + "长度不能少于" + fe.Param() + "个字符"
			default:
				return fe.Field() + "验证失败"
			}
		}
	}
	
	// 处理非验证错误
	if strings.Contains(err.Error(), "json") {
		return "请求数据格式错误"
	}
	
	return "输入参数无效"
}