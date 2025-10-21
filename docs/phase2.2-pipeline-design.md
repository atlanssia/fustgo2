# Phase 2.2: Pipeline 引擎设计方案

## 📋 需求分析

### 核心需求

Pipeline 引擎是 FustGo 的核心执行引擎，负责将插件组合成数据处理管道并执行。

**主要职责**：
1. **管道构建**: 将 Reader、Processor、Writer 组合成完整的数据流管道
2. **执行调度**: 管理管道的启动、停止、暂停、恢复
3. **并发控制**: 控制并发度、协程池、背压处理
4. **监控指标**: 收集执行进度、性能指标、错误信息
5. **错误处理**: 错误捕获、重试策略、故障恢复
6. **状态管理**: 断点续传、状态持久化

### 使用场景

```
场景 1: 简单的批量同步
Reader (MySQL) → Writer (PostgreSQL)

场景 2: 带数据转换的同步
Reader (CSV) → Processor (FieldMapper) → Writer (Stdout)

场景 3: 多处理器链式处理
Reader (File) → Processor (Filter) → Processor (FieldMapper) → Writer (Database)

场景 4: 并行处理
Reader (Database) → [Processor1, Processor2, Processor3] → Writer (Multiple)
```

### 性能目标

- **吞吐量**: 单机 10万条/分钟（基础目标），50万条/分钟（优化目标）
- **延迟**: 端到端延迟 < 100ms（P99）
- **内存**: 稳定在 200MB 以内（处理 100万条数据）
- **并发**: 支持最多 1000 个并发 goroutine
- **可靠性**: 99.9% 成功率

## 🏗️ 架构设计

### 总体架构

```
┌─────────────────────────────────────────────────────────────┐
│                       Pipeline                              │
│                                                             │
│  ┌────────────┐    ┌────────────┐    ┌────────────┐       │
│  │   Reader   │───▶│ Processor  │───▶│   Writer   │       │
│  │            │    │   Chain    │    │            │       │
│  └────────────┘    └────────────┘    └────────────┘       │
│         │                 │                  │             │
│         └─────────────────┴──────────────────┘             │
│                          │                                 │
│                     ┌────▼────┐                            │
│                     │ Channel │  (data flow)               │
│                     └─────────┘                            │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ 管理
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                   Pipeline Executor                         │
│                                                             │
│  ├─ Lifecycle Manager   (启动、停止、暂停、恢复)            │
│  ├─ Concurrency Control (协程池、背压控制)                 │
│  ├─ Metrics Collector   (指标收集、进度跟踪)               │
│  ├─ Error Handler       (错误捕获、重试策略)               │
│  └─ State Manager       (状态持久化、断点续传)             │
└─────────────────────────────────────────────────────────────┘
```

### 核心组件

#### 1. Pipeline（管道）

**职责**: 管道定义和配置

```go
type Pipeline struct {
    ID          string              // 管道唯一标识
    Name        string              // 管道名称
    Description string              // 描述
    
    Reader      plugin.Reader       // 数据源
    Processors  []plugin.Processor  // 处理器链
    Writer      plugin.Writer       // 目标写入器
    
    Config      *PipelineConfig     // 配置
    State       *PipelineState      // 状态
    Metrics     *PipelineMetrics    // 指标
    
    ctx         context.Context     // 上下文
    cancel      context.CancelFunc  // 取消函数
    
    // 内部 channel
    channels    []chan *record.Record
}

type PipelineConfig struct {
    // 并发控制
    Parallelism      int           // 并行度
    ChannelSize      int           // Channel 缓冲区大小
    MaxWorkers       int           // 最大工作协程数
    
    // 错误处理
    ErrorStrategy    string        // fail-fast | continue | retry
    MaxRetries       int           // 最大重试次数
    RetryInterval    time.Duration // 重试间隔
    
    // 性能优化
    BatchSize        int           // 批处理大小
    FlushInterval    time.Duration // 刷新间隔
    
    // 状态管理
    EnableCheckpoint bool          // 启用断点续传
    CheckpointPath   string        // 断点保存路径
}

type PipelineState struct {
    Status           Status        // 状态: pending/running/paused/stopped/failed
    StartTime        time.Time     // 开始时间
    EndTime          time.Time     // 结束时间
    LastCheckpoint   *Checkpoint   // 最后的检查点
    Error            error         // 错误信息
    
    mu               sync.RWMutex  // 锁
}

type Status string

const (
    StatusPending  Status = "pending"
    StatusRunning  Status = "running"
    StatusPaused   Status = "paused"
    StatusStopped  Status = "stopped"
    StatusFailed   Status = "failed"
    StatusComplete Status = "complete"
)

type PipelineMetrics struct {
    RecordsRead      int64         // 已读取记录数
    RecordsProcessed int64         // 已处理记录数
    RecordsWritten   int64         // 已写入记录数
    RecordsFiltered  int64         // 已过滤记录数
    RecordsError     int64         // 错误记录数
    
    BytesRead        int64         // 读取字节数
    BytesWritten     int64         // 写入字节数
    
    Duration         time.Duration // 执行时长
    TPS              float64       // 每秒事务数
    
    mu               sync.RWMutex
}
```

