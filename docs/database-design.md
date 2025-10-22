# FustGo 数据库设计文档

## 概述

FustGo 使用关系型数据库存储元数据，包括任务配置、执行记录、连接信息等。默认使用 SQLite 以支持开发环境的快速启动，生产环境推荐使用 PostgreSQL。

## 数据库表结构

### 1. jobs - 数据同步任务表

存储数据同步任务的基本信息和配置。

```sql
CREATE TABLE jobs (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    enabled BOOLEAN DEFAULT true,
    status TEXT DEFAULT 'pending',
    pipeline_config TEXT,
    schedule_config TEXT,
    retry_config TEXT,
    notification_config TEXT,
    tags TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT,
    last_execution_id TEXT,
    last_execution_time TIMESTAMP,
    next_execution_time TIMESTAMP
);
```

字段说明：
- `id`: 任务唯一标识符
- `name`: 任务名称，唯一
- `description`: 任务描述
- `enabled`: 是否启用
- `status`: 任务状态 (pending, running, succeeded, failed, cancelled, paused)
- `pipeline_config`: 管道配置 (JSON格式)
- `schedule_config`: 调度配置 (JSON格式)
- `retry_config`: 重试配置 (JSON格式)
- `notification_config`: 通知配置 (JSON格式)
- `tags`: 标签 (JSON格式)
- `created_at`: 创建时间
- `updated_at`: 更新时间
- `created_by`: 创建人
- `last_execution_id`: 最后执行ID
- `last_execution_time`: 最后执行时间
- `next_execution_time`: 下次执行时间

### 2. job_executions - 任务执行记录表

存储任务执行的详细记录和指标。

```sql
CREATE TABLE job_executions (
    id TEXT PRIMARY KEY,
    job_id TEXT NOT NULL,
    status TEXT DEFAULT 'pending',
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
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    triggered_by TEXT
);
```

字段说明：
- `id`: 执行记录唯一标识符
- `job_id`: 关联的任务ID
- `status`: 执行状态 (pending, running, succeeded, failed, cancelled)
- `start_time`: 开始时间
- `end_time`: 结束时间
- `duration`: 执行时长(毫秒)
- `records_read`: 读取记录数
- `records_written`: 写入记录数
- `records_filtered`: 过滤记录数
- `records_error`: 错误记录数
- `bytes_read`: 读取字节数
- `bytes_written`: 写入字节数
- `error_message`: 错误信息
- `error_stack`: 错误堆栈
- `context`: 执行上下文 (JSON格式)
- `metrics`: 执行指标 (JSON格式)
- `created_at`: 创建时间
- `updated_at`: 更新时间
- `triggered_by`: 触发方式

### 3. connections - 数据源连接配置表

存储数据源连接配置信息。

```sql
CREATE TABLE connections (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    type TEXT NOT NULL,
    description TEXT,
    config TEXT,
    tags TEXT,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT,
    last_tested_at TIMESTAMP,
    last_test_status TEXT
);
```

字段说明：
- `id`: 连接配置唯一标识符
- `name`: 连接名称，唯一
- `type`: 连接类型 (mysql, postgresql, mongodb等)
- `description`: 连接描述
- `config`: 连接配置 (JSON格式，加密存储)
- `tags`: 标签 (JSON格式)
- `enabled`: 是否启用
- `created_at`: 创建时间
- `updated_at`: 更新时间
- `created_by`: 创建人
- `last_tested_at`: 最后测试时间
- `last_test_status`: 最后测试状态

### 4. plugins - 插件信息表

存储插件的基本信息。

