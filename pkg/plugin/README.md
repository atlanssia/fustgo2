# FustGo 插件系统

FustGo 的插件系统是整个 ETL/ELT 平台的核心，提供了灵活、高效、易扩展的数据处理能力。

## 🏗️ 架构概览

```
┌─────────────────────────────────────────────────────────────┐
│                     插件注册表 (Registry)                     │
│  - 全局单例                                                   │
│  - 工厂模式创建插件                                           │
│  - 类型安全 (Reader/Writer/Processor)                        │
└─────────────────────────────────────────────────────────────┘
                              │
                              │ 注册/创建
                              ▼
┌─────────────┐      ┌─────────────┐      ┌─────────────┐
│   Reader    │─────▶│  Processor  │─────▶│   Writer    │
│             │      │             │      │             │
│ • File      │      │ • FieldMap  │      │ • Stdout    │
│ • Database  │      │ • Filter    │      │ • Database  │
│ • Kafka     │      │ • Aggregate │      │ • Kafka     │
└─────────────┘      └─────────────┘      └─────────────┘
      │                     │                     │
      └─────────────────────┴─────────────────────┘
                            │
                  channel 数据流 (Record)
```

## 📦 核心组件

### 1. 插件接口 (`interface.go`)

所有插件都实现基础 `Plugin` 接口：

```go
type Plugin interface {
    Name() string                              // 插件名称
    Type() Type                                // 插件类型
    Init(config map[string]interface{}) error // 初始化
    Validate() error                           // 验证配置
    Close() error                              // 关闭资源
}
```

#### Reader 接口
```go
type Reader interface {
    Plugin
    Read(ctx context.Context, output chan<- *record.Record) error
    GetProgress() *Progress
}
```

#### Writer 接口
```go
type Writer interface {
    Plugin
    Write(ctx context.Context, input <-chan *record.Record) error
    Flush() error
    GetMetrics() *Metrics
}
```

#### Processor 接口
```go
type Processor interface {
    Plugin
    Process(ctx context.Context, input <-chan *record.Record, output chan<- *record.Record) error
    GetMetrics() *Metrics
}
```

### 2. 基础实现 (`base.go`)

提供了可复用的基础功能：

- **BasePlugin**: 配置管理、名称/类型存储
- **BaseMetrics**: 原子操作的指标收集（无锁设计）
- **BaseProgress**: 线程安全的进度跟踪

### 3. 注册表 (`registry.go`)

全局插件注册和发现：

```go
// 注册插件（通常在 init() 中调用）
plugin.Register("file", NewFileReader, &plugin.Info{...})

// 创建插件实例
reader, err := plugin.CreateReader("file")
writer, err := plugin.CreateWriter("stdout")
processor, err := plugin.CreateProcessor("filter")

// 列出所有插件
readers := plugin.ListReaders()
writers := plugin.ListWriters()
processors := plugin.ListProcessors()
```

### 4. 配置验证 (`validator.go`)

强大的配置验证框架：

```go
validator := plugin.NewConfigValidator(config)

// 内置验证器
validator.ValidateRequired("field_name")
validator.ValidateType("port", "int")
validator.ValidateRange("timeout", 1, 300)
validator.ValidateEnum("format", []string{"csv", "json"})
validator.ValidatePattern("email", `^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)

// 自定义验证函数
validator.AddRule("custom_field", func(v interface{}) error {
    // 自定义验证逻辑
    return nil
})
```

## 🔌 内置插件

### Reader

#### FileReader (`pkg/plugin/reader/file.go`)

从文件读取数据，支持多种格式。

**配置参数**:
```yaml
path: "/path/to/file.csv"      # 必填：文件路径
format: "csv"                   # 可选：csv|json|text，默认自动检测
encoding: "utf-8"               # 可选：文件编码，默认 UTF-8
delimiter: ","                  # 可选：CSV 分隔符，默认逗号
has_header: true                # 可选：CSV 是否有表头，默认 true
```

**支持格式**:
- **CSV**: 自动解析表头，支持自定义分隔符
- **JSON/JSONL**: 每行一个 JSON 对象
- **Text**: 纯文本，每行作为一条记录

**使用示例**:
```go
reader, _ := plugin.CreateReader("file")
reader.Init(map[string]interface{}{
    "path": "data.csv",
    "format": "csv",
})

output := make(chan *record.Record)
go reader.Read(context.Background(), output)

