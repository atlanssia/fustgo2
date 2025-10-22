package types

import (
	"time"
)

// ExecutionStatus 执行状态
type ExecutionStatus string

const (
	ExecutionStatusPending   ExecutionStatus = "pending"   // 待执行
	ExecutionStatusRunning   ExecutionStatus = "running"   // 运行中
	ExecutionStatusSucceeded ExecutionStatus = "succeeded" // 成功
	ExecutionStatusFailed    ExecutionStatus = "failed"    // 失败
	ExecutionStatusCancelled ExecutionStatus = "cancelled" // 已取消
)

// Execution 任务执行记录
type Execution struct {
	// ID 执行唯一标识
	ID string `json:"id" gorm:"primaryKey"`
	
	// JobID 所属任务 ID
	JobID string `json:"job_id" gorm:"index;not null"`
	
	// Status 执行状态
	Status ExecutionStatus `json:"status" gorm:"default:pending;index"`
	
	// StartTime 开始时间
	StartTime time.Time `json:"start_time"`
	
	// EndTime 结束时间
	EndTime *time.Time `json:"end_time"`
	
	// Duration 执行时长(毫秒)
	Duration int64 `json:"duration"`
	
	// RecordsRead 读取记录数
	RecordsRead int64 `json:"records_read" gorm:"default:0"`
	
	// RecordsWritten 写入记录数
	RecordsWritten int64 `json:"records_written" gorm:"default:0"`
	
	// RecordsFiltered 过滤记录数
	RecordsFiltered int64 `json:"records_filtered" gorm:"default:0"`
	
	// RecordsError 错误记录数
	RecordsError int64 `json:"records_error" gorm:"default:0"`
	
	// BytesRead 读取字节数
	BytesRead int64 `json:"bytes_read" gorm:"default:0"`
	
	// BytesWritten 写入字节数
	BytesWritten int64 `json:"bytes_written" gorm:"default:0"`
	
	// ErrorMessage 错误信息
	ErrorMessage string `json:"error_message" gorm:"type:text"`
	
	// ErrorStack 错误堆栈
	ErrorStack string `json:"error_stack" gorm:"type:text"`
	
	// Context 执行上下文(JSON 格式)
	Context string `json:"context" gorm:"type:text"`
	
	// Metrics 执行指标(JSON 格式)
	Metrics string `json:"metrics" gorm:"type:text"`
	
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	
	// UpdatedAt 更新时间
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	
	// TriggeredBy 触发方式: cron, manual, event, dependency
	TriggeredBy string `json:"triggered_by"`
}

// TableName 指定表名
func (Execution) TableName() string {
	return "job_executions"
}

// ExecutionMetrics 执行指标
type ExecutionMetrics struct {
	// TotalRecords 总记录数
	TotalRecords int64 `json:"total_records"`
	
	// ProcessedRecords 已处理记录数
	ProcessedRecords int64 `json:"processed_records"`
	
	// SuccessRecords 成功记录数
	SuccessRecords int64 `json:"success_records"`
	
	// FailedRecords 失败记录数
	FailedRecords int64 `json:"failed_records"`
	
	// SkippedRecords 跳过记录数
	SkippedRecords int64 `json:"skipped_records"`
	
	// TotalBytes 总字节数
	TotalBytes int64 `json:"total_bytes"`
	
	// ProcessedBytes 已处理字节数
	ProcessedBytes int64 `json:"processed_bytes"`
	
	// ReadSpeed 读取速度(记录/秒)
	ReadSpeed float64 `json:"read_speed"`
	
	// WriteSpeed 写入速度(记录/秒)
	WriteSpeed float64 `json:"write_speed"`
	
	// Throughput 吞吐量(字节/秒)
	Throughput float64 `json:"throughput"`
	
	// Progress 进度百分比 (0-100)
	Progress float64 `json:"progress"`
	
	// StartTime 开始时间
	StartTime time.Time `json:"start_time"`
	
	// CurrentTime 当前时间
	CurrentTime time.Time `json:"current_time"`
	
	// EstimatedEndTime 预计结束时间
	EstimatedEndTime *time.Time `json:"estimated_end_time,omitempty"`
	
	// ReaderMetrics Reader 指标
	ReaderMetrics *ComponentMetrics `json:"reader_metrics,omitempty"`
	
	// TransformerMetrics Transformer 指标
	TransformerMetrics *ComponentMetrics `json:"transformer_metrics,omitempty"`
	
	// WriterMetrics Writer 指标
	WriterMetrics *ComponentMetrics `json:"writer_metrics,omitempty"`
}

