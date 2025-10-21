package types

import (
	"time"
)

// JobStatus 任务状态
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"   // 待执行
	JobStatusRunning   JobStatus = "running"   // 运行中
	JobStatusSucceeded JobStatus = "succeeded" // 成功
	JobStatusFailed    JobStatus = "failed"    // 失败
	JobStatusCancelled JobStatus = "cancelled" // 已取消
	JobStatusPaused    JobStatus = "paused"    // 已暂停
)

// ScheduleType 调度类型
type ScheduleType string

const (
	ScheduleTypeCron       ScheduleType = "cron"       // Cron 定时调度
	ScheduleTypeManual     ScheduleType = "manual"     // 手动触发
	ScheduleTypeEvent      ScheduleType = "event"      // 事件触发
	ScheduleTypeDependency ScheduleType = "dependency" // 依赖调度
)

// Job 数据同步任务
type Job struct {
	// ID 任务唯一标识
	ID string `json:"id" gorm:"primaryKey"`
	
	// Name 任务名称
	Name string `json:"name" gorm:"uniqueIndex;not null" validate:"required"`
	
	// Description 任务描述
	Description string `json:"description"`
	
	// Enabled 是否启用
	Enabled bool `json:"enabled" gorm:"default:true"`
	
	// Status 任务状态
	Status JobStatus `json:"status" gorm:"default:pending"`
	
	// PipelineConfig 管道配置(JSON 格式)
	PipelineConfig string `json:"pipeline_config" gorm:"type:text"`
	
	// ScheduleConfig 调度配置(JSON 格式)
	ScheduleConfig string `json:"schedule_config" gorm:"type:text"`
	
	// RetryConfig 重试配置(JSON 格式)
	RetryConfig string `json:"retry_config" gorm:"type:text"`
	
	// NotificationConfig 通知配置(JSON 格式)
	NotificationConfig string `json:"notification_config" gorm:"type:text"`
	
	// Tags 标签
	Tags string `json:"tags" gorm:"type:text"`
	
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
	
	// NextExecutionTime 下次执行时间
	NextExecutionTime *time.Time `json:"next_execution_time"`
}

// TableName 指定表名
func (Job) TableName() string {
	return "jobs"
}

// PipelineConfig 管道配置
type PipelineConfig struct {
	// Source 数据源配置
	Source PluginConfig `json:"source" validate:"required"`
	
	// Transform 转换配置列表
	Transform []PluginConfig `json:"transform,omitempty"`
	
	// Sink 目标配置
	Sink PluginConfig `json:"sink" validate:"required"`
	
	// Options 执行选项
	Options *ExecutionOptions `json:"options,omitempty"`
}

// PluginConfig 插件配置
type PluginConfig struct {
	// Plugin 插件名称
	Plugin string `json:"plugin" validate:"required"`
	
	// Connection 连接配置
	Connection map[string]interface{} `json:"connection,omitempty"`
	
	// Config 插件特定配置
	Config map[string]interface{} `json:"config,omitempty"`
}

// ExecutionOptions 执行选项
type ExecutionOptions struct {
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

// ConcurrencyConfig 并发配置
type ConcurrencyConfig struct {
	// ReaderWorkers Reader 并发数
	ReaderWorkers int `json:"reader_workers,omitempty" validate:"omitempty,min=1"`
	
	// TransformerWorkers Transformer 并发数
	TransformerWorkers int `json:"transformer_workers,omitempty" validate:"omitempty,min=1"`
	
	// WriterWorkers Writer 并发数
	WriterWorkers int `json:"writer_workers,omitempty" validate:"omitempty,min=1"`
}

// RetryConfig 重试配置
type RetryConfig struct {
	// MaxAttempts 最大重试次数
	MaxAttempts int `json:"max_attempts" validate:"min=0"`
	
	// InitialDelay 初始延迟(秒)
	InitialDelay int `json:"initial_delay" validate:"min=0"`
	
	// MaxDelay 最大延迟(秒)
	MaxDelay int `json:"max_delay" validate:"min=0"`
	
	// Multiplier 延迟倍数(指数退避)
	Multiplier float64 `json:"multiplier" validate:"min=1"`
}

// ErrorHandlingConfig 错误处理配置
type ErrorHandlingConfig struct {
	// Strategy 错误处理策略: abort, skip, retry
	Strategy string `json:"strategy" validate:"oneof=abort skip retry"`
	
	// MaxErrors 最大错误数(超过则中止)
	MaxErrors int `json:"max_errors,omitempty" validate:"min=0"`
	
	// LogErrors 是否记录错误数据
	LogErrors bool `json:"log_errors,omitempty"`
	
	// ErrorOutput 错误数据输出路径
	ErrorOutput string `json:"error_output,omitempty"`
}

// ValidationConfig 数据验证配置
type ValidationConfig struct {
	// Enabled 是否启用验证
	Enabled bool `json:"enabled"`
	
	// Mode 验证模式: count, sample, full
	Mode string `json:"mode" validate:"oneof=count sample full"`
	
	// SampleRate 抽样比例(sample 模式)
	SampleRate float64 `json:"sample_rate,omitempty" validate:"min=0,max=1"`
}

// ScheduleConfig 调度配置
type ScheduleConfig struct {
	// Type 调度类型
	Type ScheduleType `json:"type" validate:"required,oneof=cron manual event dependency"`
	
	// Cron Cron 表达式
	Cron string `json:"cron,omitempty"`
	
	// Timezone 时区
	Timezone string `json:"timezone,omitempty"`
	
	// Events 事件配置
	Events []EventConfig `json:"events,omitempty"`
	
	// Dependencies 依赖的任务 ID 列表
	Dependencies []string `json:"dependencies,omitempty"`
}

// EventConfig 事件配置
type EventConfig struct {
	// Type 事件类型: file, webhook, database
	Type string `json:"type" validate:"required"`
	
	// Config 事件配置
	Config map[string]interface{} `json:"config"`
}

// NotificationConfig 通知配置
type NotificationConfig struct {
	// Enabled 是否启用通知
	Enabled bool `json:"enabled"`
	
	// OnSuccess 成功时通知
	OnSuccess bool `json:"on_success"`
	
	// OnFailure 失败时通知
	OnFailure bool `json:"on_failure"`
	
	// OnRetry 重试时通知
	OnRetry bool `json:"on_retry"`
	
	// Channels 通知渠道
	Channels []NotificationChannel `json:"channels,omitempty"`
}

// NotificationChannel 通知渠道
type NotificationChannel struct {
	// Type 渠道类型: email, webhook, slack
	Type string `json:"type" validate:"required"`
	
	// Config 渠道配置
	Config map[string]interface{} `json:"config"`
}
