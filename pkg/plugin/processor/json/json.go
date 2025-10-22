package json

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fustgo/fustgo2/pkg/plugin"
	"github.com/fustgo/fustgo2/pkg/record"
)

// JSONProcessor JSON处理器
type JSONProcessor struct {
	plugin.BasePlugin
	config *Config
}

// Config JSON处理器配置
type Config struct {
	Field    string `json:"field" validate:"required"` // 要解析的字段名
	Target   string `json:"target"`                    // 解析后的目标字段名，如果为空则替换原字段
	Remove   bool   `json:"remove"`                    // 是否移除原字段
}

// NewJSONProcessor 创建JSON处理器实例
func NewJSONProcessor() plugin.Processor {
	return &JSONProcessor{
		BasePlugin: plugin.BasePlugin{
			PluginInfo: plugin.Info{
				Name:        "json-processor",
				Type:        plugin.TypeProcessor,
				Version:     "1.0.0",
				Author:      "FustGo Team",
				Description: "JSON数据解析处理器插件",
			},
		},
	}
}

// Init 初始化插件
func (p *JSONProcessor) Init(config map[string]interface{}) error {
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
func (p *JSONProcessor) validateConfig(cfg *Config) error {
	if cfg.Field == "" {
		return fmt.Errorf("field is required")
	}
	return nil
}

// Process 处理数据
func (p *JSONProcessor) Process(ctx context.Context, input <-chan *record.Record, output chan<- *record.Record) error {
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

			// 处理记录
			if err := p.processRecord(rec); err != nil {
				p.GetMetrics().ErrorCount++
				// 错误记录发送到输出通道（根据需求决定是否过滤掉错误记录）
				select {
				case output <- rec:
				case <-ctx.Done():
					return ctx.Err()
				}
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

// processRecord 处理单条记录
func (p *JSONProcessor) processRecord(rec *record.Record) error {
	fields := rec.GetFields()

	// 获取要解析的字段值
	value, exists := fields[p.config.Field]
	if !exists {
		// 字段不存在，不处理
		return nil
	}

	// 将字段值转换为字符串
	var jsonStr string
	switch v := value.(type) {
	case string:
		jsonStr = v
	case []byte:
		jsonStr = string(v)
	default:
		// 其他类型不处理
		return nil
	}

	// 解析JSON
	var parsed interface{}
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// 设置解析后的字段
	targetField := p.config.Target
	if targetField == "" {
		// 如果没有指定目标字段，替换原字段
		targetField = p.config.Field
	}

	// 如果指定了移除原字段且目标字段不同
	if p.config.Remove && targetField != p.config.Field {
		delete(fields, p.config.Field)
	}

	// 设置目标字段
	fields[targetField] = parsed

	return nil
}

// Close 关闭插件
func (p *JSONProcessor) Close() error {
	return nil
}

// Validate 验证配置
func (p *JSONProcessor) Validate() error {
	if p.config == nil {
		return fmt.Errorf("config is nil")
	}
	return p.validateConfig(p.config)
}

// GetInfo 获取插件信息
func (p *JSONProcessor) GetInfo() *plugin.Info {
	return &p.PluginInfo
}

// GetName 获取插件名称
func (p *JSONProcessor) GetName() string {
	return p.PluginInfo.Name
}

// GetType 获取插件类型
func (p *JSONProcessor) GetType() plugin.Type {
	return p.PluginInfo.Type
}