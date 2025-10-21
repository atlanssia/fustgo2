# FustGo 核心模块接口定义

## 1. 插件系统接口

### 1.1 基础插件接口

```go
package plugin

import (
    "context"
)

// PluginType 插件类型
type PluginType string

const (
    PluginTypeConnector   PluginType = "connector"
    PluginTypeReader      PluginType = "reader"
    PluginTypeTransformer PluginType = "transformer"
    PluginTypeWriter      PluginType = "writer"
)

// Plugin 所有插件的基础接口
type Plugin interface {
    // Name 返回插件名称
    Name() string
    
    // Type 返回插件类型
    Type() PluginType
    
    // Init 初始化插件
    Init(config map[string]interface{}) error
    
    // Validate 验证配置
    Validate() error
    
    // Close 关闭插件并释放资源
    Close() error
}

// PluginInfo 插件元信息
type PluginInfo struct {
    Name        string
    Version     string
    Author      string
    Description string
    Type        PluginType
    ConfigSchema map[string]interface{} // JSON Schema
}
```

### 1.2 连接器接口

```go
package connector

// Connector 数据源连接器接口
type Connector interface {
    plugin.Plugin
    
    // Connect 建立连接
    Connect(ctx context.Context) error
    
    // Disconnect 断开连接
    Disconnect() error
    
    // Ping 检查连接是否有效
    Ping(ctx context.Context) error
    
    // GetConnection 获取底层连接对象
    GetConnection() interface{}
    
    // GetMetadata 获取数据源元信息
    GetMetadata(ctx context.Context) (*Metadata, error)
}

// Metadata 数据源元信息
type Metadata struct {
    Type    string
    Version string
    Schemas []SchemaInfo
}

// SchemaInfo 数据库模式信息
type SchemaInfo struct {
    Name   string
    Tables []TableInfo
}

// TableInfo 表信息
type TableInfo struct {
    Schema  string
    Name    string
    Columns []ColumnInfo
    Indexes []IndexInfo
}

// ColumnInfo 列信息
type ColumnInfo struct {
    Name     string
    Type     string
    Nullable bool
    Default  interface{}
    Comment  string
}
```

### 1.3 读取器接口

```go
package reader

import (
    "context"
    "github.com/fustgo/fustgo2/pkg/record"
)

// Reader 数据读取器接口
type Reader interface {
    plugin.Plugin
    
    // Read 读取数据并发送到channel
    Read(ctx context.Context, output chan<- *record.Record) error
    
    // GetProgress 获取读取进度
    GetProgress() *Progress
    
    // Checkpoint 保存检查点
    Checkpoint() (*Checkpoint, error)
    
    // Resume 从检查点恢复
    Resume(checkpoint *Checkpoint) error
}

// Progress 读取进度
type Progress struct {
    TotalRecords     int64
    ProcessedRecords int64
    BytesRead        int64
    StartTime        time.Time
    LastUpdateTime   time.Time
}

// Checkpoint 检查点数据
type Checkpoint struct {
    Offset   int64
    Position interface{} // 位置信息(如binlog位置)
    Metadata map[string]interface{}
}

// BatchReader 批量读取器接口
type BatchReader interface {
    Reader
    
    // GetBatchSize 获取批次大小
    GetBatchSize() int
    
    // SetBatchSize 设置批次大小
    SetBatchSize(size int)
}

// StreamReader 流式读取器接口
type StreamReader interface {
    Reader
    
    // Subscribe 订阅数据流
    Subscribe(ctx context.Context) error
    
    // Unsubscribe 取消订阅
    Unsubscribe() error
}
```

### 1.4 转换器接口

