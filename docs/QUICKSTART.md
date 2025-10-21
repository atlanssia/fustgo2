# FustGo 快速开始指南

## 📦 项目已完成的内容

### ✅ 架构设计文档
- `README.md` - 项目概览和介绍
- `docs/architecture.md` - 完整的系统架构设计
- `docs/module-interfaces.md` - 核心模块接口定义
- `docs/tech-selection.md` - 技术选型论证
- `docs/roadmap.md` - 详细的开发路线图
- `docs/project-structure.md` - 项目目录结构说明
- `docs/SUMMARY.md` - 架构设计总结

### ✅ 基础设施文件
- `go.mod` - Go 模块依赖定义
- `Makefile` - 构建和开发脚本
- `.gitignore` - Git 忽略规则
- `LICENSE` - Apache 2.0 开源许可证

### ✅ 部署配置
- `deployments/docker/Dockerfile` - 多阶段 Docker 构建配置
- `deployments/docker/docker-compose.yml` - 完整的容器编排配置

### ✅ 配置模板
- `configs/system.yaml` - 系统配置模板
- `configs/jobs/example-job.yaml` - 任务配置示例
- `configs/connections/example-conn.yaml` - 连接配置示例

### ✅ 目录结构
完整的项目目录已创建,包括:
- `cmd/` - 应用入口
- `internal/` - 私有代码
- `pkg/` - 公共库
- `web/` - 前端代码
- `configs/` - 配置文件
- `deployments/` - 部署配置
- `docs/` - 文档
- `scripts/` - 脚本

## 🚀 下一步开发指引

### Step 1: 初始化开发环境

```bash
# 1. 进入项目目录
cd /data/workspace/fustgo2

# 2. 下载 Go 依赖
go mod tidy

# 3. 安装前端依赖(如果还没有初始化前端项目)
# 首先需要创建前端项目,参考下面的 Step 2
```

### Step 2: 初始化前端项目

```bash
# 进入 web 目录
cd web

# 使用 Vite 创建 React + TypeScript 项目
pnpm create vite . --template react-ts

# 安装依赖
pnpm install

# 安装必要的库
pnpm add react-router-dom zustand @tanstack/react-query
pnpm add antd @ant-design/icons
pnpm add tailwindcss postcss autoprefixer
pnpm add -D @types/node

# 初始化 Tailwind CSS
npx tailwindcss init -p

# 返回项目根目录
cd ..
```

### Step 3: 开始开发 v0.1

参考 `docs/roadmap.md` 中的 v0.1 任务列表:

**后端核心任务:**
1. ✅ 项目目录结构设计 - 已完成
2. ☐ Go mod 初始化 - 执行 `go mod tidy`
3. ☐ Gin 框架集成
4. ☐ GORM ORM 集成
5. ☐ Zap 日志系统集成
6. ☐ Viper 配置管理集成
7. ☐ 插件系统框架开发
8. ☐ Pipeline 管道引擎骨架

**前端核心任务:**
1. ☐ React + Vite 项目初始化
2. ☐ TypeScript 配置
3. ☐ Tailwind CSS 集成
4. ☐ 路由配置
5. ☐ 基础布局组件

### Step 4: 运行开发服务器

```bash
# 后端开发服务器
make run

# 或者直接运行
go run ./cmd/server

# 前端开发服务器(另一个终端)
cd web
pnpm dev
```

### Step 5: 构建和打包

```bash
# 构建前端
make build-web

# 构建后端
make build

# 构建前后端(前端会嵌入到后端)
make build-all

# 运行打包后的二进制
./bin/fustgo server
```

## 📖 推荐阅读顺序

建议按以下顺序阅读文档,快速理解项目:

1. **README.md** - 了解项目概览和特性
2. **docs/SUMMARY.md** - 查看架构设计总结
3. **docs/architecture.md** - 深入理解系统架构
4. **docs/tech-selection.md** - 了解技术选型理由
5. **docs/module-interfaces.md** - 熟悉核心接口定义
6. **docs/roadmap.md** - 了解开发计划
7. **docs/project-structure.md** - 熟悉项目结构

## 🛠️ 开发工具推荐

### 必备工具
- **Go 1.21+** - 后端开发
- **Node.js 18+** - 前端开发
- **pnpm** - 前端包管理器
- **Docker** - 容器化部署
- **Docker Compose** - 本地开发环境

### 推荐的 IDE
- **VSCode** + Go 插件 + ESLint + Prettier
- **GoLand** (JetBrains)

### 有用的 VSCode 插件
- Go (官方)
- ESLint
- Prettier
- Tailwind CSS IntelliSense
- Docker
- YAML

## 📝 代码规范

### Go 代码规范
```bash
# 格式化代码
make fmt

# 代码检查
make lint

# 运行测试
make test
```

### 前端代码规范
```bash
cd web

# 格式化代码
pnpm format

# 代码检查
pnpm lint

# 运行测试
pnpm test
```

## 🐛 常见问题

### Q: go mod tidy 失败?
A: 确保 Go 版本 >= 1.21,并且网络可以访问 Go 模块代理

### Q: 前端依赖安装慢?
A: 使用 pnpm 并配置国内镜像:
```bash
pnpm config set registry https://registry.npmmirror.com
```

### Q: Docker 构建失败?
A: 确保先构建前端,再运行 Docker 构建

## 🎯 本周目标 (v0.1 Week 1)

1. ✅ 架构设计完成
2. ☐ 初始化 Go 项目
3. ☐ 初始化前端项目
4. ☐ 集成 Gin 框架
5. ☐ 实现基础 API 路由
6. ☐ 实现插件系统框架

## 📞 获取帮助

- 查看 `docs/` 目录下的详细文档
- 参考 `configs/` 目录下的配置示例
- 查看 `Makefile` 了解所有可用命令

## 🎉 开始开发吧!

所有的架构设计和基础设施已经准备就绪,现在可以开始实际的代码开发了。

建议从实现插件系统的基础接口开始,然后逐步完成 MySQL 和 PostgreSQL 插件。

祝开发顺利! 🚀