#### 2. PipelineBuilder（构建器）

**职责**: 流式构建管道

```go
type PipelineBuilder struct {
    pipeline *Pipeline
    err      error
}

// 使用示例
pipeline, err := NewPipelineBuilder("my-pipeline").
    SetReader(reader).
    AddProcessor(fieldMapper).
    AddProcessor(filter).
    SetWriter(writer).
    WithConfig(&PipelineConfig{
        Parallelism: 4,
        ChannelSize: 1000,
    }).
    Build()
```

#### 3. PipelineExecutor（执行器）

**职责**: 管道执行和生命周期管理

```go
type PipelineExecutor struct {
    pipeline     *Pipeline
    workerPool   *WorkerPool       // 协程池
    errorHandler *ErrorHandler     // 错误处理器
    stateManager *StateManager     // 状态管理器
    
    wg           sync.WaitGroup    // 等待组
}

// 核心方法
func (e *PipelineExecutor) Start() error
func (e *PipelineExecutor) Stop() error
func (e *PipelineExecutor) Pause() error
func (e *PipelineExecutor) Resume() error
func (e *PipelineExecutor) Wait() error
func (e *PipelineExecutor) GetMetrics() *PipelineMetrics
```

#### 4. WorkerPool（协程池）

**职责**: 管理工作协程，控制并发度

```go
type WorkerPool struct {
    size     int                  // 池大小
    tasks    chan func()          // 任务队列
    workers  []*Worker            // 工作协程
    active   int32                // 活跃数（atomic）
    stopped  int32                // 停止标志
}

func (p *WorkerPool) Submit(task func()) error
func (p *WorkerPool) Stop()
func (p *WorkerPool) ActiveCount() int
```

#### 5. ErrorHandler（错误处理器）

**职责**: 错误捕获、重试、记录

```go
type ErrorHandler struct {
    strategy     string            // 错误策略
    maxRetries   int               // 最大重试次数
    retryInterval time.Duration    // 重试间隔
    
    errors       []error           // 错误列表
    mu           sync.Mutex
}

func (h *ErrorHandler) Handle(err error, record *record.Record) error
func (h *ErrorHandler) ShouldRetry(err error, attemptCount int) bool
func (h *ErrorHandler) GetErrors() []error
```

#### 6. StateManager（状态管理器）

**职责**: 状态持久化、断点续传

```go
type StateManager struct {
    enabled        bool
    checkpointPath string
    lastCheckpoint *Checkpoint
}

type Checkpoint struct {
    PipelineID     string
    Timestamp      time.Time
    RecordsRead    int64
    LastRecordID   string        // 最后处理的记录ID
    ReaderState    interface{}   // Reader 状态
    WriterState    interface{}   // Writer 状态
}

func (m *StateManager) SaveCheckpoint(checkpoint *Checkpoint) error
func (m *StateManager) LoadCheckpoint(pipelineID string) (*Checkpoint, error)
func (m *StateManager) ClearCheckpoint(pipelineID string) error
```

## 🔄 数据流设计

### 数据流模式

```
Reader 输出
    │
    ▼
┌─────────────────┐
│  Channel 1      │ (buffer: 1000)
│  (raw records)  │
└─────────────────┘
    │
    ▼
Processor 1
    │
    ▼
┌─────────────────┐
│  Channel 2      │
│  (processed)    │
└─────────────────┘
    │
    ▼
Processor 2
    │
    ▼
┌─────────────────┐
│  Channel 3      │
│  (filtered)     │
└─────────────────┘
    │
    ▼
Writer 输入
```

### Channel 管理策略

1. **缓冲区大小**: 根据内存和性能平衡
   - 默认: 1000
   - 高吞吐: 10000
   - 低内存: 100

2. **背压处理**: Channel 满时自动阻塞，提供反压
   