```go
package transformer

import (
    "context"
    "github.com/fustgo/fustgo2/pkg/record"
)

// Transformer 数据转换器接口
type Transformer interface {
    plugin.Plugin
    
    // Transform 转换数据
    Transform(ctx context.Context, input <-chan *record.Record, output chan<- *record.Record) error
    
    // GetMetrics 获取转换指标
    GetMetrics() *Metrics
}

// Metrics 转换指标
type Metrics struct {
    InputRecords    int64
    OutputRecords   int64
    FilteredRecords int64
    ErrorRecords    int64
    TransformTime   time.Duration
}

// FieldMapper 字段映射转换器
type FieldMapper interface {
    Transformer
    
    // MapFields 映射字段
    MapFields(input map[string]interface{}) (map[string]interface{}, error)
}

// Filter 过滤器
type Filter interface {
    Transformer
    
    // Match 检查记录是否匹配过滤条件
    Match(record *record.Record) bool
}

// Aggregator 聚合器
type Aggregator interface {
    Transformer
    
    // Aggregate 聚合数据
    Aggregate(records []*record.Record) (*record.Record, error)
    
    // GetWindowSize 获取窗口大小
    GetWindowSize() time.Duration
}
```

### 1.5 写入器接口

```go
package writer

import (
    "context"
    "github.com/fustgo/fustgo2/pkg/record"
)

// Writer 数据写入器接口
type Writer interface {
    plugin.Plugin
    
    // Write 写入数据
    Write(ctx context.Context, input <-chan *record.Record) error
    
    // Flush 刷新缓冲区
    Flush(ctx context.Context) error
    
    // GetMetrics 获取写入指标
    GetMetrics() *Metrics
}

// Metrics 写入指标
type Metrics struct {
    RecordsWritten int64
    BytesWritten   int64
    ErrorRecords   int64
    WriteTime      time.Duration
}

// BatchWriter 批量写入器
type BatchWriter interface {
    Writer
    
    // GetBatchSize 获取批次大小
    GetBatchSize() int
    
    // SetBatchSize 设置批次大小
    SetBatchSize(size int)
    
    // WriteBatch 批量写入
    WriteBatch(ctx context.Context, records []*record.Record) error
}

// UpsertWriter 更新或插入写入器
type UpsertWriter interface {
    Writer
    
    // Upsert 更新或插入数据
    Upsert(ctx context.Context, record *record.Record) error
    
    // SetConflictKeys 设置冲突键
    SetConflictKeys(keys []string)
}
```

## 2. 数据管道接口

### 2.1 管道核心接口

```go
package pipeline

import (
    "context"
)

// Pipeline 数据管道接口
type Pipeline interface {
    // Init 初始化管道
    Init(config *Config) error
    
    // Validate 验证管道配置
    Validate() error
    
    // Execute 执行管道
    Execute(ctx context.Context) error
    
    // Pause 暂停管道
    Pause() error
    
    // Resume 恢复管道
    Resume() error
    
    // Stop 停止管道
    Stop() error
    
    // GetState 获取管道状态
    GetState() State
    
    // GetMetrics 获取管道指标
    GetMetrics() *Metrics
}

// State 管道状态
type State string

const (
    StateCreated     State = "created"
    StateValidating  State = "validating"
    StateValidated   State = "validated"
    StateInitializing State = "initializing"
    StateReady       State = "ready"
    StateRunning     State = "running"
    StatePaused      State = "paused"
    StateStopping    State = "stopping"
    StateStopped     State = "stopped"
    StateFailed      State = "failed"
)

// Config 管道配置
type Config struct {
    Name        string
    Description string
    Source      SourceConfig
    Transform   []TransformConfig
    Sink        SinkConfig
    Options     Options
}

// Metrics 管道指标
type Metrics struct {
    RecordsRead    int64
    RecordsWritten int64
    BytesRead      int64
    BytesWritten   int64
    ErrorCount     int64
    StartTime      time.Time
    EndTime        time.Time
    Duration       time.Duration
}
```

### 2.2 记录接口

