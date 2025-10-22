package pipeline

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// Status 管道状态
type Status string

const (
	StatusPending  Status = "pending"  // 等待执行
	StatusRunning  Status = "running"  // 正在运行
	StatusPaused   Status = "paused"   // 已暂停
	StatusStopped  Status = "stopped"  // 已停止
	StatusFailed   Status = "failed"   // 执行失败
	StatusComplete Status = "complete" // 执行完成
)

// String 返回状态字符串
func (s Status) String() string {
	return string(s)
}

// IsTerminal 判断是否为终止状态
func (s Status) IsTerminal() bool {
	return s == StatusStopped || s == StatusFailed || s == StatusComplete
}

// ErrorStrategy 错误处理策略
type ErrorStrategy string

const (
	// StrategyFailFast 快速失败，遇到错误立即停止
	StrategyFailFast ErrorStrategy = "fail-fast"

	// StrategyContinue 继续执行，跳过错误记录
	StrategyContinue ErrorStrategy = "continue"

	// StrategyRetry 重试，达到最大次数后失败
	StrategyRetry ErrorStrategy = "retry"
)

// Config 管道配置
type Config struct {
	// 并发控制
	Parallelism  int `json:"parallelism"`   // 并行度（processor 工作协程数）
	ChannelSize  int `json:"channel_size"`  // Channel 缓冲区大小
	MaxWorkers   int `json:"max_workers"`   // 最大工作协程数
	
	// 错误处理
	ErrorStrategy ErrorStrategy `json:"error_strategy"` // 错误策略
	MaxRetries    int           `json:"max_retries"`    // 最大重试次数
	RetryInterval time.Duration `json:"retry_interval"` // 重试间隔
	
	// 性能优化
	BatchSize     int           `json:"batch_size"`     // 批处理大小
	FlushInterval time.Duration `json:"flush_interval"` // 刷新间隔
	
	// 状态管理
	EnableCheckpoint bool   `json:"enable_checkpoint"` // 启用断点续传
	CheckpointPath   string `json:"checkpoint_path"`   // 断点保存路径
	CheckpointInterval time.Duration `json:"checkpoint_interval"` // 检查点间隔
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Parallelism:        4,
		ChannelSize:        1000,
		MaxWorkers:         100,
		ErrorStrategy:      StrategyContinue,
		MaxRetries:         3,
		RetryInterval:      time.Second,
		BatchSize:          100,
		FlushInterval:      time.Second,
		EnableCheckpoint:   false,
		CheckpointPath:     "./checkpoints",
		CheckpointInterval: 10 * time.Second,
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Parallelism < 1 {
		c.Parallelism = 1
	}
	if c.ChannelSize < 1 {
		c.ChannelSize = 100
	}
	if c.MaxWorkers < 1 {
		c.MaxWorkers = 100
	}
	if c.MaxRetries < 0 {
		c.MaxRetries = 0
	}
	if c.BatchSize < 1 {
		c.BatchSize = 1
	}
	return nil
}

// State 管道状态
type State struct {
	Status         Status      `json:"status"`          // 状态
	StartTime      time.Time   `json:"start_time"`      // 开始时间
	EndTime        time.Time   `json:"end_time"`        // 结束时间
	LastCheckpoint *Checkpoint `json:"last_checkpoint"` // 最后的检查点
	Error          string      `json:"error,omitempty"` // 错误信息
	
	mu sync.RWMutex `json:"-"` // 读写锁
}

// NewState 创建新状态
func NewState() *State {
	return &State{
		Status: StatusPending,
	}
}

// SetStatus 设置状态
func (s *State) SetStatus(status Status) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Status = status
}

// GetStatus 获取状态
func (s *State) GetStatus() Status {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Status
}

// SetError 设置错误
func (s *State) SetError(err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err != nil {
		s.Error = err.Error()
		s.Status = StatusFailed
	}
}

// GetError 获取错误
func (s *State) GetError() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.Error
}

// SetStartTime 设置开始时间
func (s *State) SetStartTime(t time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.StartTime = t
}

// SetEndTime 设置结束时间
func (s *State) SetEndTime(t time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.EndTime = t
}

// Duration 获取执行时长
func (s *State) Duration() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.StartTime.IsZero() {
		return 0
	}
	if s.EndTime.IsZero() {
		return time.Since(s.StartTime)
	}
	return s.EndTime.Sub(s.StartTime)
}

