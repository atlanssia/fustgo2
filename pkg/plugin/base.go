package plugin

import (
	"sync"
	"sync/atomic"
)

// BasePlugin 插件基础实现
// 提供通用功能，具体插件可以嵌入此结构
type BasePlugin struct {
	name   string
	typ    Type
	config map[string]interface{}
	mu     sync.RWMutex
}

// NewBasePlugin 创建基础插件
func NewBasePlugin(name string, typ Type) *BasePlugin {
	return &BasePlugin{
		name: name,
		typ:  typ,
	}
}

// Name 返回插件名称
func (b *BasePlugin) Name() string {
	return b.name
}

// Type 返回插件类型
func (b *BasePlugin) Type() Type {
	return b.typ
}

// Init 初始化插件（保存配置）
func (b *BasePlugin) Init(config map[string]interface{}) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.config = config
	return nil
}

// GetConfig 获取配置
func (b *BasePlugin) GetConfig() map[string]interface{} {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.config
}

// GetConfigString 获取字符串配置
func (b *BasePlugin) GetConfigString(key string) (string, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	
	if val, ok := b.config[key]; ok {
		if str, ok := val.(string); ok {
			return str, true
		}
	}
	return "", false
}

// GetConfigInt 获取整数配置
func (b *BasePlugin) GetConfigInt(key string) (int, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	
	if val, ok := b.config[key]; ok {
		switch v := val.(type) {
		case int:
			return v, true
		case int64:
			return int(v), true
		case float64:
			return int(v), true
		}
	}
	return 0, false
}

// GetConfigBool 获取布尔配置
func (b *BasePlugin) GetConfigBool(key string) (bool, bool) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	
	if val, ok := b.config[key]; ok {
		if b, ok := val.(bool); ok {
			return b, true
		}
	}
	return false, false
}

// Validate 默认验证实现（子类应该重写）
func (b *BasePlugin) Validate() error {
	return nil
}

// Close 默认关闭实现（子类可以重写）
func (b *BasePlugin) Close() error {
	return nil
}

// BaseMetrics 基础指标实现
type BaseMetrics struct {
	recordsProcessed atomic.Int64
	bytesProcessed   atomic.Int64
	errorCount       atomic.Int64
	successCount     atomic.Int64
	filteredCount    atomic.Int64
}

// NewBaseMetrics 创建基础指标
func NewBaseMetrics() *BaseMetrics {
	return &BaseMetrics{}
}

// IncrRecordsProcessed 增加已处理记录数
func (m *BaseMetrics) IncrRecordsProcessed(delta int64) {
	m.recordsProcessed.Add(delta)
}

// IncrBytesProcessed 增加已处理字节数
func (m *BaseMetrics) IncrBytesProcessed(delta int64) {
	m.bytesProcessed.Add(delta)
}

// IncrErrorCount 增加错误数
func (m *BaseMetrics) IncrErrorCount(delta int64) {
	m.errorCount.Add(delta)
}

// IncrSuccessCount 增加成功数
func (m *BaseMetrics) IncrSuccessCount(delta int64) {
	m.successCount.Add(delta)
}

// IncrFilteredCount 增加过滤数
func (m *BaseMetrics) IncrFilteredCount(delta int64) {
	m.filteredCount.Add(delta)
}

// GetMetrics 获取指标快照
func (m *BaseMetrics) GetMetrics() *Metrics {
	return &Metrics{
		RecordsProcessed: m.recordsProcessed.Load(),
		BytesProcessed:   m.bytesProcessed.Load(),
		ErrorCount:       m.errorCount.Load(),
		SuccessCount:     m.successCount.Load(),
		FilteredCount:    m.filteredCount.Load(),
	}
}

// Reset 重置指标
func (m *BaseMetrics) Reset() {
	m.recordsProcessed.Store(0)
	m.bytesProcessed.Store(0)
	m.errorCount.Store(0)
	m.successCount.Store(0)
	m.filteredCount.Store(0)
}

// BaseProgress 基础进度实现
type BaseProgress struct {
	totalRecords     atomic.Int64
	processedRecords atomic.Int64
	bytesProcessed   atomic.Int64
	message          atomic.Value // string
}

// NewBaseProgress 创建基础进度
func NewBaseProgress() *BaseProgress {
	p := &BaseProgress{}
	p.message.Store("")
	return p
}

// SetTotalRecords 设置总记录数
func (p *BaseProgress) SetTotalRecords(total int64) {
	p.totalRecords.Store(total)
}

// IncrProcessedRecords 增加已处理记录数
func (p *BaseProgress) IncrProcessedRecords(delta int64) {
	p.processedRecords.Add(delta)
}

// IncrBytesProcessed 增加已处理字节数
func (p *BaseProgress) IncrBytesProcessed(delta int64) {
	p.bytesProcessed.Add(delta)
}

// SetMessage 设置进度消息
func (p *BaseProgress) SetMessage(msg string) {
	p.message.Store(msg)
}

// GetProgress 获取进度快照
func (p *BaseProgress) GetProgress() *Progress {
	total := p.totalRecords.Load()
	processed := p.processedRecords.Load()
	
	var percentage float64
	if total > 0 {
		percentage = float64(processed) / float64(total) * 100
	}
	
	return &Progress{
		TotalRecords:     total,
		ProcessedRecords: processed,
		BytesProcessed:   p.bytesProcessed.Load(),
		Percentage:       percentage,
		Message:          p.message.Load().(string),
	}
}
