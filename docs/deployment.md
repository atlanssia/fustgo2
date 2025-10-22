# FustGo 部署指南

## 概述

FustGo 支持多种部署方式，包括 Docker 容器化部署和二进制部署。推荐使用 Docker 方式进行生产环境部署，使用开发环境配置进行本地开发。

## 系统要求

### 最低配置
- CPU: 2 核
- 内存: 4 GB
- 磁盘: 10 GB 可用空间

### 推荐配置
- CPU: 4 核或更多
- 内存: 8 GB 或更多
- 磁盘: 50 GB 或更多 SSD 存储

## 部署方式

### 1. Docker 部署（推荐）

FustGo 提供了完整的 Docker 部署方案，包括开发环境和生产环境两种配置。

#### 开发环境部署（单机版，无外部依赖）

开发环境使用 SQLite 数据库和内存缓存，适合本地开发和测试。

```bash
# 克隆仓库
git clone https://github.com/yourusername/fustgo2.git
cd fustgo2

# 启动开发环境
docker-compose -f deployments/docker/docker-compose.dev.yml up -d

# 访问应用
open http://localhost:8080
```

#### 生产环境部署（完整版，需要外部依赖）

生产环境使用 PostgreSQL 数据库和 Redis 缓存，适合生产环境部署。

```bash
# 克隆仓库
git clone https://github.com/yourusername/fustgo2.git
cd fustgo2

# 启动生产环境
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 2. 二进制部署

#### 下载预编译版本

```bash
# 下载最新版本
wget https://github.com/yourusername/fustgo2/releases/latest/download/fustgo-linux-amd64.tar.gz

# 解压
tar -xzf fustgo-linux-amd64.tar.gz

# 运行
./fustgo server --config config.yaml
```

#### 从源码编译

```bash
# 克隆仓库
git clone https://github.com/yourusername/fustgo2.git
cd fustgo2

# 编译前端（可选，已有预编译版本）
cd web && npm install && npm run build && cd ..

# 编译后端
make build

# 运行
./bin/fustgo server --config configs/system.yaml
```

## 配置说明

### 系统配置文件

FustGo 使用 YAML 格式的配置文件，支持环境变量覆盖。

#### 开发环境配置 (configs/system.dev.yaml)

```yaml
# 服务器配置
server:
  port: 8080
  mode: debug
  read_timeout: 60
  write_timeout: 60

# 数据库配置 (SQLite)
database:
  type: sqlite
  sqlite:
    path: ./data/fustgo_dev.db

# 缓存配置 (内存)
cache:
  type: memory
  memory:
    max_size_mb: 512
    cleanup_interval: 300

# 日志配置
logging:
  level: debug
  format: console
  outputs:
    - type: stdout
```

#### 生产环境配置 (configs/system.yaml)

```yaml
# 服务器配置
server:
  port: 8080
  mode: release
  read_timeout: 60
  write_timeout: 60

# 数据库配置 (PostgreSQL)
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

# 缓存配置 (Redis)
cache:
  type: redis
  redis:
    host: localhost
    port: 6379
    password: ${REDIS_PASSWORD}
    db: 0
    max_retries: 3

# 日志配置
logging:
  level: info
  format: json
  outputs:
    - type: stdout
    - type: file
      path: ./logs/fustgo.log
      max_size: 100
      max_backups: 7
      max_age: 30
      compress: true
```

## 环境变量

FustGo 支持通过环境变量覆盖配置文件中的值：

| 环境变量 | 说明 | 默认值 |
|---------|------|-------|
| SERVER_PORT | 服务器端口 | 8080 |
| DB_TYPE | 数据库类型 | sqlite |
| DB_HOST | 数据库主机 | localhost |
| DB_PORT | 数据库端口 | 5432 |
| DB_NAME | 数据库名 | fustgo |
| DB_USER | 数据库用户名 | fustgo |
| DB_PASSWORD | 数据库密码 |  |
| REDIS_HOST | Redis主机 | localhost |
| REDIS_PORT | Redis端口 | 6379 |
| REDIS_PASSWORD | Redis密码 |  |
| LOG_LEVEL | 日志级别 | info |

## 数据库初始化

### SQLite (开发环境)

SQLite 数据库会自动创建，无需手动初始化。

### PostgreSQL (生产环境)

```sql
-- 创建数据库
CREATE DATABASE fustgo;

-- 创建用户
CREATE USER fustgo WITH PASSWORD 'your_password';

-- 授权
GRANT ALL PRIVILEGES ON DATABASE fustgo TO fustgo;
```

## 监控和日志

### 日志查看

```bash
# Docker 环境
docker-compose logs -f fustgo

# 二进制部署
tail -f logs/fustgo.log
```

### 健康检查

```bash
# 检查服务状态
curl http://localhost:8080/health

# 返回示例
{
  "status": "healthy",
  "version": "v1.0.0"
}
```

## 备份和恢复

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

## 故障排除

### 常见问题

1. **端口被占用**
   ```bash
   # 查看端口占用
   lsof -i :8080
   
   # 杀死占用进程
   kill -9 <PID>
   ```

2. **数据库连接失败**
   - 检查数据库服务是否启动
   - 检查数据库配置是否正确
   - 检查网络连接

3. **权限问题**
   ```bash
   # 确保数据目录有写权限
   chmod 755 ./data/
   ```

### 日志分析

查看日志文件以获取更多错误信息：

```bash
# 查看错误日志
grep -i error logs/fustgo.log

# 查看警告日志
grep -i warn logs/fustgo.log
```

## 升级指南

### Docker 升级

```bash
# 拉取最新镜像
docker-compose pull

# 重启服务
docker-compose up -d
```

### 二进制升级

```bash
# 停止服务
pkill fustgo

# 下载新版本
wget https://github.com/yourusername/fustgo2/releases/latest/download/fustgo-linux-amd64.tar.gz

# 替换二进制文件
tar -xzf fustgo-linux-amd64.tar.gz
mv fustgo /usr/local/bin/

# 启动服务
fustgo server --config config.yaml
```

## 安全建议

1. **使用强密码**
   - 为数据库用户设置强密码
   - 为 Redis 设置密码认证

2. **网络隔离**
   - 不要将数据库端口暴露到公网
   - 使用防火墙限制访问

3. **定期更新**
   - 定期更新 FustGo 到最新版本
   - 定期更新操作系统和依赖组件

4. **备份策略**
   - 定期备份数据库
   - 保留多个时间点的备份

## 性能优化

1. **数据库优化**
   - 根据负载调整连接池大小
   - 定期清理过期数据
   - 创建合适的索引

2. **缓存优化**
   - 合理设置缓存大小
   - 使用 Redis 而不是内存缓存处理大量数据

3. **资源监控**
   - 监控 CPU 和内存使用情况
   - 监控磁盘空间和 I/O 性能