```sql
CREATE TABLE plugins (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    type TEXT NOT NULL,
    version TEXT,
    description TEXT,
    author TEXT,
    config_schema TEXT,
    enabled BOOLEAN DEFAULT true,
    built_in BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

字段说明：
- `id`: 插件唯一标识符
- `name`: 插件名称，唯一
- `type`: 插件类型 (reader, writer, processor)
- `version`: 版本
- `description`: 描述
- `author`: 作者
- `config_schema`: 配置模式 (JSON Schema格式)
- `enabled`: 是否启用
- `built_in`: 是否内置插件
- `created_at`: 创建时间
- `updated_at`: 更新时间

### 5. plugin_instances - 插件实例配置表

存储插件实例的具体配置。

```sql
CREATE TABLE plugin_instances (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    plugin_id TEXT NOT NULL,
    config TEXT,
    description TEXT,
    tags TEXT,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT
);
```

字段说明：
- `id`: 实例唯一标识符
- `name`: 实例名称，唯一
- `plugin_id`: 关联的插件ID
- `config`: 配置信息 (JSON格式)
- `description`: 描述
- `tags`: 标签 (JSON格式)
- `enabled`: 是否启用
- `created_at`: 创建时间
- `updated_at`: 更新时间
- `created_by`: 创建人

### 6. pipelines - 管道配置表

存储数据管道的配置信息。

```sql
CREATE TABLE pipelines (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    status TEXT DEFAULT 'draft',
    config TEXT,
    schedule TEXT,
    tags TEXT,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by TEXT,
    last_execution_id TEXT,
    last_execution_time TIMESTAMP
);
```

字段说明：
- `id`: 管道唯一标识符
- `name`: 管道名称，唯一
- `description`: 描述
- `status`: 状态 (draft, active, inactive, archived)
- `config`: 管道配置 (JSON格式)
- `schedule`: 调度配置 (JSON格式)
- `tags`: 标签 (JSON格式)
- `enabled`: 是否启用
- `created_at`: 创建时间
- `updated_at`: 更新时间
- `created_by`: 创建人
- `last_execution_id`: 最后执行ID
- `last_execution_time`: 最后执行时间

## 索引设计

为提高查询性能，创建以下索引：

```sql
-- jobs 表索引
CREATE INDEX idx_jobs_status ON jobs(status);
CREATE INDEX idx_jobs_enabled ON jobs(enabled);
CREATE INDEX idx_jobs_next_execution ON jobs(next_execution_time);

-- job_executions 表索引
CREATE INDEX idx_executions_job_id ON job_executions(job_id);
CREATE INDEX idx_executions_status ON job_executions(status);
CREATE INDEX idx_executions_start_time ON job_executions(start_time);

-- connections 表索引
CREATE INDEX idx_connections_type ON connections(type);
CREATE INDEX idx_connections_enabled ON connections(enabled);

-- plugins 表索引
CREATE INDEX idx_plugins_type ON plugins(type);
CREATE INDEX idx_plugins_enabled ON plugins(enabled);

-- plugin_instances 表索引
CREATE INDEX idx_plugin_instances_plugin_id ON plugin_instances(plugin_id);
CREATE INDEX idx_plugin_instances_enabled ON plugin_instances(enabled);

-- pipelines 表索引
CREATE INDEX idx_pipelines_status ON pipelines(status);
CREATE INDEX idx_pipelines_enabled ON pipelines(enabled);
```

## 数据库配置

### SQLite (开发环境默认)

```yaml
database:
  type: sqlite
  sqlite:
    path: ./data/fustgo.db
```

### PostgreSQL (生产环境推荐)

```yaml
database:
  type: postgresql
  postgresql:
    host: localhost
    port: 5432
    database: fustgo
    username: fustgo
    password: ${DB_PASSWORD}
    ssl_mode: disable
    max_open_conns: 100
    max_idle_conns: 10
    conn_max_lifetime: 3600
```

## 数据库迁移

FustGo 使用 GORM 的自动迁移功能来管理数据库模式变更。在应用启动时会自动创建缺失的表和字段。

对于生产环境，建议使用专业的数据库迁移工具如 [Golang Migrate](https://github.com/golang-migrate/migrate) 来管理数据库变更。

## 备份与恢复

建议定期备份数据库文件：

### SQLite 备份
```bash
# 备份
cp ./data/fustgo.db ./backups/fustgo_$(date +%Y%m%d_%H%M%S).db

# 恢复
cp ./backups/fustgo_backup.db ./data/fustgo.db
```

### PostgreSQL 备份
```bash
# 备份
pg_dump -h localhost -U fustgo -d fustgo > ./backups/fustgo_$(date +%Y%m%d_%H%M%S).sql

# 恢复
psql -h localhost -U fustgo -d fustgo < ./backups/fustgo_backup.sql
```

## 性能优化建议

1. **索引优化**: 根据实际查询模式创建合适的索引
2. **连接池配置**: 合理配置数据库连接池大小
3. **查询优化**: 避免 N+1 查询问题，使用预加载
4. **分页查询**: 对于大量数据的查询使用分页
5. **定期维护**: 定期清理过期数据，优化表结构