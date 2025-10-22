# Phase 1.1 快速参考指南

## 📦 核心类型定义

### Record (数据记录)

```go
import "github.com/fustgo/fustgo2/pkg/record"

// 创建记录
rec := record.NewRecord(map[string]interface{}{
    "id": 1,
    "name": "张三",
})

// 设置元数据
rec.Meta.Source = "mysql"
rec.Meta.Table = "users"
rec.Meta.Operation = record.OperationInsert

// 获取字段
name, ok := rec.GetString("name")
age, ok := rec.GetInt64("age")

// JSON 序列化
jsonStr, err := rec.ToJSON()
rec2, err := record.FromJSON(jsonStr)

// 克隆
cloned := rec.Clone()
```

### Job (任务定义)

```go
import "github.com/fustgo/fustgo2/internal/core/types"

job := &types.Job{
    ID:          "job-12345",
    Name:        "MySQL 到 PostgreSQL 同步",
    Description: "用户数据同步",
    Enabled:     true,
    Status:      types.JobStatusPending,
    // ... 配置 JSON 字符串
}

// 任务状态
const (
    JobStatusPending   // 待执行
    JobStatusRunning   // 运行中
    JobStatusSucceeded // 成功
    JobStatusFailed    // 失败
    JobStatusCancelled // 已取消
    JobStatusPaused    // 已暂停
)

// 调度类型
const (
    ScheduleTypeCron       // Cron 定时
    ScheduleTypeManual     // 手动触发
    ScheduleTypeEvent      // 事件触发
    ScheduleTypeDependency // 依赖调度
)
```

### Execution (执行记录)

```go
import "github.com/fustgo/fustgo2/internal/core/types"

execution := &types.Execution{
    ID:             "exec-12345",
    JobID:          "job-12345",
    Status:         types.ExecutionStatusRunning,
    StartTime:      time.Now(),
    RecordsRead:    1000,
    RecordsWritten: 950,
    RecordsError:   50,
}

// 执行状态
const (
    ExecutionStatusPending   // 待执行
    ExecutionStatusRunning   // 运行中
    ExecutionStatusSucceeded // 成功
    ExecutionStatusFailed    // 失败
    ExecutionStatusCancelled // 已取消
)
```

### Config (系统配置)

```go
import "github.com/fustgo/fustgo2/internal/config"

// 加载配置
cfg, err := config.Load("./configs/system.yaml")
if err != nil {
    log.Fatal(err)
}

// 访问配置
port := cfg.Server.Port
dbType := cfg.Database.Type
logLevel := cfg.Logging.Level

// 获取时间间隔
timeout := cfg.GetDuration(cfg.Executor.DefaultTimeout)
```

## 🔧 常用命令

### 开发命令

```bash
# 下载依赖
go mod tidy

# 运行服务器
go run ./cmd/server

# 使用自定义配置
go run ./cmd/server --config ./configs/system.yaml

# 查看版本
go run ./cmd/server --version

# 运行测试
go test ./...

# 查看覆盖率
go test ./... -cover

# 运行特定测试
go test ./pkg/record/ -v

# 基准测试
go test ./pkg/record/ -bench=.

# 格式化代码
make fmt

# 代码检查
make lint

# 构建
make build

# 清理
make clean
```

### Docker 命令

```bash
# 构建镜像
docker build -t fustgo:latest -f deployments/docker/Dockerfile .

# 使用 docker-compose 启动
cd deployments/docker
docker-compose up -d

# 查看日志
docker-compose logs -f fustgo

# 停止服务
docker-compose down
```

## 📊 API 端点

### Health Check

```bash
GET /health

Response:
{
  "status": "healthy",
  "version": "dev"
}
```

### Jobs API

```bash
# 列出任务
GET /api/v1/jobs

# 创建任务
POST /api/v1/jobs
Content-Type: application/json
{
  "name": "数据同步任务",
  "description": "MySQL 到 PostgreSQL",
  "enabled": true
}

# 获取任务详情
GET /api/v1/jobs/:id

# 更新任务
PUT /api/v1/jobs/:id

# 删除任务
DELETE /api/v1/jobs/:id
```

### Executions API

```bash
# 列出执行记录
GET /api/v1/executions

# 获取执行详情
GET /api/v1/executions/:id
```

### Connections API

```bash
# 列出连接
GET /api/v1/connections

# 创建连接
POST /api/v1/connections

# 测试连接
POST /api/v1/connections/:id/test

# 获取连接详情
GET /api/v1/connections/:id

# 更新连接
PUT /api/v1/connections/:id

# 删除连接
DELETE /api/v1/connections/:id
```

