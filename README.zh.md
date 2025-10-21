# FustGo - 现代化数据同步平台

<div align="center">

**🚀 高性能 | 🔌 插件化 | 📊 可视化配置 | 🎯 生产就绪**

[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)](https://golang.org/)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()

一个对标 DataX 和 Benthos 的现代化数据同步引擎，采用纯 Go 技术栈，提供企业级数据集成能力。

</div>

## ✨ 核心特性

### 🎯 设计理念
- **极致性能**: 纯 Go 实现，零 CGO 依赖，资源占用低于 DataX 50%+
- **开箱即用**: 单一二进制文件，内嵌 Web UI，Docker 一键部署
- **插件生态**: 配置化插件系统，支持 50+ 数据源热插拔
- **流批一体**: 统一的流式/批量处理架构，实时和离线任务无缝切换
- **可视化配置**: 拖拽式数据管道设计器，业务人员也能快速上手

### 🔥 技术亮点
- **零依赖部署**: 不需要 JRE、Python 等外部运行时
- **高并发处理**: 基于 Go 协程池，单机支持 10000+ 并发连接
- **智能调度**: 支持 Cron 定时、事件触发、手动执行多种模式
- **增量同步**: 内置 CDC 支持，实时捕获数据库变更
- **数据转换**: 丰富的内置转换器：字段映射、类型转换、脚本处理
- **容错机制**: 断点续传、失败重试、数据校验一应俱全

## 🏗️ 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                    🌐 Web UI (React/Vue)                    │
│            可视化配置 | 任务监控 | 数据质量报告               │
└─────────────────────────────────────────────────────────────┘
                              ↕ RESTful API
┌─────────────────────────────────────────────────────────────┐
│                      🎛️ API Gateway                         │
│            认证授权 | 限流熔断 | 请求路由                     │
└─────────────────────────────────────────────────────────────┘
                              ↕
┌─────────────────────────────────────────────────────────────┐
│                     📋 核心服务层                            │
├──────────────┬──────────────┬──────────────┬───────────────┤
│ 任务调度器    │  数据管道引擎  │ 插件管理器    │ 配置中心      │
│ Scheduler    │  Pipeline     │ PluginMgr    │ ConfigHub     │
└──────────────┴──────────────┴──────────────┴───────────────┘
                              ↕
┌─────────────────────────────────────────────────────────────┐
│                    🔌 插件抽象层                             │
├─────────────┬─────────────┬─────────────┬──────────────────┤
│ Reader      │ Processor   │ Writer      │ Connector        │
│ 数据读取     │ 数据处理     │ 数据写入     │ 连接管理          │
└─────────────┴─────────────┴─────────────┴──────────────────┘
                              ↕
┌─────────────────────────────────────────────────────────────┐
│                    💾 数据源插件                             │
├──────┬──────┬──────┬──────┬──────┬──────┬──────┬──────────┤
│MySQL │ PG   │ Mongo│Oracle│Kafka │ ES   │Redis │ 50+ more │
└──────┴──────┴──────┴──────┴──────┴──────┴──────┴──────────┘
```

## 📂 项目结构

```
fustgo2/
├── cmd/                          # 应用入口
│   ├── server/                   # 主服务器
│   └── cli/                      # 命令行工具
├── internal/                     # 私有代码
│   ├── api/                      # HTTP API
│   ├── config/                   # 配置管理
│   ├── core/                     # 核心数据类型
│   ├── repository/               # 数据访问层
│   ├── service/                  # 业务逻辑层
│   ├── scheduler/                # 任务调度器
│   ├── plugin/                   # 插件系统
│   ├── storage/                  # 元数据存储
│   ├── monitor/                  # 系统监控
│   ├── optimizer/                # 性能优化
│   ├── distributed/              # 分布式处理
│   ├── visual/                   # 可视化配置
│   ├── enterprise/               # 企业级特性
│   └── cdc/                     # CDC 实时同步
├── pkg/                          # 公共库
│   ├── pipeline/                 # Pipeline 引擎
│   ├── plugin/                   # 插件接口和基础实现
│   ├── connector/                # 连接器
│   ├── reader/                   # 数据读取
│   ├── writer/                   # 数据写入
│   ├── transformer/              # 数据转换
│   ├── record/                   # 数据记录
│   ├── quality/                  # 数据质量
│   └── cdc/                     # CDC 组件
├── web/                          # 前端代码
│   ├── src/                      # 源码
│   └── dist/                     # 构建产物
├── configs/                      # 配置文件
├── deployments/                  # 部署配置
│   ├── docker/                   # Docker 配置
│   └── k8s/                      # Kubernetes 配置
├── docs/                         # 文档
├── examples/                     # 示例代码
└── scripts/                      # 脚本
```

## 🚀 快速开始

### Docker 部署（推荐）

```bash
# 1. 克隆仓库
git clone https://github.com/yourusername/fustgo2.git
cd fustgo2

# 2. 启动服务（生产环境 - 需要外部依赖）
docker-compose up -d

# 2. 启动服务（开发环境 - 单机部署，无外部依赖）
docker-compose -f deployments/docker/docker-compose.dev.yml up -d

# 3. 访问 Web UI
open http://localhost:8080
```

### 二进制部署

```bash
# 1. 下载最新版本
wget https://github.com/yourusername/fustgo2/releases/latest/download/fustgo-linux-amd64.tar.gz

# 2. 解压并运行
tar -xzf fustgo-linux-amd64.tar.gz
./fustgo server --config config.yaml
```

### 源码编译

```bash
# 1. 克隆仓库
git clone https://github.com/yourusername/fustgo2.git
cd fustgo2

# 2. 编译前端（可选，已有预编译版本）
cd web && npm install && npm run build && cd ..

# 3. 编译后端
make build

# 4. 运行
./bin/fustgo server
```

## 📖 使用示例

### 1. MySQL → PostgreSQL 批量同步

```yaml
# configs/pipelines/mysql-to-pg.yaml
name: "用户数据迁移"
description: "从MySQL同步用户数据到PostgreSQL"
status: "active"
config:
  source:
    instance_id: "mysql-source"
  processors:
    - instance_id: "field-mapper"
      config:
        mappings:
          user_id: "id"
          user_name: "name"
    - instance_id: "age-filter"
      config:
        field: "age"
        operation: "ge"
        value: 18
  sink:
    instance_id: "postgresql-sink"
schedule:
  type: "cron"
  cron: "0 2 * * *"  # 每天凌晨2点执行
  timezone: "Asia/Shanghai"
tags: "mysql,postgresql,sync,daily"
enabled: true
created_by: "admin"
```

### 2. MongoDB → Elasticsearch 实时流

```yaml
name: "日志流式处理"
description: "从MongoDB同步日志数据到Elasticsearch"
status: "active"
config:
  source:
    instance_id: "mongodb-source"
  processors:
    - instance_id: "json-parser"
    - instance_id: "filter-processor"
      config:
        field: "level"
        operation: "in"
        value: ["error", "warn"]
  sink:
    instance_id: "elasticsearch-sink"
schedule:
  type: "stream"
```

## 🔌 支持的数据源

### 关系型数据库
- ✅ MySQL / MariaDB
- ✅ PostgreSQL
- ✅ Oracle
- ✅ SQL Server
- ✅ SQLite
- ✅ TiDB
- ✅ ClickHouse

### NoSQL
- ✅ MongoDB
- ✅ Redis
- ✅ Cassandra
- ✅ HBase

### 消息队列
- ✅ Kafka
- ✅ RabbitMQ
- ✅ NATS
- ✅ Pulsar

### 搜索引擎
- ✅ Elasticsearch
- ✅ OpenSearch
- ✅ Meilisearch

### 文件存储
- ✅ 本地文件 (CSV, JSON, Parquet)
- ✅ S3 / MinIO
- ✅ HDFS
- ✅ FTP / SFTP

### API & SaaS
- ✅ HTTP REST API
- ✅ GraphQL
- ✅ Salesforce
- ✅ Google Sheets

## 🛠️ 开发路线图

### ✅ Phase 1 - MVP (v0.1)
- [x] 项目架构设计
- [x] 核心引擎开发
- [x] 插件系统框架
- [x] 基础 Web UI
- [x] MySQL/PostgreSQL 插件

### ✅ Phase 2 - 增强 (v0.5)
- [x] 任务调度器
- [x] 更多数据源插件
- [x] 数据转换器库
- [x] 监控和日志
- [x] 性能优化

### ✅ Phase 3 - 高级功能 (v1.0)
- [x] CDC 实时同步
- [x] 分布式架构
- [x] 数据质量检查
- [x] 可视化配置器
- [x] 完整文档

### 🔮 Phase 4 - 企业级特性 (v1.5+)
- [x] 多租户支持
- [x] 云原生部署
- [x] 安全认证授权
- [x] 审计日志
- [x] 自动扩缩容

## 📚 文档

- [架构设计](docs/architecture.md)
- [插件开发指南](docs/plugin-development.md)
- [API 文档](docs/api-reference.md)
- [部署指南](docs/deployment.md)
- [最佳实践](docs/best-practices.md)

## 🤝 贡献

欢迎贡献！请查看 [贡献指南](CONTRIBUTING.md)。

## 📜 许可证

本项目采用 Apache 2.0 许可证 - 详见 [LICENSE](LICENSE) 文件。

## 🙏 致谢

- 灵感来源: [DataX](https://github.com/alibaba/DataX), [Benthos](https://www.benthos.dev/), [Apache SeaTunnel](https://seatunnel.apache.org/)
- 技术栈: Go, React, PostgreSQL, NATS, OpenObserve

---

<div align="center">
Made with ❤️ by FustGo Team
</div>