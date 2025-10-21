# FustGo 系统架构设计

## 1. 总体架构

FustGo 采用分层架构设计，从上到下分为五层：

1. **展示层**: Web UI (React + TypeScript + Tailwind CSS)
2. **接口层**: API Gateway (认证、限流、路由)
3. **应用层**: 核心业务服务
4. **领域层**: 数据管道引擎和插件系统
5. **基础设施层**: 存储、缓存、消息队列等

## 2. 核心模块设计

### 2.1 数据管道引擎

```go
// Pipeline 接口定义
type Pipeline interface {
    Init(config PipelineConfig) error
    Validate() error
    Execute(ctx context.Context) error
    Pause() error
    Resume() error
    Stop() error
    GetMetrics() Metrics
}

// Record 数据记录
type Record struct {
    Data      map[string]interface{}
    Meta      RecordMeta
    Timestamp time.Time
}
```

### 2.2 插件系统

插件分为四类:
- **Connector**: 连接器 (管理数据源连接)
- **Reader**: 读取器 (从数据源读取数据)
- **Transformer**: 转换器 (数据处理和转换)
- **Writer**: 写入器 (将数据写入目标)

```go
// 插件注册机制
func init() {
    registry.Register("mysql", &MySQLPlugin{})
    registry.Register("postgresql", &PostgreSQLPlugin{})
}
```

### 2.3 任务调度器

支持四种调度模式:
- Cron 定时调度
- 事件触发调度
- 手动立即执行
- DAG 依赖调度

### 2.4 配置中心

配置层级:
- 系统配置 (服务器、数据库)
- 连接配置 (数据源连接信息)
- 任务配置 (管道定义)
- 插件配置 (插件参数)

## 3. 数据流设计

### 3.1 批处理流程
源数据库 → Reader Pool → Channel → Transformer Pool → Channel → Writer Pool → 目标数据库

### 3.2 流处理流程
消息队列 → Consumer Group → Stream Buffer → Transformer → Stream Buffer → Writer → 目标系统

### 3.3 CDC 增量同步
Binlog/WAL → Change Events → Parser → Transformer → Upsert Writer → 目标数据库

## 4. 存储设计

核心表:
- jobs: 任务定义
- job_executions: 执行历史
- connections: 数据源连接
- datasource_metadata: 数据源元信息
- job_dependencies: 任务依赖关系

## 5. 安全设计

- 认证: JWT Token / API Key / OAuth 2.0
- 授权: RBAC (Admin, Developer, Operator, Viewer)
- 加密: 敏感信息加密存储
- 审计: 完整的操作审计日志

## 6. 性能优化

- 协程池: 动态伸缩的 Worker Pool
- 批量处理: 批次大小和超时可配置
- 连接池: 数据库连接复用
- 数据压缩: gzip/snappy 压缩传输

## 7. 容错设计

- 错误分类: 可重试/不可重试/致命错误
- 重试策略: 指数退避重试
- 断点续传: Checkpoint 机制
- 一致性保证: At-Least-Once / Exactly-Once

## 8. 部署架构

### 单机模式
单一二进制文件包含所有组件 + SQLite

### 集群模式
API Server (多实例) + Scheduler (Leader选举) + Worker Pool + PostgreSQL HA

## 9. 技术栈

| 层次 | 组件 | 选型 | 理由 |
|------|------|------|------|
| 后端 | 语言 | Go 1.21+ | 高性能、并发友好 |
| 后端 | Web框架 | Gin | 轻量、性能优秀 |
| 后端 | ORM | GORM | 功能完善、易用 |
| 前端 | 框架 | React 18 | 生态成熟、组件丰富 |
| 前端 | 样式 | Tailwind CSS | 快速开发、体积小 |
| 前端 | 构建 | Vite | 极速热更新 |
| 存储 | 数据库 | PostgreSQL | 功能强大、JSON支持 |
| 存储 | 开发库 | SQLite | 轻量、零配置 |
| 缓存 | 缓存 | Redis | 高性能、数据结构丰富 |
| 消息 | 消息队列 | NATS | 轻量、云原生 |
| 监控 | 日志/指标 | OpenObserve | 一体化、资源占用低 |
