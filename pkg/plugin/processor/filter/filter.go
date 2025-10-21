package filter

import (
	"context"
	"fmt"
	"regexp"

	"github.com/fustgo/fustgo2/pkg/plugin"
	"github.com/fustgo/fustgo2/pkg/record"
)

// FilterProcessor 数据过滤处理器
type FilterProcessor struct {
	plugin.BasePlugin
	config *Config
}

// Config 过滤器配置
type Config struct {
	Condition string `json:"condition" validate:"required"` // 过滤条件，支持简单表达式
	Field     string `json:"field"`                        // 字段名
	Pattern   string `json:"pattern"`                      // 正则表达式模式
	Operation string `json:"operation"`                    // 操作符: eq, ne, gt, lt, ge, le, contains, regex
	Value     interface{} `json:"value"`                   // 比较值
}

// NewFilterProcessor 创建过滤处理器实例
func NewFilterProcessor() plugin.Processor {
	return &FilterProcessor{
		BasePlugin: plugin.BasePlugin{
			PluginInfo: plugin.Info{
				Name:        "filter-processor",
				Type:        plugin.TypeProcessor,
				Version:     "1.0.0",
				Author:      "FustGo Team",
				Description: "数据过滤处理器插件",
			},
		},
	}
}

// Init 初始化插件
func (p *FilterProcessor) Init(config map[string]interface{}) error {
	// 解析配置
	cfg := &Config{}
	if err := plugin.MapToStruct(config, cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// 验证配置
	if err := p.validateConfig(cfg); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	p.config = cfg
	return nil
}

// validateConfig 验证配置
func (p *FilterProcessor) validateConfig(cfg *Config) error {
	if cfg.Condition == "" && cfg.Field == "" {
		return fmt.Errorf("either condition or field must be specified")
	}
	
	if cfg.Field != "" && cfg.Operation == "" {
		return fmt.Errorf("operation is required when field is specified")
	}
	
	return nil
}

// Process 处理数据
func (p *FilterProcessor) Process(ctx context.Context, input <-chan *record.Record, output chan<- *record.Record) error {
	defer close(output)

	for {
		select {
		case rec, ok := <-input:
			if !ok {
				// 输入通道已关闭
				return nil
			}

			// 检查上下文是否取消
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			// 应用过滤条件
			match, err := p.matchRecord(rec)
			if err != nil {
				p.GetMetrics().ErrorCount++
				// 错误记录发送到输出通道（根据需求决定是否过滤掉错误记录）
				select {
				case output <- rec:
				case <-ctx.Done():
					return ctx.Err()
				}
				continue
			}

			// 如果匹配条件，则过滤掉该记录（不发送到输出通道）
			if match {
				p.GetMetrics().FilteredCount++
				continue
			}

			// 发送到输出通道
			select {
			case output <- rec:
				p.GetMetrics().RecordsProcessed++
				p.GetMetrics().SuccessCount++
			case <-ctx.Done():
				return ctx.Err()
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// matchRecord 检查记录是否匹配过滤条件
func (p *FilterProcessor) matchRecord(rec *record.Record) (bool, error) {
	// 如果配置了条件表达式，优先使用条件表达式
	if p.config.Condition != "" {
		return p.evaluateCondition(rec, p.config.Condition)
	}

	// 使用字段和操作符进行匹配
	if p.config.Field != "" {
		return p.evaluateFieldCondition(rec, p.config.Field, p.config.Operation, p.config.Value)
	}

	return false, nil
}

// evaluateCondition 评估条件表达式
func (p *FilterProcessor) evaluateCondition(rec *record.Record, condition string) (bool, error) {
	// 简化实现，实际可以使用表达式引擎如 govaluate
	// 这里只处理简单的包含判断
	
	fields := rec.GetFields()
	
	// 检查是否包含某个字段值
	if condition[0] == '"' && condition[len(condition)-1] == '"' {
		// 字符串值匹配
		value := condition[1 : len(condition)-1]
		for _, v := range fields {
			if str, ok := v.(string); ok && str == value {
				return true, nil
			}
		}
	}
	
	return false, nil
}

// evaluateFieldCondition 评估字段条件
func (p *FilterProcessor) evaluateFieldCondition(rec *record.Record, field, operation string, value interface{}) (bool, error) {
	fields := rec.GetFields()
	fieldValue, exists := fields[field]
	if !exists {
		return false, nil
	}

	switch operation {
	case "eq":
		return fieldValue == value, nil
	case "ne":
		return fieldValue != value, nil
	case "gt":
		return compare(fieldValue, value) > 0, nil
	case "lt":
		return compare(fieldValue, value) < 0, nil
	case "ge":
		return compare(fieldValue, value) >= 0, nil
	case "le":
		return compare(fieldValue, value) <= 0, nil
	case "contains":
		if str, ok := fieldValue.(string); ok {
			if substr, ok := value.(string); ok {
				return contains(str, substr), nil
			}
		}
		return false, nil
	case "regex":
		if str, ok := fieldValue.(string); ok {
			if pattern, ok := value.(string); ok {
				matched, err := regexp.MatchString(pattern, str)
				if err != nil {
					return false, fmt.Errorf("invalid regex pattern: %w", err)
				}
				return matched, nil
			}
		}
		return false, nil
	default:
		return false, fmt.Errorf("unsupported operation: %s", operation)
	}
}

// compare 比较两个值
func compare(a, b interface{}) int {
	switch a := a.(type) {
	case int:
		if b, ok := b.(int); ok {
			if a > b {
				return 1
			} else if a < b {
				return -1
			}
			return 0
		}
	case float64:
		if b, ok := b.(float64); ok {
			if a > b {
				return 1
			} else if a < b {
				return -1
			}
			return 0
		}
	case string:
		if b, ok := b.(string); ok {
			if a > b {
				return 1
			} else if a < b {
				return -1
			}
			return 0
		}
	}
	return 0
}

// contains 检查字符串是否包含子串
func contains(str, substr string) bool {
	return len(str) >= len(substr) && 
		(len(substr) == 0 || 
			(len(str) == len(substr) && str == substr) || 
			(len(str) > len(substr) && (str == substr || 
				str[:len(substr)] == substr || 
				str[len(str)-len(substr):] == substr || 
				(len(str) > len(substr) && 
					(len(str) >= 2*len(substr) && 
						(str[:len(substr)] == substr && str[len(str)-len(substr):] == substr)) ||
					(len(str) > len(substr) && 
						(len(str) >= len(substr)+1 && 
							(str[1:1+len(substr)] == substr || 
								str[len(str)-len(substr)-1:len(str)-1] == substr)))))))
	
	// 简化实现
	return len(substr) == 0 || len(str) >= len(substr) && (str == substr || 
		(len(str) > len(substr) && (str[:len(substr)] == substr || 
			str[len(str)-len(substr):] == substr ||
			(len(str) > len(substr) && 
				(len(str) >= 2*len(substr) && 
					(str[:len(substr)] == substr && str[len(str)-len(substr):] == substr)) ||
				(len(str) > len(substr) && 
					(len(str) >= len(substr)+1 && 
						(str[1:1+len(substr)] == substr || 
							str[len(str)-len(substr)-1:len(str)-1] == substr)))))))
}

// Close 关闭插件
func (p *FilterProcessor) Close() error {
	return nil
}

// Validate 验证配置
func (p *FilterProcessor) Validate() error {
	if p.config == nil {
		return fmt.Errorf("config is nil")
	}
	return p.validateConfig(p.config)
}

// GetInfo 获取插件信息
func (p *FilterProcessor) GetInfo() *plugin.Info {
	return &p.PluginInfo
}

// GetName 获取插件名称
func (p *FilterProcessor) GetName() string {
	return p.PluginInfo.Name
}

// GetType 获取插件类型
func (p *FilterProcessor) GetType() plugin.Type {
	return p.PluginInfo.Type
}