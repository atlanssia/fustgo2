package processor

import (
	"context"
	"fmt"
	"strings"

	"github.com/fustgo/fustgo2/pkg/plugin"
	"github.com/fustgo/fustgo2/pkg/record"
)

// Filter 过滤器处理器
// 根据条件过滤记录
type Filter struct {
	*plugin.BasePlugin
	metrics *plugin.BaseMetrics

	// 配置
	condition string            // 过滤条件（简化版）
	rules     []FilterRule      // 过滤规则
	mode      string            // include 或 exclude
}

// FilterRule 过滤规则
type FilterRule struct {
	Field    string
	Operator string // eq, ne, gt, lt, contains, startswith, endswith
	Value    interface{}
}

// NewFilter 创建过滤器
func NewFilter() plugin.Plugin {
	return &Filter{
		BasePlugin: plugin.NewBasePlugin("filter", plugin.TypeProcessor),
		metrics:    plugin.NewBaseMetrics(),
		mode:       "include",
	}
}

// Init 初始化插件
func (p *Filter) Init(config map[string]interface{}) error {
	if err := p.BasePlugin.Init(config); err != nil {
		return err
	}

	// 获取模式
	if mode, ok := p.GetConfigString("mode"); ok {
		p.mode = mode
	}

	// 获取条件
	p.condition, _ = p.GetConfigString("condition")

	// 解析规则
	if rules, ok := config["rules"].([]interface{}); ok {
		p.rules = make([]FilterRule, 0, len(rules))
		for _, r := range rules {
			if ruleMap, ok := r.(map[string]interface{}); ok {
				rule := FilterRule{
					Field:    ruleMap["field"].(string),
					Operator: ruleMap["operator"].(string),
					Value:    ruleMap["value"],
				}
				p.rules = append(p.rules, rule)
			}
		}
	}

	return nil
}

// Validate 验证配置
func (p *Filter) Validate() error {
	if p.mode != "include" && p.mode != "exclude" {
		return fmt.Errorf("mode must be 'include' or 'exclude'")
	}

	if len(p.rules) == 0 && p.condition == "" {
		return fmt.Errorf("either rules or condition must be specified")
	}

	return nil
}

// Process 处理数据
func (p *Filter) Process(ctx context.Context, input <-chan *record.Record, output chan<- *record.Record) error {
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

			p.metrics.IncrRecordsProcessed(1)

			// 检查是否匹配过滤条件
			match := p.matchRecord(rec)

			// 根据模式决定是否输出
			shouldOutput := (p.mode == "include" && match) || (p.mode == "exclude" && !match)

			if shouldOutput {
				select {
				case output <- rec:
					p.metrics.IncrSuccessCount(1)
				case <-ctx.Done():
					return ctx.Err()
				}
			} else {
				p.metrics.IncrFilteredCount(1)
			}
		}
	}
}

// matchRecord 检查记录是否匹配规则
func (p *Filter) matchRecord(rec *record.Record) bool {
	if len(p.rules) == 0 {
		// 简单条件匹配（这里简化处理）
		return true
	}

	// 所有规则都必须匹配
	for _, rule := range p.rules {
		if !p.matchRule(rec, rule) {
			return false
		}
	}

	return true
}

// matchRule 检查单个规则
func (p *Filter) matchRule(rec *record.Record, rule FilterRule) bool {
	value, ok := rec.Data[rule.Field]
	if !ok {
		return false
	}

	switch rule.Operator {
	case "eq", "==":
		return value == rule.Value

	case "ne", "!=":
		return value != rule.Value

	case "gt", ">":
		return compareNumbers(value, rule.Value) > 0

	case "lt", "<":
		return compareNumbers(value, rule.Value) < 0

	case "gte", ">=":
		return compareNumbers(value, rule.Value) >= 0

	case "lte", "<=":
		return compareNumbers(value, rule.Value) <= 0

	case "contains":
		str, ok1 := value.(string)
		substr, ok2 := rule.Value.(string)
		if ok1 && ok2 {
			return strings.Contains(str, substr)
		}

	case "startswith":
		str, ok1 := value.(string)
		prefix, ok2 := rule.Value.(string)
		if ok1 && ok2 {
			return strings.HasPrefix(str, prefix)
		}

	case "endswith":
		str, ok1 := value.(string)
		suffix, ok2 := rule.Value.(string)
		if ok1 && ok2 {
			return strings.HasSuffix(str, suffix)
		}
	}

	return false
}

// compareNumbers 比较数字
func compareNumbers(a, b interface{}) int {
	aNum := toFloat64(a)
	bNum := toFloat64(b)

	if aNum < bNum {
		return -1
	} else if aNum > bNum {
		return 1
	}
	return 0
}

// toFloat64 转换为 float64
func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case float64:
		return val
	case float32:
		return float64(val)
	default:
		return 0
	}
}

// GetMetrics 获取指标
func (p *Filter) GetMetrics() *plugin.Metrics {
	return p.metrics.GetMetrics()
}

// Close 关闭插件
func (p *Filter) Close() error {
	return nil
}

// 注册插件
func init() {
	info := &plugin.Info{
		Name:        "filter",
		Type:        plugin.TypeProcessor,
		Version:     "1.0.0",
		Author:      "FustGo Team",
		Description: "Filter records based on conditions",
		ConfigSchema: map[string]interface{}{
			"mode": "string (optional) - include or exclude. Default: include",
			"rules": []map[string]interface{}{
				{
					"field":    "string - Field name",
					"operator": "string - eq, ne, gt, lt, gte, lte, contains, startswith, endswith",
					"value":    "any - Value to compare",
				},
			},
		},
	}

	plugin.Register("filter", NewFilter, info)
}