3. **超时控制**: Context 控制超时和取消

### 并发模型

```
Reader (1 goroutine)
    │
    ├──▶ Worker 1 ─┐
    ├──▶ Worker 2 ─┤
    ├──▶ Worker 3 ─┼──▶ Processor Chain
    └──▶ Worker 4 ─┘
            │
            ▼
        Writer (1 goroutine)
```

- **Reader**: 单协程读取，避免并发问题
- **Processor**: 多协程并行处理
- **Writer**: 单协程写入（可选批量）

## 📊 监控设计

### 指标类型

1. **吞吐量指标**
   - TPS (Transactions Per Second)
   - 读取速率 (records/sec)
   - 写入速率 (records/sec)

2. **延迟指标**
   - 端到端延迟 (P50, P90, P99)
   - 各阶段延迟

3. **资源指标**
   - CPU 使用率
   - 内存使用量
   - 协程数量
   - Channel 队列长度

4. **业务指标**
   - 成功记录数
   - 失败记录数
   - 过滤记录数
   - 错误率

### 监控接口

```go
type MetricsCollector interface {
    RecordRead(count int64)
    RecordProcessed(count int64)
    RecordWritten(count int64)
    RecordError(count int64)
    
    RecordLatency(stage string, duration time.Duration)
    
    Snapshot() *MetricsSnapshot
}

type MetricsSnapshot struct {
    Timestamp        time.Time
    RecordsRead      int64
    RecordsProcessed int64
    RecordsWritten   int64
    TPS              float64
    AvgLatency       time.Duration
}
```

## 🔧 错误处理策略

### 错误类型

1. **可恢复错误**
   - 网络超时
   - 数据库死锁
   - 临时性错误
   → 策略: 重试

2. **不可恢复错误**
   - 配置错误
   - 权限错误
   - 数据格式错误
   → 策略: 失败快速返回

3. **业务错误**
   - 数据校验失败
   - 过滤条件不匹配
   → 策略: 记录并继续

### 错误处理流程

```
错误发生
    │
    ▼
判断错误类型
    │
    ├─ 可恢复? ──Yes──▶ 重试（指数退避）
    │                       │
    │                   达到最大次数?
    │                       │
    │                   ──Yes──▶ 记录并失败
    │                       │
    │                   ──No───▶ 继续重试
    │
    └─ 不可恢复? ──Yes──▶ 立即失败
                │
            记录错误日志
                │
                ▼
        根据策略处理
            │
            ├─ fail-fast: 停止整个管道
            ├─ continue: 跳过该记录继续
            └─ retry: 重试该记录
```

### 重试策略

```go
type RetryStrategy struct {
    MaxRetries    int           // 最大重试次数
    InitialDelay  time.Duration // 初始延迟
    MaxDelay      time.Duration // 最大延迟
    Multiplier    float64       // 延迟倍数（指数退避）
    Jitter        bool          // 是否添加随机抖动
}

// 指数退避示例
// 第1次: 1s
// 第2次: 2s
// 第3次: 4s
// 第4次: 8s
// 第5次: 16s (cap at MaxDelay)
```

## 💾 状态管理设计

### Checkpoint 机制

**触发时机**:
1. 定时触发（每 N 秒）
2. 记录数触发（每处理 N 条记录）
3. 手动触发（API 调用）
4. 管道停止前

**保存内容**:
```json
{
  "pipeline_id": "pipeline-123",
  "timestamp": "2025-10-21T08:00:00Z",
  "metrics": {
    "records_read": 100000,
    "records_written": 99500,
    "records_error": 500
  },
  "reader_state": {
    "last_offset": 100000,
    "last_id": "user_123456"
  },
  "writer_state": {
    "batch_id": "batch_001"
  }
}
```

**恢复流程**:
1. 加载最后的 checkpoint
2. 恢复 Reader 状态（从断点位置继续读取）
3. 恢复 Writer 状态
4. 继续执行

### 状态转换

```
         ┌─────────┐
    ┌───▶│ PENDING │
    │    └────┬────┘
    │         │ Start()
    │         ▼
    │    ┌─────────┐
    │    │ RUNNING │◀─────────┐
    │    └────┬────┘          │
    │         │                │ Resume()
    │         ├─ Pause() ──▶┌──────┐
    │         │              │PAUSED│
    │         │              └──────┘
    │         │
    │         ├─ Stop() ───▶┌───────┐
    │         │             │STOPPED│
    │         │             └───────┘
    │         │
    │         ├─ Error ────▶┌──────┐
    │         │             │FAILED│
    │         │             └──────┘
    │         │
    │         ▼ Complete
    │    ┌─────────┐
    └────│COMPLETE │
         └─────────┘
```