// ComponentMetrics 组件指标
type ComponentMetrics struct {
	// Name 组件名称
	Name string `json:"name"`
	
	// RecordsProcessed 已处理记录数
	RecordsProcessed int64 `json:"records_processed"`
	
	// BytesProcessed 已处理字节数
	BytesProcessed int64 `json:"bytes_processed"`
	
	// ErrorCount 错误数
	ErrorCount int64 `json:"error_count"`
	
	// AvgProcessingTime 平均处理时间(毫秒)
	AvgProcessingTime float64 `json:"avg_processing_time"`
	
	// MaxProcessingTime 最大处理时间(毫秒)
	MaxProcessingTime float64 `json:"max_processing_time"`
	
	// MinProcessingTime 最小处理时间(毫秒)
	MinProcessingTime float64 `json:"min_processing_time"`
}

// Connection 数据源连接配置
type Connection struct {
	// ID 连接唯一标识
	ID string `json:"id" gorm:"primaryKey"`
	
	// Name 连接名称
	Name string `json:"name" gorm:"uniqueIndex;not null" validate:"required"`
	
	// Type 连接类型: mysql, postgresql, mongodb, etc.
	Type string `json:"type" gorm:"not null" validate:"required"`
	
	// Description 连接描述
	Description string `json:"description"`
	
	// Config 连接配置(JSON 格式,加密存储)
	Config string `json:"config" gorm:"type:text"`
	
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
	
	// LastTestedAt 最后测试时间
	LastTestedAt *time.Time `json:"last_tested_at"`
	
	// LastTestStatus 最后测试状态
	LastTestStatus string `json:"last_test_status"`
}

// TableName 指定表名
func (Connection) TableName() string {
	return "connections"
}

// ConnectionConfig 连接配置
type ConnectionConfig struct {
	// Host 主机地址
	Host string `json:"host" validate:"required"`
	
	// Port 端口
	Port int `json:"port" validate:"required,min=1,max=65535"`
	
	// Database 数据库名
	Database string `json:"database,omitempty"`
	
	// Username 用户名
	Username string `json:"username,omitempty"`
	
	// Password 密码(加密存储)
	Password string `json:"password,omitempty"`
	
	// Params 额外参数
	Params map[string]string `json:"params,omitempty"`
	
	// SSL SSL/TLS 配置
	SSL *SSLConfig `json:"ssl,omitempty"`
	
	// Pool 连接池配置
	Pool *PoolConfig `json:"pool,omitempty"`
}

// SSLConfig SSL/TLS 配置
type SSLConfig struct {
	// Enabled 是否启用
	Enabled bool `json:"enabled"`
	
	// CACert CA 证书路径
	CACert string `json:"ca_cert,omitempty"`
	
	// ClientCert 客户端证书路径
	ClientCert string `json:"client_cert,omitempty"`
	
	// ClientKey 客户端密钥路径
	ClientKey string `json:"client_key,omitempty"`
	
	// VerifyServerCert 是否验证服务器证书
	VerifyServerCert bool `json:"verify_server_cert"`
}

// PoolConfig 连接池配置
type PoolConfig struct {
	// MaxOpenConns 最大打开连接数
	MaxOpenConns int `json:"max_open_conns,omitempty" validate:"omitempty,min=1"`
	
	// MaxIdleConns 最大空闲连接数
	MaxIdleConns int `json:"max_idle_conns,omitempty" validate:"omitempty,min=0"`
	
	// ConnMaxLifetime 连接最大生命周期(秒)
	ConnMaxLifetime int `json:"conn_max_lifetime,omitempty" validate:"omitempty,min=0"`
	
	// ConnMaxIdleTime 连接最大空闲时间(秒)
	ConnMaxIdleTime int `json:"conn_max_idle_time,omitempty" validate:"omitempty,min=0"`
}