```go
package record

import "time"

// Record 数据记录
type Record struct {
    Data      map[string]interface{}
    Meta      *RecordMeta
    Timestamp time.Time
}

// RecordMeta 记录元数据
type RecordMeta struct {
    Source    string
    Schema    string
    Table     string
    Operation Operation
    Offset    int64
    Extra     map[string]interface{}
}

// Operation 操作类型
type Operation string

const (
    OperationInsert Operation = "INSERT"
    OperationUpdate Operation = "UPDATE"
    OperationDelete Operation = "DELETE"
    OperationRead   Operation = "READ"
)

// NewRecord 创建新记录
func NewRecord(data map[string]interface{}) *Record {
    return &Record{
        Data:      data,
        Meta:      &RecordMeta{},
        Timestamp: time.Now(),
    }
}

// Clone 克隆记录
func (r *Record) Clone() *Record {
    data := make(map[string]interface{})
    for k, v := range r.Data {
        data[k] = v
    }
    
    return &Record{
        Data:      data,
        Meta:      r.Meta,
        Timestamp: r.Timestamp,
    }
}

// GetField 获取字段值
func (r *Record) GetField(name string) (interface{}, bool) {
    val, ok := r.Data[name]
    return val, ok
}

// SetField 设置字段值
func (r *Record) SetField(name string, value interface{}) {
    r.Data[name] = value
}
```

## 3. 调度器接口

```go
package scheduler

import (
    "context"
    "time"
)

// Scheduler 任务调度器接口
type Scheduler interface {
    // Schedule 调度任务
    Schedule(job *Job) error
    
    // Unschedule 取消调度
    Unschedule(jobID string) error
    
    // Trigger 立即触发任务
    Trigger(jobID string) error
    
    // Pause 暂停任务
    Pause(jobID string) error
    
    // Resume 恢复任务
    Resume(jobID string) error
    
    // GetJob 获取任务
    GetJob(jobID string) (*Job, error)
    
    // ListJobs 列出所有任务
    ListJobs() ([]*Job, error)
    
    // Start 启动调度器
    Start(ctx context.Context) error
    
    // Stop 停止调度器
    Stop() error
}

// Job 任务定义
type Job struct {
    ID          string
    Name        string
    Description string
    Pipeline    *pipeline.Config
    Schedule    *ScheduleConfig
    Retry       *RetryConfig
    Notify      *NotifyConfig
    Enabled     bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// ScheduleConfig 调度配置
type ScheduleConfig struct {
    Type         ScheduleType
    Cron         string
    Timezone     string
    Events       []EventConfig
    Dependencies []string
}

// ScheduleType 调度类型
type ScheduleType string

const (
    ScheduleTypeCron       ScheduleType = "cron"
    ScheduleTypeEvent      ScheduleType = "event"
    ScheduleTypeManual     ScheduleType = "manual"
    ScheduleTypeDependency ScheduleType = "dependency"
)

// RetryConfig 重试配置
type RetryConfig struct {
    MaxAttempts  int
    InitialDelay time.Duration
    MaxDelay     time.Duration
    Multiplier   float64
}

// Execution 任务执行
type Execution struct {
    ID          string
    JobID       string
    Status      ExecutionStatus
    StartTime   time.Time
    EndTime     time.Time
    Metrics     *pipeline.Metrics
    Error       string
}

// ExecutionStatus 执行状态
type ExecutionStatus string

const (
    ExecutionStatusPending   ExecutionStatus = "pending"
    ExecutionStatusRunning   ExecutionStatus = "running"
    ExecutionStatusSucceeded ExecutionStatus = "succeeded"
    ExecutionStatusFailed    ExecutionStatus = "failed"
    ExecutionStatusCancelled ExecutionStatus = "cancelled"
)
```

## 4. 配置管理接口

