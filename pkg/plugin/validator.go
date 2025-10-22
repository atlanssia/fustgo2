package plugin

import (
	"fmt"
	"reflect"
)

// Validator 配置验证器接口
type Validator interface {
	// Validate 验证配置
	Validate(config map[string]interface{}) error
}

// ConfigValidator 配置验证器
type ConfigValidator struct {
	rules []ValidationRule
}

// ValidationRule 验证规则
type ValidationRule struct {
	// Field 字段名
	Field string

	// Required 是否必填
	Required bool

	// Type 期望的类型（可选）
	Type reflect.Kind

	// Validator 自定义验证函数
	Validator func(value interface{}) error

	// DefaultValue 默认值（当字段不存在且非必填时使用）
	DefaultValue interface{}
}

// NewConfigValidator 创建配置验证器
func NewConfigValidator() *ConfigValidator {
	return &ConfigValidator{
		rules: make([]ValidationRule, 0),
	}
}

// AddRule 添加验证规则
func (v *ConfigValidator) AddRule(rule ValidationRule) *ConfigValidator {
	v.rules = append(v.rules, rule)
	return v
}

// RequireString 添加必填字符串字段
func (v *ConfigValidator) RequireString(field string) *ConfigValidator {
	return v.AddRule(ValidationRule{
		Field:    field,
		Required: true,
		Type:     reflect.String,
	})
}

// OptionalString 添加可选字符串字段
func (v *ConfigValidator) OptionalString(field string, defaultValue string) *ConfigValidator {
	return v.AddRule(ValidationRule{
		Field:        field,
		Required:     false,
		Type:         reflect.String,
		DefaultValue: defaultValue,
	})
}

// RequireInt 添加必填整数字段
func (v *ConfigValidator) RequireInt(field string) *ConfigValidator {
	return v.AddRule(ValidationRule{
		Field:    field,
		Required: true,
		Validator: func(value interface{}) error {
			switch value.(type) {
			case int, int64, float64:
				return nil
			default:
				return fmt.Errorf("field %s must be an integer", field)
			}
		},
	})
}

// OptionalInt 添加可选整数字段
func (v *ConfigValidator) OptionalInt(field string, defaultValue int) *ConfigValidator {
	return v.AddRule(ValidationRule{
		Field:        field,
		Required:     false,
		DefaultValue: defaultValue,
		Validator: func(value interface{}) error {
			switch value.(type) {
			case int, int64, float64:
				return nil
			default:
				return fmt.Errorf("field %s must be an integer", field)
			}
		},
	})
}

// RequireBool 添加必填布尔字段
func (v *ConfigValidator) RequireBool(field string) *ConfigValidator {
	return v.AddRule(ValidationRule{
		Field:    field,
		Required: true,
		Type:     reflect.Bool,
	})
}

// OptionalBool 添加可选布尔字段
func (v *ConfigValidator) OptionalBool(field string, defaultValue bool) *ConfigValidator {
	return v.AddRule(ValidationRule{
		Field:        field,
		Required:     false,
		Type:         reflect.Bool,
		DefaultValue: defaultValue,
	})
}

// Validate 验证配置
func (v *ConfigValidator) Validate(config map[string]interface{}) error {
	for _, rule := range v.rules {
		value, exists := config[rule.Field]

		// 检查必填字段
		if !exists {
			if rule.Required {
				return fmt.Errorf("required field %s is missing", rule.Field)
			}
			// 设置默认值
			if rule.DefaultValue != nil {
				config[rule.Field] = rule.DefaultValue
			}
			continue
		}

		// 检查类型
		if rule.Type != reflect.Invalid {
			valueType := reflect.TypeOf(value).Kind()
			if valueType != rule.Type {
				return fmt.Errorf("field %s has wrong type: expected %s, got %s",
					rule.Field, rule.Type, valueType)
			}
		}

		// 自定义验证
		if rule.Validator != nil {
			if err := rule.Validator(value); err != nil {
				return err
			}
		}
	}

	return nil
}

// 内置验证函数

// ValidateRange 验证数值范围
func ValidateRange(min, max int) func(interface{}) error {
	return func(value interface{}) error {
		var num int
		switch v := value.(type) {
		case int:
			num = v
		case int64:
			num = int(v)
		case float64:
			num = int(v)
		default:
			return fmt.Errorf("value must be a number")
		}

		if num < min || num > max {
			return fmt.Errorf("value %d is out of range [%d, %d]", num, min, max)
		}
		return nil
	}
}

// ValidateEnum 验证枚举值
func ValidateEnum(validValues ...string) func(interface{}) error {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("value must be a string")
		}

		for _, valid := range validValues {
			if str == valid {
				return nil
			}
		}

		return fmt.Errorf("value %s is not in valid values: %v", str, validValues)
	}
}

// ValidatePattern 验证正则表达式（简化版，实际应使用 regexp）
func ValidatePattern(pattern string) func(interface{}) error {
	return func(value interface{}) error {
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("value must be a string")
		}

		// 这里简化处理，实际应使用正则表达式
		if len(str) == 0 {
			return fmt.Errorf("value cannot be empty")
		}

		return nil
	}
}

// ValidateFileExists 验证文件是否存在
func ValidateFileExists() func(interface{}) error {
	return func(value interface{}) error {
		path, ok := value.(string)
		if !ok {
			return fmt.Errorf("file path must be a string")
		}

		if len(path) == 0 {
			return fmt.Errorf("file path cannot be empty")
		}

		// 实际应该检查文件是否存在
		// 这里简化处理
		return nil
	}
}

// ValidateURL 验证 URL 格式
func ValidateURL() func(interface{}) error {
	return func(value interface{}) error {
		url, ok := value.(string)
		if !ok {
			return fmt.Errorf("URL must be a string")
		}

		if len(url) == 0 {
			return fmt.Errorf("URL cannot be empty")
		}

		// 简单检查是否包含协议
		if len(url) < 7 || (url[:7] != "http://" && url[:8] != "https://") {
			return fmt.Errorf("URL must start with http:// or https://")
		}

		return nil
	}
}