// Checkpoint 检查点
type Checkpoint struct {
	PipelineID   string                 `json:"pipeline_id"`   // 管道ID
	Timestamp    time.Time              `json:"timestamp"`     // 时间戳
	RecordsRead  int64                  `json:"records_read"`  // 已读取记录数
	LastRecordID string                 `json:"last_record_id"` // 最后处理的记录ID
	ReaderState  map[string]interface{} `json:"reader_state"`  // Reader 状态
	WriterState  map[string]interface{} `json:"writer_state"`  // Writer 状态
}

// Metrics 管道指标
type Metrics struct {
	// 记录数统计
	RecordsRead      int64 `json:"records_read"`      // 已读取记录数
	RecordsProcessed int64 `json:"records_processed"` // 已处理记录数
	RecordsWritten   int64 `json:"records_written"`   // 已写入记录数
	RecordsFiltered  int64 `json:"records_filtered"`  // 已过滤记录数
	RecordsError     int64 `json:"records_error"`     // 错误记录数
	
	// 字节数统计
	BytesRead    int64 `json:"bytes_read"`    // 读取字节数
	BytesWritten int64 `json:"bytes_written"` // 写入字节数
	
	// 时间统计
	StartTime time.Time     `json:"start_time"` // 开始时间
	EndTime   time.Time     `json:"end_time"`   // 结束时间
	Duration  time.Duration `json:"duration"`   // 执行时长
	
	// 性能指标
	TPS          float64 `json:"tps"`            // 每秒事务数
	AvgLatency   int64   `json:"avg_latency_ns"` // 平均延迟（纳秒）
	
	mu sync.RWMutex `json:"-"` // 读写锁
}

// NewMetrics 创建新指标
func NewMetrics() *Metrics {
	return &Metrics{
		StartTime: time.Now(),
	}
}

// IncrRecordsRead 增加读取记录数
func (m *Metrics) IncrRecordsRead(n int64) {
	atomic.AddInt64(&m.RecordsRead, n)
}

// IncrRecordsProcessed 增加处理记录数
func (m *Metrics) IncrRecordsProcessed(n int64) {
	atomic.AddInt64(&m.RecordsProcessed, n)
}

// IncrRecordsWritten 增加写入记录数
func (m *Metrics) IncrRecordsWritten(n int64) {
	atomic.AddInt64(&m.RecordsWritten, n)
}

// IncrRecordsFiltered 增加过滤记录数
func (m *Metrics) IncrRecordsFiltered(n int64) {
	atomic.AddInt64(&m.RecordsFiltered, n)
}

// IncrRecordsError 增加错误记录数
func (m *Metrics) IncrRecordsError(n int64) {
	atomic.AddInt64(&m.RecordsError, n)
}

// IncrBytesRead 增加读取字节数
func (m *Metrics) IncrBytesRead(n int64) {
	atomic.AddInt64(&m.BytesRead, n)
}

// IncrBytesWritten 增加写入字节数
func (m *Metrics) IncrBytesWritten(n int64) {
	atomic.AddInt64(&m.BytesWritten, n)
}

// CalculateTPS 计算 TPS
func (m *Metrics) CalculateTPS() float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	duration := time.Since(m.StartTime).Seconds()
	if duration == 0 {
		return 0
	}
	
	totalRecords := atomic.LoadInt64(&m.RecordsProcessed)
	return float64(totalRecords) / duration
}

// UpdateTPS 更新 TPS
func (m *Metrics) UpdateTPS() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.TPS = m.CalculateTPS()
}

// Snapshot 获取指标快照
func (m *Metrics) Snapshot() *Metrics {
	return &Metrics{
		RecordsRead:      atomic.LoadInt64(&m.RecordsRead),
		RecordsProcessed: atomic.LoadInt64(&m.RecordsProcessed),
		RecordsWritten:   atomic.LoadInt64(&m.RecordsWritten),
		RecordsFiltered:  atomic.LoadInt64(&m.RecordsFiltered),
		RecordsError:     atomic.LoadInt64(&m.RecordsError),
		BytesRead:        atomic.LoadInt64(&m.BytesRead),
		BytesWritten:     atomic.LoadInt64(&m.BytesWritten),
		StartTime:        m.StartTime,
		EndTime:          m.EndTime,
		Duration:         m.Duration,
		TPS:              m.CalculateTPS(),
	}
}

// String 返回指标字符串表示
func (m *Metrics) String() string {
	snapshot := m.Snapshot()
	return fmt.Sprintf(
		"Metrics{Read:%d, Processed:%d, Written:%d, Filtered:%d, Error:%d, TPS:%.2f}",
		snapshot.RecordsRead,
		snapshot.RecordsProcessed,
		snapshot.RecordsWritten,
		snapshot.RecordsFiltered,
		snapshot.RecordsError,
		snapshot.TPS,
	)
}