## 🗄️ 数据库表结构

### jobs 表

```sql
CREATE TABLE jobs (
    id VARCHAR PRIMARY KEY,
    name VARCHAR UNIQUE NOT NULL,
    description TEXT,
    enabled BOOLEAN DEFAULT true,
    status VARCHAR DEFAULT 'pending',
    pipeline_config TEXT,
    schedule_config TEXT,
    retry_config TEXT,
    notification_config TEXT,
    tags TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    created_by VARCHAR,
    last_execution_id VARCHAR,
    last_execution_time TIMESTAMP,
    next_execution_time TIMESTAMP
);
```

### job_executions 表

```sql
CREATE TABLE job_executions (
    id VARCHAR PRIMARY KEY,
    job_id VARCHAR NOT NULL,
    status VARCHAR DEFAULT 'pending',
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    duration BIGINT,
    records_read BIGINT DEFAULT 0,
    records_written BIGINT DEFAULT 0,
    records_filtered BIGINT DEFAULT 0,
    records_error BIGINT DEFAULT 0,
    bytes_read BIGINT DEFAULT 0,
    bytes_written BIGINT DEFAULT 0,
    error_message TEXT,
    error_stack TEXT,
    context TEXT,
    metrics TEXT,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    triggered_by VARCHAR
);
```

### connections 表

```sql
CREATE TABLE connections (
    id VARCHAR PRIMARY KEY,
    name VARCHAR UNIQUE NOT NULL,
    type VARCHAR NOT NULL,
    description TEXT,
    config TEXT,
    tags TEXT,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    created_by VARCHAR,
    last_tested_at TIMESTAMP,
    last_test_status VARCHAR
);
```

## 📝 配置文件示例

### system.yaml (最小配置)

```yaml
server:
  port: 8080
  mode: release

database:
  type: sqlite
  sqlite:
    path: ./data/fustgo.db

logging:
  level: info
  format: json

executor:
  max_concurrent_jobs: 50
  
scheduler:
  enabled: true
  timezone: Asia/Shanghai
```

### 任务配置 JSON

```json
{
  "source": {
    "plugin": "mysql",
    "connection": {
      "host": "localhost",
      "port": 3306,
      "database": "source_db"
    },
    "config": {
      "table": "users"
    }
  },
  "transform": [
    {
      "plugin": "field_mapper",
      "config": {
        "mappings": {
          "user_id": "id",
          "user_name": "name"
        }
      }
    }
  ],
  "sink": {
    "plugin": "postgresql",
    "connection": {
      "host": "localhost",
      "port": 5432,
      "database": "target_db"
    },
    "config": {
      "table": "dim_users",
      "mode": "upsert"
    }
  }
}
```

## 🧪 测试示例

### 单元测试

```go
package mypackage

import (
    "testing"
    "github.com/fustgo/fustgo2/pkg/record"
)

func TestMyFunction(t *testing.T) {
    // 准备测试数据
    rec := record.NewRecord(map[string]interface{}{
        "id": 1,
        "name": "test",
    })
    
    // 执行测试
    result := MyFunction(rec)
    
    // 验证结果
    if result == nil {
        t.Error("Expected non-nil result")
    }
}
```

### 基准测试

```go
func BenchmarkMyFunction(b *testing.B) {
    rec := record.NewRecord(map[string]interface{}{
        "id": 1,
    })
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        MyFunction(rec)
    }
}
```

## 🔍 调试技巧

### 1. 查看日志

```bash
# 查看服务器日志
tail -f ./logs/fustgo.log

# JSON 日志格式化
tail -f ./logs/fustgo.log | jq .
```

### 2. 使用 Delve 调试

```bash
# 安装 Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 调试服务器
dlv debug ./cmd/server

# 设置断点
(dlv) break main.main
(dlv) continue
```

### 3. pprof 性能分析

```go
import _ "net/http/pprof"

// 访问 http://localhost:8080/debug/pprof/
```

## 📚 相关资源

### 文档
- [架构设计](../architecture.md)
- [技术选型](../tech-selection.md)
- [开发路线图](../roadmap.md)

### 官方文档
- [Gin 文档](https://gin-gonic.com/docs/)
- [GORM 文档](https://gorm.io/docs/)
- [Viper 文档](https://github.com/spf13/viper)
- [Zap 文档](https://github.com/uber-go/zap)

### Go 学习资源
- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example](https://gobyexample.com/)
- [Go 标准库](https://pkg.go.dev/std)
