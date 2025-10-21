package types

import (
	"time"
)

// PipelineStatus 管道状态
type PipelineStatus string

const (
	PipelineStatusDraft     PipelineStatus = "draft"     // 草稿
	PipelineStatusActive    PipelineStatus = "active"    // 激活
	PipelineStatusInactive  PipelineStatus = "inactive"  // 停用
	PipelineStatusArchived  PipelineStatus = "archived"  // 归档
)

// Pipeline 管道配置
type Pipeline struct {
	// ID 管道唯一标识
	ID string `json:"id" gorm:"primaryKey"`
	
	// Name 管道名称
	Name string `json:"name" gorm:"uniqueIndex;not null"`
	
	// Description 描述
	Description string `json:"description"`
	
	// Status 状态
	Status PipelineStatus `json:"status" gorm:"default:draft"`
	
	// Config 管道配置(JSON格式)
	Config string `json:"config" gorm:"type:text"`
	
	// Schedule 调度配置(JSON格式)
	Schedule string `json:"schedule" gorm:"type:text"`
	
	// Tags 标签
	Tags string `json:"tags" gorm:"type:text"`
	
	// Enabled 是否启用
	Enabled bool `json:"enabled" gorm:"default:true"`
	
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	
	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// CreatedBy 创建人
	CreatedBy string `json:"created_by"`
	
	// LastExecutionID 最后一次执行ID
	LastExecutionID string `json:"last_execution_id"`
	
	// LastExecutionTime 最后一次执行时间
	LastExecutionTime *time.Time `json:"last_execution_time"`
}

// TableName 指定表名
func (Pipeline) TableName() string {
	return "pipelines"
}

// PipelineConfig 管道配置结构
type PipelineConfig struct {
	// Source 数据源配置
	Source *PluginInstanceRef `json:"source"`
	
	// Processors 处理器配置列表
	Processors []*PluginInstanceRef `json:"processors,omitempty"`
	
	// Sink 数据目标配置
	Sink *PluginInstanceRef `json:"sink"`
	
	// Options 执行选项
	Options *PipelineOptions `json:"options,omitempty"`
}

// PluginInstanceRef 插件实例引用
type PluginInstanceRef struct {
	// InstanceID 实例ID
	InstanceID string `json:"instance_id"`
	
	// Config 覆盖配置(JSON格式)
	Config string `json:"config,omitempty" gorm:"type:text"`
}

// PipelineOptions 管道执行选项
type PipelineOptions struct {
	// Concurrency 并发配置
	Concurrency *ConcurrencyConfig `json:"concurrency,omitempty"`
	
	// Retry 重试策略
	Retry *RetryConfig `json:"retry,omitempty"`
	
	// ErrorHandling 错误处理
	ErrorHandling *ErrorHandlingConfig `json:"error_handling,omitempty"`
	
	// Validation 数据验证
	Validation *ValidationConfig `json:"validation,omitempty"`
	
	// Timeout 超时时间(秒)
	Timeout int `json:"timeout,omitempty"`
}