package pipeline

import (
	"fmt"

	"github.com/fustgo/fustgo2/pkg/plugin"
)

// Builder 管道构建器
type Builder struct {
	pipeline *Pipeline
	err      error
}

// NewBuilder 创建新的管道构建器
func NewBuilder(id, name string) *Builder {
	return &Builder{
		pipeline: NewPipeline(id, name),
	}
}

// SetReader 设置 Reader
func (b *Builder) SetReader(reader plugin.Reader) *Builder {
	if b.err != nil {
		return b
	}
	
	if reader == nil {
		b.err = fmt.Errorf("reader cannot be nil")
		return b
	}
	
	b.pipeline.SetReader(reader)
	return b
}

// AddProcessor 添加 Processor
func (b *Builder) AddProcessor(processor plugin.Processor) *Builder {
	if b.err != nil {
		return b
	}
	
	if processor == nil {
		b.err = fmt.Errorf("processor cannot be nil")
		return b
	}
	
	b.pipeline.AddProcessor(processor)
	return b
}

// SetWriter 设置 Writer
func (b *Builder) SetWriter(writer plugin.Writer) *Builder {
	if b.err != nil {
		return b
	}
	
	if writer == nil {
		b.err = fmt.Errorf("writer cannot be nil")
		return b
	}
	
	b.pipeline.SetWriter(writer)
	return b
}

// WithConfig 设置配置
func (b *Builder) WithConfig(config *Config) *Builder {
	if b.err != nil {
		return b
	}
	
	if config == nil {
		b.err = fmt.Errorf("config cannot be nil")
		return b
	}
	
	b.pipeline.SetConfig(config)
	return b
}

// WithDescription 设置描述
func (b *Builder) WithDescription(description string) *Builder {
	if b.err != nil {
		return b
	}
	
	b.pipeline.Description = description
	return b
}

// Build 构建管道
func (b *Builder) Build() (*Pipeline, error) {
	if b.err != nil {
		return nil, b.err
	}
	
	// 验证管道配置
	if err := b.pipeline.Validate(); err != nil {
		return nil, fmt.Errorf("pipeline validation failed: %w", err)
	}
	
	return b.pipeline, nil
}

// MustBuild 构建管道，如果失败则 panic
func (b *Builder) MustBuild() *Pipeline {
	pipeline, err := b.Build()
	if err != nil {
		panic(fmt.Sprintf("failed to build pipeline: %v", err))
	}
	return pipeline
}