for rec := range output {
    fmt.Printf("Record: %+v\n", rec.Data)
}
```

### Writer

#### StdoutWriter (`pkg/plugin/writer/stdout.go`)

输出到标准输出（控制台），支持多种格式。

**配置参数**:
```yaml
format: "json"    # 可选：json|text|table，默认 json
pretty: true      # 可选：JSON 美化，默认 false
```

**输出格式**:
- **JSON**: 紧凑或美化的 JSON
- **Text**: key=value 格式
- **Table**: 表格格式

**使用示例**:
```go
writer, _ := plugin.CreateWriter("stdout")
writer.Init(map[string]interface{}{
    "format": "json",
    "pretty": true,
})

input := make(chan *record.Record)
go writer.Write(context.Background(), input)

input <- record.NewRecord(map[string]interface{}{
    "id": 1,
    "name": "Alice",
})
close(input)
```

### Processor

#### FieldMapper (`pkg/plugin/processor/field_mapper.go`)

字段名映射处理器。

**配置参数**:
```yaml
mappings:           # 必填：字段映射规则
  old_name: new_name
  user_id: id
drop_unmapped: false  # 可选：是否删除未映射字段，默认 false
```

**使用示例**:
```go
mapper, _ := plugin.CreateProcessor("field_mapper")
mapper.Init(map[string]interface{}{
    "mappings": map[string]interface{}{
        "user_id": "id",
        "user_name": "name",
    },
})

// 输入: {user_id: 123, user_name: "John"}
// 输出: {id: 123, name: "John"}
```

#### Filter (`pkg/plugin/processor/filter.go`)

基于条件过滤记录。

**配置参数**:
```yaml
mode: "include"     # 可选：include|exclude，默认 include
rules:              # 必填：过滤规则列表
  - field: "age"
    operator: "gte"
    value: 18
  - field: "name"
    operator: "contains"
    value: "test"
```

**支持的操作符**:
- 比较: `eq`, `ne`, `gt`, `lt`, `gte`, `lte`
- 字符串: `contains`, `startswith`, `endswith`

**使用示例**:
```go
filter, _ := plugin.CreateProcessor("filter")
filter.Init(map[string]interface{}{
    "mode": "include",
    "rules": []interface{}{
        map[string]interface{}{
            "field": "age",
            "operator": "gte",
            "value": 18,
        },
    },
})

// 只保留 age >= 18 的记录
```

## 🚀 快速开始

### 1. 基本使用

```go
package main

import (
    "context"
    "github.com/fustgo/fustgo2/pkg/plugin"
    _ "github.com/fustgo/fustgo2/pkg/plugin/reader"  // 导入 readers
    _ "github.com/fustgo/fustgo2/pkg/plugin/writer"  // 导入 writers
)

func main() {
    // 创建 Reader
    reader, _ := plugin.CreateReader("file")
    reader.Init(map[string]interface{}{
        "path": "data.csv",
    })
    
    // 创建 Writer
    writer, _ := plugin.CreateWriter("stdout")
    writer.Init(map[string]interface{}{
        "format": "json",
    })
    
    // 连接管道
    channel := make(chan *record.Record, 10)
    ctx := context.Background()
    
    go reader.Read(ctx, channel)
    writer.Write(ctx, channel)
}
```

### 2. 构建数据管道

```go
// Reader -> Processor -> Writer 管道
readerOutput := make(chan *record.Record)
processorOutput := make(chan *record.Record)

go reader.Read(ctx, readerOutput)
go processor.Process(ctx, readerOutput, processorOutput)
writer.Write(ctx, processorOutput)
```

### 3. 运行演示

```bash
# 运行完整演示程序
go run examples/plugin_demo.go
```

## 🧪 测试

### 运行所有测试

```bash
# 运行测试
go test ./pkg/plugin/... -v

# 测试覆盖率
go test ./pkg/plugin/... -cover

# 基准测试
go test ./pkg/plugin -bench=. -benchmem
```

### 测试结果

```
pkg/plugin           36.7% coverage  ✅
pkg/plugin/processor 79.3% coverage  ✅
pkg/plugin/reader    81.5% coverage  ✅
pkg/plugin/writer    84.5% coverage  ✅

21/21 测试全部通过 ✅
```

### 性能指标

```
BenchmarkMetricsIncrement  153,089,654 次/秒  (7.8ns/op, 0 allocs)
BenchmarkProgressUpdate    153,011,031 次/秒  (7.8ns/op, 0 allocs)
```

## 🛠️ 开发自定义插件

### 1. 创建 Reader 插件

```go
package reader

