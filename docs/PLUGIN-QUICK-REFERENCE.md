# FustGo 插件系统 - 快速参考

## 🚀 快速开始

### 导入插件
```go
import (
    "github.com/fustgo/fustgo2/pkg/plugin"
    _ "github.com/fustgo/fustgo2/pkg/plugin/reader"
    _ "github.com/fustgo/fustgo2/pkg/plugin/writer"
    _ "github.com/fustgo/fustgo2/pkg/plugin/processor"
)
```

### 创建插件
```go
// Reader
reader, _ := plugin.CreateReader("file")
reader.Init(map[string]interface{}{"path": "data.csv"})

// Writer
writer, _ := plugin.CreateWriter("stdout")
writer.Init(map[string]interface{}{"format": "json"})

// Processor
mapper, _ := plugin.CreateProcessor("field_mapper")
mapper.Init(map[string]interface{}{
    "mappings": map[string]interface{}{"old": "new"},
})
```

### 构建管道
```go
ch1 := make(chan *record.Record)
ch2 := make(chan *record.Record)

go reader.Read(ctx, ch1)
go processor.Process(ctx, ch1, ch2)
writer.Write(ctx, ch2)
```

---

## 📦 内置插件

### FileReader
```yaml
path: "file.csv"        # 必填
format: "csv"           # 可选: csv|json|text
has_header: true        # 可选，默认 true
delimiter: ","          # 可选，默认 ","
```

### StdoutWriter
```yaml
format: "json"          # 可选: json|text|table
pretty: false           # 可选，默认 false
```

### FieldMapper
```yaml
mappings:
  old_name: new_name
  user_id: id
drop_unmapped: false    # 可选，默认 false
```

### Filter
```yaml
mode: "include"         # 可选: include|exclude
rules:
  - field: "age"
    operator: "gte"     # eq|ne|gt|lt|gte|lte|contains|startswith|endswith
    value: 18
```

---

## 🛠️ 开发插件

### Reader 模板
```go
type MyReader struct {
    *plugin.BasePlugin
    metrics  *plugin.BaseMetrics
    progress *plugin.BaseProgress
}

func (r *MyReader) Read(ctx context.Context, output chan<- *record.Record) error {
    defer close(output)
    // 读取数据并发送到 output
    return nil
}

func init() {
    plugin.Register("my_reader", NewMyReader, &plugin.Info{...})
}
```

### Writer 模板
```go
type MyWriter struct {
    *plugin.BasePlugin
    metrics *plugin.BaseMetrics
}

func (w *MyWriter) Write(ctx context.Context, input <-chan *record.Record) error {
    for rec := range input {
        // 写入数据
        w.metrics.IncrRecordsProcessed(1)
    }
    return nil
}
```

### Processor 模板
```go
type MyProcessor struct {
    *plugin.BasePlugin
    metrics *plugin.BaseMetrics
}

func (p *MyProcessor) Process(ctx context.Context, input <-chan *record.Record, output chan<- *record.Record) error {
    defer close(output)
    for rec := range input {
        // 处理数据
        output <- processed
    }
    return nil
}
```

---

## 🧪 测试

### 运行测试
```bash
# 所有测试
go test ./pkg/plugin/... -v

# 覆盖率
go test ./pkg/plugin/... -cover

# 基准测试
go test ./pkg/plugin -bench=. -benchmem
```

### 测试结果
```
✅ 21/21 测试通过
✅ 平均覆盖率 70.4%
✅ 性能: 1.5亿次/秒
```

---

## 📊 API 参考

### 插件注册
```go
plugin.Register(name string, factory Factory, info *Info)
plugin.CreateReader(name string) (Reader, error)
plugin.CreateWriter(name string) (Writer, error)
plugin.CreateProcessor(name string) (Processor, error)
```

### 插件发现
```go
plugin.ListReaders() []string
plugin.ListWriters() []string
plugin.ListProcessors() []string
plugin.GetInfo(name string) (*Info, error)
```

### 配置验证
```go
validator := plugin.NewConfigValidator(config)
validator.ValidateRequired(field string)
validator.ValidateType(field, expectedType string)
validator.ValidateRange(field string, min, max int)
validator.ValidateEnum(field string, values []string)
validator.ValidatePattern(field, pattern string)
```

### 指标和进度
```go
// Metrics
metrics.IncrRecordsProcessed(n int64)
metrics.IncrSuccessCount(n int64)
metrics.IncrErrorCount(n int64)
metrics.IncrFilteredCount(n int64)

// Progress
progress.SetProcessedRecords(n int64)
progress.SetTotalRecords(n int64)
progress.SetMessage(msg string)
progress.IncrProcessedRecords(n int64)
```

---

## 🎯 最佳实践

### 1. 使用 Context
```go
func (r *Reader) Read(ctx context.Context, output chan<- *record.Record) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()  // 优雅退出
        default:
            // 处理数据
        }
    }
}
```

### 2. 关闭 Channel
```go
func (p *Processor) Process(...) error {
    defer close(output)  // 确保关闭
    // 处理逻辑
}
```

### 3. 错误处理
```go
if err != nil {
    return plugin.NewError("my_plugin", "operation", err, retryable)
}
```

### 4. 线程安全
```go
// 使用 atomic
atomic.AddInt64(&m.recordsProcessed, n)

// 或使用 mutex
m.mu.Lock()
defer m.mu.Unlock()
```

---

## 🔍 调试技巧

### 1. 启用详细日志
```go
reader.SetLogLevel("debug")
```

### 2. 检查指标
```go
metrics := reader.GetMetrics()
fmt.Printf("Processed: %d, Errors: %d\n", 
    metrics.RecordsProcessed, metrics.ErrorCount)
```

### 3. 监控进度
```go
progress := reader.GetProgress()
fmt.Printf("Progress: %.2f%% (%d/%d)\n",
    progress.Percentage, 
    progress.ProcessedRecords,
    progress.TotalRecords)
```

---

## ⚡ 性能提示

1. **使用缓冲 Channel**: `make(chan *record.Record, 100)`
2. **批量处理**: 积累一批记录后批量写入
3. **并发处理**: 多个 Processor 并行处理
4. **避免复制**: 尽量传递指针
5. **复用对象**: 使用 sync.Pool

---

## 📚 更多资源

- 完整文档: [pkg/plugin/README.md](../pkg/plugin/README.md)
- 演示程序: [examples/plugin_demo.go](../examples/plugin_demo.go)
- 完成报告: [docs/phase2.1-completion-report.md](./phase2.1-completion-report.md)

---

**版本**: 1.0.0  
**更新**: 2025-10-21
