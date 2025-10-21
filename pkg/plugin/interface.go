package plugin

import (
	"context"
	"fmt"

	"github.com/fustgo/fustgo2/pkg/record"
)

// Type 插件类型
type Type string

const (
	TypeReader    Type = "reader"
	TypeWriter    Type = "writer"
	TypeProcessor Type = "processor"
)

// Plugin 所有插件的基础接口
type Plugin interface {
	// Name 返回插件名称
	Name() string

	// Type 返回插件类型
	Type() Type

	// Init 初始化插件
	// config: 插件配置（从 YAML/JSON 解析而来）
	Init(config map[string]interface{}) error

	// Validate 验证配置是否有效
	Validate() error

	// Close 关闭插件并释放资源
	Close() error
}

// Reader 数据读取器接口
type Reader interface {
	Plugin

	// Read 从数据源读取数据
	// ctx: 上下文，用于取消操作
	// output: 输出通道，读取的记录发送到此通道
	// 返回错误时应该关闭输出通道
	Read(ctx context.Context, output chan<- *record.Record) error

	// GetProgress 获取读取进度（可选）
	GetProgress() *Progress
}

// Writer 数据写入器接口
type Writer interface {
	Plugin

	// Write 将数据写入目标
	// ctx: 上下文，用于取消操作
	// input: 输入通道，从此通道读取要写入的记录
	Write(ctx context.Context, input <-chan *record.Record) error

	// Flush 刷新缓冲区，确保所有数据已写入
	Flush() error

	// GetMetrics 获取写入指标
	GetMetrics() *Metrics
}

// Processor 数据处理器接口
type Processor interface {
	Plugin

	// Process 处理数据
	// ctx: 上下文，用于取消操作
	// input: 输入通道
	// output: 输出通道
	Process(ctx context.Context, input <-chan *record.Record, output chan<- *record.Record) error

	// GetMetrics 获取处理指标
	GetMetrics() *Metrics
}

// Progress 进度信息
type Progress struct {
	// TotalRecords 总记录数（如果已知）
	TotalRecords int64

	// ProcessedRecords 已处理记录数
	ProcessedRecords int64

	// BytesProcessed 已处理字节数
	BytesProcessed int64

	// Percentage 进度百分比 (0-100)
	Percentage float64

	// Message 进度消息
	Message string
}

// Metrics 指标信息
type Metrics struct {
	// RecordsProcessed 已处理记录数
	RecordsProcessed int64

	// BytesProcessed 已处理字节数
	BytesProcessed int64

	// ErrorCount 错误数
	ErrorCount int64

	// SuccessCount 成功数
	SuccessCount int64

	// FilteredCount 过滤数（Processor）
	FilteredCount int64
}

// Info 插件元信息
type Info struct {
	// Name 插件名称
	Name string

	// Type 插件类型
	Type Type

	// Version 版本
	Version string

	// Author 作者
	Author string

	// Description 描述
	Description string

	// ConfigSchema 配置 Schema（JSON Schema 格式）
	ConfigSchema map[string]interface{}
}

// Error 插件错误
type Error struct {
	// Plugin 插件名称
	Plugin string

	// Operation 操作名称
	Operation string

	// Err 底层错误
	Err error

	// Retryable 是否可重试
	Retryable bool
}

func (e *Error) Error() string {
	return fmt.Sprintf("plugin %s error in %s: %v", e.Plugin, e.Operation, e.Err)
}

func (e *Error) Unwrap() error {
	return e.Err
}

// NewError 创建插件错误
func NewError(plugin, operation string, err error, retryable bool) *Error {
	return &Error{
		Plugin:    plugin,
		Operation: operation,
		Err:       err,
		Retryable: retryable,
	}
}

// IsRetryable 判断错误是否可重试
func IsRetryable(err error) bool {
	if pluginErr, ok := err.(*Error); ok {
		return pluginErr.Retryable
	}
	return false
}