```go
package config

// Provider 配置提供者接口
type Provider interface {
    // Get 获取配置值
    Get(key string) (interface{}, error)
    
    // GetString 获取字符串配置
    GetString(key string) (string, error)
    
    // GetInt 获取整数配置
    GetInt(key string) (int, error)
    
    // GetBool 获取布尔配置
    GetBool(key string) (bool, error)
    
    // Set 设置配置值
    Set(key string, value interface{}) error
    
    // Watch 监听配置变化
    Watch(key string, callback func(value interface{})) error
    
    // Load 加载配置
    Load() error
    
    // Save 保存配置
    Save() error
}

// Manager 配置管理器
type Manager interface {
    // RegisterProvider 注册配置提供者
    RegisterProvider(name string, provider Provider) error
    
    // GetProvider 获取配置提供者
    GetProvider(name string) (Provider, error)
    
    // Merge 合并多个配置源
    Merge(providers ...Provider) (Provider, error)
}
```

## 5. 存储接口

```go
package storage

import "context"

// MetadataStore 元数据存储接口
type MetadataStore interface {
    // Jobs 任务相关操作
    CreateJob(ctx context.Context, job *Job) error
    GetJob(ctx context.Context, id string) (*Job, error)
    UpdateJob(ctx context.Context, job *Job) error
    DeleteJob(ctx context.Context, id string) error
    ListJobs(ctx context.Context, filter *JobFilter) ([]*Job, error)
    
    // Executions 执行历史相关操作
    CreateExecution(ctx context.Context, exec *Execution) error
    GetExecution(ctx context.Context, id string) (*Execution, error)
    UpdateExecution(ctx context.Context, exec *Execution) error
    ListExecutions(ctx context.Context, filter *ExecutionFilter) ([]*Execution, error)
    
    // Connections 连接相关操作
    CreateConnection(ctx context.Context, conn *Connection) error
    GetConnection(ctx context.Context, id string) (*Connection, error)
    UpdateConnection(ctx context.Context, conn *Connection) error
    DeleteConnection(ctx context.Context, id string) error
    ListConnections(ctx context.Context) ([]*Connection, error)
}

// CheckpointStore 检查点存储接口
type CheckpointStore interface {
    // Save 保存检查点
    Save(ctx context.Context, jobID, executionID string, checkpoint interface{}) error
    
    // Load 加载检查点
    Load(ctx context.Context, jobID, executionID string) (interface{}, error)
    
    // Delete 删除检查点
    Delete(ctx context.Context, jobID, executionID string) error
}

// CacheStore 缓存存储接口
type CacheStore interface {
    // Get 获取缓存
    Get(ctx context.Context, key string) (interface{}, error)
    
    // Set 设置缓存
    Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    
    // Delete 删除缓存
    Delete(ctx context.Context, key string) error
    
    // Exists 检查键是否存在
    Exists(ctx context.Context, key string) (bool, error)
}
```

## 6. 监控接口

```go
package monitor

import "context"

// MetricsCollector 指标收集器接口
type MetricsCollector interface {
    // RecordCounter 记录计数器
    RecordCounter(name string, value int64, tags map[string]string)
    
    // RecordGauge 记录仪表盘
    RecordGauge(name string, value float64, tags map[string]string)
    
    // RecordHistogram 记录直方图
    RecordHistogram(name string, value float64, tags map[string]string)
    
    // RecordTiming 记录耗时
    RecordTiming(name string, duration time.Duration, tags map[string]string)
}

// Logger 日志接口
type Logger interface {
    Debug(msg string, fields map[string]interface{})
    Info(msg string, fields map[string]interface{})
    Warn(msg string, fields map[string]interface{})
    Error(msg string, fields map[string]interface{})
    Fatal(msg string, fields map[string]interface{})
}

// Tracer 链路追踪接口
type Tracer interface {
    // StartSpan 开始span
    StartSpan(ctx context.Context, operationName string) (Span, context.Context)
}

// Span span接口
type Span interface {
    // SetTag 设置标签
    SetTag(key string, value interface{})
    
    // SetError 设置错误
    SetError(err error)
    
    // Finish 结束span
    Finish()
}
```
