package processor

import (
	"context"
	"fmt"

	"github.com/fustgo/fustgo2/pkg/plugin"
	"github.com/fustgo/fustgo2/pkg/record"
)

// FieldMapper 字段映射处理器
// 将输入字段映射到输出字段
type FieldMapper struct {
	*plugin.BasePlugin
	metrics *plugin.BaseMetrics

	// 配置
	mappings     map[string]string // 源字段 -> 目标字段
	dropUnmapped bool              // 是否丢弃未映射的字段
}

// NewFieldMapper 创建字段映射处理器
func NewFieldMapper() plugin.Plugin {
	return &FieldMapper{
		BasePlugin:   plugin.NewBasePlugin("field_mapper", plugin.TypeProcessor),
		metrics:      plugin.NewBaseMetrics(),
		mappings:     make(map[string]string),
		dropUnmapped: false,
	}
}

// Init 初始化插件
func (p *FieldMapper) Init(config map[string]interface{}) error {
	if err := p.BasePlugin.Init(config); err != nil {
		return err
	}

	// 获取映射配置
	if mappings, ok := config["mappings"].(map[string]interface{}); ok {
		for src, dst := range mappings {
			if dstStr, ok := dst.(string); ok {
				p.mappings[src] = dstStr
			}
		}
	}

	// 是否丢弃未映射字段
	if drop, ok := p.GetConfigBool("drop_unmapped"); ok {
		p.dropUnmapped = drop
	}

	return nil
}

// Validate 验证配置
func (p *FieldMapper) Validate() error {
	if len(p.mappings) == 0 {
		return fmt.Errorf("mappings configuration is required")
	}
	return nil
}

// Process 处理数据
func (p *FieldMapper) Process(ctx context.Context, input <-chan *record.Record, output chan<- *record.Record) error {
	defer close(output)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case rec, ok := <-input:
			if !ok {
				// 输入通道已关闭
				return nil
			}

			// 映射字段
			newData := make(map[string]interface{})

			// 处理映射的字段
			for srcField, dstField := range p.mappings {
				if value, ok := rec.Data[srcField]; ok {
					newData[dstField] = value
				}
			}

			// 如果不丢弃未映射字段，保留原始字段
			if !p.dropUnmapped {
				for key, value := range rec.Data {
					// 如果字段未被映射，保留原字段名
					if _, mapped := p.mappings[key]; !mapped {
						if _, exists := newData[key]; !exists {
							newData[key] = value
						}
					}
				}
			}

			// 创建新记录
			newRec := record.NewRecordWithMeta(newData, rec.Meta)
			newRec.Timestamp = rec.Timestamp

			// 发送记录
			select {
			case output <- newRec:
				p.metrics.IncrRecordsProcessed(1)
				p.metrics.IncrSuccessCount(1)
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
}

// GetMetrics 获取指标
func (p *FieldMapper) GetMetrics() *plugin.Metrics {
	return p.metrics.GetMetrics()
}

// Close 关闭插件
func (p *FieldMapper) Close() error {
	return nil
}

// 注册插件
func init() {
	info := &plugin.Info{
		Name:        "field_mapper",
		Type:        plugin.TypeProcessor,
		Version:     "1.0.0",
		Author:      "FustGo Team",
		Description: "Map input fields to output fields",
		ConfigSchema: map[string]interface{}{
			"mappings": map[string]interface{}{
				"description": "Field mappings (source -> destination)",
				"example": map[string]interface{}{
					"old_name": "new_name",
					"user_id":  "id",
				},
			},
			"drop_unmapped": "bool (optional) - Drop unmapped fields. Default: false",
		},
	}

	plugin.Register("field_mapper", NewFieldMapper, info)
}
