package quality

import (
	"context"
	"fmt"

	"github.com/fustgo/fustgo2/pkg/record"
)

// Validator 数据质量验证器
type Validator struct {
	rules []*ValidationRule
}

// ValidationRule 验证规则
type ValidationRule struct {
	Name        string              `json:"name"`
	Field       string              `json:"field"`
	Type        ValidationType      `json:"type"`
	Parameters  map[string]interface{} `json:"parameters"`
	Description string              `json:"description"`
}

// ValidationType 验证类型
type ValidationType string

const (
	ValidationTypeNotNull    ValidationType = "not_null"
	ValidationTypeMinLength  ValidationType = "min_length"
	ValidationTypeMaxLength  ValidationType = "max_length"
	ValidationTypeMinValue   ValidationType = "min_value"
	ValidationTypeMaxValue   ValidationType = "max_value"
	ValidationTypePattern    ValidationType = "pattern"
	ValidationTypeInList     ValidationType = "in_list"
	ValidationTypeEmail      ValidationType = "email"
	ValidationTypeURL        ValidationType = "url"
	ValidationTypeCustom     ValidationType = "custom"
)

// ValidationResult 验证结果
type ValidationResult struct {
	RecordID  string         `json:"record_id"`
	Field     string         `json:"field"`
	Valid     bool           `json:"valid"`
	Error     string         `json:"error,omitempty"`
	RuleName  string         `json:"rule_name"`
	Timestamp int64          `json:"timestamp"`
}

// NewValidator 创建数据质量验证器实例
func NewValidator() *Validator {
	return &Validator{
		rules: make([]*ValidationRule, 0),
	}
}

// AddRule 添加验证规则
func (v *Validator) AddRule(rule *ValidationRule) error {
	if rule.Name == "" {
		return fmt.Errorf("rule name is required")
	}
	
	if rule.Field == "" {
		return fmt.Errorf("field is required")
	}
	
	if rule.Type == "" {
		return fmt.Errorf("validation type is required")
	}
	
	v.rules = append(v.rules, rule)
	return nil
}

// Validate 验证记录
func (v *Validator) Validate(ctx context.Context, rec *record.Record) ([]*ValidationResult, error) {
	results := make([]*ValidationResult, 0, len(v.rules))
	
	fields := rec.GetFields()
	recordID := rec.GetID()
	
	for _, rule := range v.rules {
		// 检查上下文是否取消
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		
		// 获取字段值
		value, exists := fields[rule.Field]
		if !exists {
			// 字段不存在，根据规则决定是否验证失败
			result := &ValidationResult{
				RecordID:  recordID,
				Field:     rule.Field,
				Valid:     false,
				Error:     "field not found",
				RuleName:  rule.Name,
				Timestamp: rec.GetTimestamp(),
			}
			results = append(results, result)
			continue
		}
		
		// 执行验证
		valid, err := v.validateField(value, rule)
		result := &ValidationResult{
			RecordID:  recordID,
			Field:     rule.Field,
			Valid:     valid,
			RuleName:  rule.Name,
			Timestamp: rec.GetTimestamp(),
		}
		
		if err != nil {
			result.Error = err.Error()
		}
		
		results = append(results, result)
	}
	
	return results, nil
}

// validateField 验证字段
func (v *Validator) validateField(value interface{}, rule *ValidationRule) (bool, error) {
	switch rule.Type {
	case ValidationTypeNotNull:
		return v.validateNotNull(value)
	case ValidationTypeMinLength:
		return v.validateMinLength(value, rule.Parameters)
	case ValidationTypeMaxLength:
		return v.validateMaxLength(value, rule.Parameters)
	case ValidationTypeMinValue:
		return v.validateMinValue(value, rule.Parameters)
	case ValidationTypeMaxValue:
		return v.validateMaxValue(value, rule.Parameters)
	case ValidationTypePattern:
		return v.validatePattern(value, rule.Parameters)
	case ValidationTypeInList:
		return v.validateInList(value, rule.Parameters)
	case ValidationTypeEmail:
		return v.validateEmail(value)
	case ValidationTypeURL:
		return v.validateURL(value)
	case ValidationTypeCustom:
		return v.validateCustom(value, rule.Parameters)
	default:
		return false, fmt.Errorf("unsupported validation type: %s", rule.Type)
	}
}

// validateNotNull 验证非空
func (v *Validator) validateNotNull(value interface{}) (bool, error) {
	if value == nil {
		return false, nil
	}
	
	// 检查字符串是否为空
	if str, ok := value.(string); ok {
		if str == "" {
			return false, nil
		}
	}
	
	return true, nil
}

// validateMinLength 验证最小长度
func (v *Validator) validateMinLength(value interface{}, params map[string]interface{}) (bool, error) {
	minLength, ok := params["min_length"].(int)
	if !ok {
		return false, fmt.Errorf("min_length parameter is required")
	}
	
	switch val := value.(type) {
	case string:
		if len(val) < minLength {
			return false, nil
		}
	case []interface{}:
		if len(val) < minLength {
			return false, nil
		}
	default:
		return false, fmt.Errorf("unsupported type for min_length validation")
	}
	
	return true, nil
}