## 🎯 API 设计

### Pipeline API

```go
// 创建管道
func NewPipeline(name string) *Pipeline

// 构建器模式
type PipelineBuilder interface {
    SetReader(reader plugin.Reader) *PipelineBuilder
    AddProcessor(processor plugin.Processor) *PipelineBuilder
    SetWriter(writer plugin.Writer) *PipelineBuilder
    WithConfig(config *PipelineConfig) *PipelineBuilder
    Build() (*Pipeline, error)
}

// 执行器
type PipelineExecutor interface {
    Start() error
    Stop() error
    Pause() error
    Resume() error
    Wait() error
    
    GetState() *PipelineState
    GetMetrics() *PipelineMetrics
}
```

### 使用示例

```go
// 1. 创建插件
reader, _ := plugin.CreateReader("file")
reader.Init(map[string]interface{}{"path": "data.csv"})

mapper, _ := plugin.CreateProcessor("field_mapper")
mapper.Init(map[string]interface{}{
    "mappings": map[string]interface{}{"old": "new"},
})

writer, _ := plugin.CreateWriter("stdout")
writer.Init(map[string]interface{}{"format": "json"})

// 2. 构建管道
pipeline, _ := NewPipelineBuilder("my-pipeline").
    SetReader(reader).
    AddProcessor(mapper).
    SetWriter(writer).
    WithConfig(&PipelineConfig{
        Parallelism: 4,
        ErrorStrategy: "continue",
    }).
    Build()

// 3. 创建执行器
executor := NewPipelineExecutor(pipeline)

// 4. 启动执行
if err := executor.Start(); err != nil {
    log.Fatal(err)
}

// 5. 监控进度
go func() {
    ticker := time.NewTicker(1 * time.Second)
    for range ticker.C {
        metrics := executor.GetMetrics()
        log.Printf("Progress: %d records, TPS: %.2f",
            metrics.RecordsProcessed, metrics.TPS)
    }
}()

// 6. 等待完成
executor.Wait()

// 7. 获取最终指标
finalMetrics := executor.GetMetrics()
log.Printf("Completed: %d records in %v",
    finalMetrics.RecordsProcessed, finalMetrics.Duration)
```

## 🔍 技术选型

### 并发控制

- **协程池**: 自定义实现（简单高效）
- **Channel**: Go 原生 channel（带缓冲）
- **Context**: 超时和取消控制

### 状态持久化

- **格式**: JSON（易读易调试）
- **存储**: 
  - 文件系统（默认）
  - 数据库（可选，后续支持）

### 监控指标

- **收集**: 原子操作（atomic）
- **暴露**: 
  - 内存指标（API 查询）
  - Prometheus（可选，后续支持）

## 📈 性能优化策略

### 1. 批量处理

Reader 和 Writer 支持批量操作，减少系统调用。

### 2. 零拷贝

Record 使用指针传递，避免数据复制。

### 3. 对象池

使用 `sync.Pool` 复用 Record 对象。

### 4. 并行处理

Processor 支持并行执行，充分利用多核。

### 5. 背压控制

Channel 缓冲区提供自然的背压机制。

## 🧪 测试策略

### 单元测试

- Pipeline 构建
- 状态转换
- 错误处理
- 指标收集

### 集成测试

- 完整数据流
- 多 Processor 链
- 错误恢复
- 断点续传

### 性能测试

- 吞吐量测试（10万、50万、100万条/分钟）
- 并发测试（不同并行度）
- 内存测试（长时间运行）
- 压力测试（极限负载）

### 测试用例

```go
TestPipelineBuilder
TestPipelineExecution
TestPipelinePauseResume
TestPipelineErrorHandling
TestPipelineRetry
TestPipelineCheckpoint
TestPipelineMetrics
TestPipelineConcurrency
BenchmarkPipelineThroughput
BenchmarkPipelineLatency
```

## 📝 总结

Phase 2.2 将实现一个：
- ✅ **灵活**: 支持任意插件组合
- ✅ **高效**: 目标 50万条/分钟吞吐量
- ✅ **可靠**: 完善的错误处理和重试
- ✅ **可观测**: 详细的监控指标
- ✅ **易用**: 简洁的 API 设计

的 Pipeline 引擎，为后续的 Job 调度和实时同步打下坚实基础。
