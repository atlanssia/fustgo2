package pipeline

import (
	"context"
	"fmt"

	"github.com/fustgo/fustgo2/pkg/plugin"
	"github.com/fustgo/fustgo2/pkg/record"
)

// Pipeline 管道
type Pipeline struct {
	// 基本信息
	ID          string `json:"id"`          // 管道唯一标识
	Name        string `json:"name"`        // 管道名称
	Description string `json:"description"` // 描述
	
	// 插件组件
	Reader     plugin.Reader       `json:"-"` // 数据源
	Processors []plugin.Processor  `json:"-"` // 处理器链
	Writer     plugin.Writer       `json:"-"` // 目标写入器
	
	// 配置和状态
	Config  *Config  `json:"config"`  // 配置
	State   *State   `json:"state"`   // 状态
	Metrics *Metrics `json:"metrics"` // 指标
	
	// 内部控制
	ctx     context.Context    `json:"-"` // 上下文
	cancel  context.CancelFunc `json:"-"` // 取消函数
	
	// 数据通道
	channels []chan *record.Record `json:"-"` // 内部通道
}

// NewPipeline 创建新管道
func NewPipeline(id, name string) *Pipeline {
	return &Pipeline{
		ID:      id,
		Name:    name,
		Config:  DefaultConfig(),
		State:   NewState(),
		Metrics: NewMetrics(),
	}
}

// SetReader 设置 Reader
func (p *Pipeline) SetReader(reader plugin.Reader) {
	p.Reader = reader
}

// AddProcessor 添加 Processor
func (p *Pipeline) AddProcessor(processor plugin.Processor) {
	p.Processors = append(p.Processors, processor)
}

// SetWriter 设置 Writer
func (p *Pipeline) SetWriter(writer plugin.Writer) {
	p.Writer = writer
}

// SetConfig 设置配置
func (p *Pipeline) SetConfig(config *Config) {
	p.Config = config
}

// Validate 验证管道配置
func (p *Pipeline) Validate() error {
	if p.Reader == nil {
		return fmt.Errorf("reader is required")
	}
	
	if p.Writer == nil {
		return fmt.Errorf("writer is required")
	}
	
	if p.Config == nil {
		p.Config = DefaultConfig()
	}
	
	return p.Config.Validate()
}

// GetID 获取管道ID
func (p *Pipeline) GetID() string {
	return p.ID
}

// GetName 获取管道名称
func (p *Pipeline) GetName() string {
	return p.Name
}

// GetState 获取管道状态
func (p *Pipeline) GetState() *State {
	return p.State
}

// GetMetrics 获取管道指标
func (p *Pipeline) GetMetrics() *Metrics {
	return p.Metrics
}

// GetConfig 获取管道配置
func (p *Pipeline) GetConfig() *Config {
	return p.Config
}

// String 返回管道字符串表示
func (p *Pipeline) String() string {
	return fmt.Sprintf("Pipeline{ID:%s, Name:%s, Status:%s}", 
		p.ID, p.Name, p.State.GetStatus())
}