import (
    "context"
    "github.com/fustgo/fustgo2/pkg/plugin"
    "github.com/fustgo/fustgo2/pkg/record"
)

type MyReader struct {
    *plugin.BasePlugin
    metrics  *plugin.BaseMetrics
    progress *plugin.BaseProgress
    // 自定义字段...
}

func NewMyReader() plugin.Plugin {
    return &MyReader{
        BasePlugin: plugin.NewBasePlugin("my_reader", plugin.TypeReader),
        metrics:    plugin.NewBaseMetrics(),
        progress:   plugin.NewBaseProgress(),
    }
}

func (r *MyReader) Init(config map[string]interface{}) error {
    if err := r.BasePlugin.Init(config); err != nil {
        return err
    }
    // 初始化自定义字段...
    return nil
}

func (r *MyReader) Validate() error {
    // 验证配置...
    return nil
}

func (r *MyReader) Read(ctx context.Context, output chan<- *record.Record) error {
    defer close(output)
    
    // 读取数据并发送到 output channel
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // 读取数据...
            rec := record.NewRecord(data)
            output <- rec
            r.metrics.IncrRecordsProcessed(1)
        }
    }
}

func (r *MyReader) GetProgress() *plugin.Progress {
    return r.progress.GetProgress()
}

func (r *MyReader) Close() error {
    // 清理资源...
    return nil
}

// 注册插件
func init() {
    plugin.Register("my_reader", NewMyReader, &plugin.Info{
        Name:        "my_reader",
        Type:        plugin.TypeReader,
        Version:     "1.0.0",
        Description: "My custom reader",
    })
}
```

### 2. 创建 Processor 插件

```go
type MyProcessor struct {
    *plugin.BasePlugin
    metrics *plugin.BaseMetrics
}

func (p *MyProcessor) Process(ctx context.Context, input <-chan *record.Record, output chan<- *record.Record) error {
    defer close(output)
    
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case rec, ok := <-input:
            if !ok {
                return nil
            }
            
            // 处理记录...
            processed := p.processRecord(rec)
            
            output <- processed
            p.metrics.IncrRecordsProcessed(1)
        }
    }
}
```

### 3. 注册和使用

```go
// 在 init() 中自动注册
func init() {
    plugin.Register("my_plugin", NewMyPlugin, &plugin.Info{...})
}

// 使用时导入包
import _ "your/package/path"

// 创建插件
plugin, _ := plugin.CreateReader("my_plugin")
```

## 📊 设计特点

### 1. **简洁的接口**
- 最少的方法数量
- 清晰的职责划分
- 易于理解和实现

### 2. **高性能**
- 无锁指标更新（atomic）
- Channel 管道模式（天然背压）
- 零内存分配（基准测试验证）

### 3. **线程安全**
- 所有公共方法并发安全
- Context 支持优雅取消
- Mutex 保护共享状态

### 4. **易于扩展**
- 插件自动发现（init 注册）
- 工厂模式创建
- 配置验证框架可复用

### 5. **完善的错误处理**
- 统一的错误类型
- 可重试错误分类
- 详细的错误上下文

## 🔮 未来扩展

基于当前架构，可以轻松添加：

### 更多 Reader
- `DatabaseReader`: MySQL/PostgreSQL/MongoDB
- `KafkaReader`: 从 Kafka 消费
- `HTTPReader`: HTTP API 数据源
- `S3Reader`: AWS S3/MinIO

### 更多 Writer
- `DatabaseWriter`: 批量写入数据库
- `KafkaWriter`: 发送到 Kafka
- `HTTPWriter`: HTTP API 目标
- `S3Writer`: 写入对象存储

### 更多 Processor
- `Aggregator`: 数据聚合
- `Enricher`: 数据增强/关联
- `Validator`: 数据验证
- `Deduplicator`: 去重
- `Transformer`: 复杂数据转换

### 高级特性
- 插件热重载
- 插件沙箱隔离
- 动态配置更新
- 插件性能分析

## 📚 更多资源

- [完成报告](../../docs/phase2.1-completion-report.md)
- [演示程序](../../examples/plugin_demo.go)
- [测试用例](./plugin_test.go)

## 📝 许可证

Copyright © 2025 FustGo Team