// validateMaxLength 验证最大长度
func (v *Validator) validateMaxLength(value interface{}, params map[string]interface{}) (bool, error) {
	maxLength, ok := params["max_length"].(int)
	if !ok {
		return false, fmt.Errorf("max_length parameter is required")
	}
	
	switch val := value.(type) {
	case string:
		if len(val) > maxLength {
			return false, nil
		}
	case []interface{}:
		if len(val) > maxLength {
			return false, nil
		}
	default:
		return false, fmt.Errorf("unsupported type for max_length validation")
	}
	
	return true, nil
}

// validateMinValue 验证最小值
func (v *Validator) validateMinValue(value interface{}, params map[string]interface{}) (bool, error) {
	minValue, ok := params["min_value"]
	if !ok {
		return false, fmt.Errorf("min_value parameter is required")
	}
	
	// 比较数值
	switch val := value.(type) {
	case int:
		if min, ok := minValue.(int); ok {
			if val < min {
				return false, nil
			}
		} else if min, ok := minValue.(float64); ok {
			if float64(val) < min {
				return false, nil
			}
		}
	case float64:
		if min, ok := minValue.(float64); ok {
			if val < min {
				return false, nil
			}
		} else if min, ok := minValue.(int); ok {
			if val < float64(min) {
				return false, nil
			}
		}
	default:
		return false, fmt.Errorf("unsupported type for min_value validation")
	}
	
	return true, nil
}

// validateMaxValue 验证最大值
func (v *Validator) validateMaxValue(value interface{}, params map[string]interface{}) (bool, error) {
	maxValue, ok := params["max_value"]
	if !ok {
		return false, fmt.Errorf("max_value parameter is required")
	}
	
	// 比较数值
	switch val := value.(type) {
	case int:
		if max, ok := maxValue.(int); ok {
			if val > max {
				return false, nil
			}
		} else if max, ok := maxValue.(float64); ok {
			if float64(val) > max {
				return false, nil
			}
		}
	case float64:
		if max, ok := maxValue.(float64); ok {
			if val > max {
				return false, nil
			}
		} else if max, ok := maxValue.(int); ok {
			if val > float64(max) {
				return false, nil
			}
		}
	default:
		return false, fmt.Errorf("unsupported type for max_value validation")
	}
	
	return true, nil
}

// validatePattern 验证正则表达式模式
func (v *Validator) validatePattern(value interface{}, params map[string]interface{}) (bool, error) {
	pattern, ok := params["pattern"].(string)
	if !ok {
		return false, fmt.Errorf("pattern parameter is required")
	}
	
	if str, ok := value.(string); ok {
		matched, err := regexpMatch(pattern, str)
		if err != nil {
			return false, fmt.Errorf("invalid pattern: %w", err)
		}
		return matched, nil
	}
	
	return false, fmt.Errorf("unsupported type for pattern validation")
}

// validateInList 验证是否在列表中
func (v *Validator) validateInList(value interface{}, params map[string]interface{}) (bool, error) {
	list, ok := params["list"].([]interface{})
	if !ok {
		return false, fmt.Errorf("list parameter is required")
	}
	
	for _, item := range list {
		if item == value {
			return true, nil
		}
	}
	
	return false, nil
}

// validateEmail 验证邮箱格式
func (v *Validator) validateEmail(value interface{}) (bool, error) {
	if str, ok := value.(string); ok {
		return isEmail(str), nil
	}
	
	return false, fmt.Errorf("unsupported type for email validation")
}

// validateURL 验证URL格式
func (v *Validator) validateURL(value interface{}) (bool, error) {
	if str, ok := value.(string); ok {
		return isURL(str), nil
	}
	
	return false, fmt.Errorf("unsupported type for url validation")
}

// validateCustom 自定义验证
func (v *Validator) validateCustom(value interface{}, params map[string]interface{}) (bool, error) {
	// 自定义验证需要实现具体的验证逻辑
	// 这里简化处理，返回true
	return true, nil
}

// regexpMatch 正则表达式匹配
func regexpMatch(pattern, str string) (bool, error) {
	// TODO: 实现正则表达式匹配
	return true, nil
}

// isEmail 检查是否为有效邮箱
func isEmail(email string) bool {
	// TODO: 实现邮箱格式检查
	return len(email) > 0 && len(email) < 255
}

// isURL 检查是否为有效URL
func isURL(url string) bool {
	// TODO: 实现URL格式检查
	return len(url) > 0 && len(url) < 2048
}

// GetRules 获取所有验证规则
func (v *Validator) GetRules() []*ValidationRule {
	rules := make([]*ValidationRule, len(v.rules))
	copy(rules, v.rules)
	return rules
}

// RemoveRule 移除验证规则
func (v *Validator) RemoveRule(name string) {
	for i, rule := range v.rules {
		if rule.Name == name {
			v.rules = append(v.rules[:i], v.rules[i+1:]...)
			return
		}
	}
}