# FustGo 项目目录结构

## 整体结构

```
fustgo2/
├── cmd/                          # 应用程序入口
│   ├── server/                   # 主服务器
│   │   └── main.go              # 服务器入口
│   └── cli/                      # CLI 命令行工具
│       └── main.go              # CLI 入口
│
├── internal/                     # 私有应用代码(不可被外部导入)
│   ├── api/                      # HTTP API 层
│   │   ├── handler/             # 请求处理器
│   │   ├── middleware/          # 中间件
│   │   ├── router/              # 路由配置
│   │   └── dto/                 # 数据传输对象
│   │
│   ├── core/                     # 核心业务逻辑
│   │   ├── pipeline/            # 数据管道引擎
│   │   ├── executor/            # 任务执行器
│   │   └── registry/            # 插件注册表
│   │
│   ├── scheduler/                # 任务调度器
│   │   ├── cron/                # Cron 调度
│   │   ├── event/               # 事件调度
│   │   └── dag/                 # DAG 工作流
│   │
│   ├── plugin/                   # 插件管理
│   │   ├── loader/              # 插件加载器
│   │   ├── manager/             # 插件管理器
│   │   └── validator/           # 插件验证器
│   │
│   └── storage/                  # 数据存储层
│       ├── model/               # 数据模型
│       ├── repository/          # 数据访问层
│       └── migration/           # 数据库迁移
│
├── pkg/                          # 公共库(可被外部导入)
│   ├── connector/                # 连接器
│   │   ├── mysql/               # MySQL 连接器
│   │   ├── postgresql/          # PostgreSQL 连接器
│   │   ├── mongodb/             # MongoDB 连接器
│   │   └── base/                # 基础连接器接口
│   │
│   ├── reader/                   # 数据读取器
│   │   ├── mysql/               # MySQL 读取器
│   │   ├── postgresql/          # PostgreSQL 读取器
│   │   ├── batch/               # 批量读取器
│   │   ├── stream/              # 流式读取器
│   │   └── base/                # 基础读取器接口
│   │
│   ├── writer/                   # 数据写入器
│   │   ├── mysql/               # MySQL 写入器
│   │   ├── postgresql/          # PostgreSQL 写入器
│   │   ├── batch/               # 批量写入器
│   │   └── base/                # 基础写入器接口
│   │
│   ├── transformer/              # 数据转换器
│   │   ├── mapper/              # 字段映射
│   │   ├── filter/              # 过滤器
│   │   ├── converter/           # 类型转换
│   │   └── base/                # 基础转换器接口
│   │
│   ├── record/                   # 数据记录
│   │   └── record.go            # Record 定义
│   │
│   ├── config/                   # 配置管理
│   │   ├── loader/              # 配置加载器
│   │   ├── provider/            # 配置提供者
│   │   └── validator/           # 配置验证器
│   │
│   └── monitor/                  # 监控
│       ├── metrics/             # 指标收集
│       ├── logger/              # 日志
│       └── tracer/              # 链路追踪
│
├── web/                          # 前端代码
│   ├── src/                      # 源代码
│   │   ├── components/          # 公共组件
│   │   ├── pages/               # 页面
│   │   ├── hooks/               # 自定义 Hooks
│   │   ├── stores/              # 状态管理
│   │   ├── services/            # API 服务
│   │   ├── utils/               # 工具函数
│   │   ├── types/               # TypeScript 类型
│   │   ├── App.tsx              # 应用根组件
│   │   └── main.tsx             # 应用入口
│   │
│   ├── public/                   # 静态资源
│   ├── dist/                     # 构建产物
│   ├── index.html               # HTML 模板
│   ├── package.json             # 前端依赖
│   ├── tsconfig.json            # TypeScript 配置
│   ├── vite.config.ts           # Vite 配置
│   └── tailwind.config.js       # Tailwind 配置
│
├── configs/                      # 配置文件
│   ├── system.yaml              # 系统配置
│   ├── jobs/                    # 任务配置目录
│   │   └── example-job.yaml    # 示例任务配置
│   ├── connections/             # 连接配置目录
│   │   └── example-conn.yaml   # 示例连接配置
│   └── plugins/                 # 插件配置目录
│       └── example-plugin.yaml # 示例插件配置
│
├── deployments/                  # 部署配置
│   ├── docker/                  # Docker 配置
│   │   ├── Dockerfile           # Dockerfile
│   │   └── docker-compose.yml  # Docker Compose
│   └── k8s/                     # Kubernetes 配置
│       ├── deployment.yaml     # 部署配置
│       ├── service.yaml        # 服务配置
│       └── configmap.yaml      # 配置映射
│
├── docs/                         # 文档
│   ├── architecture.md          # 架构设计
│   ├── module-interfaces.md    # 模块接口定义
│   ├── tech-selection.md       # 技术选型
│   ├── roadmap.md              # 开发路线图
│   ├── api-reference.md        # API 参考
│   ├── plugin-development.md  # 插件开发指南
│   └── deployment.md           # 部署指南
│
├── scripts/                      # 脚本
│   ├── build.sh                # 构建脚本
│   ├── test.sh                 # 测试脚本
│   └── release.sh              # 发布脚本
│
├── test/                         # 测试
│   ├── integration/            # 集成测试
│   ├── e2e/                    # 端到端测试
│   └── fixtures/               # 测试数据
│
├── .gitignore                   # Git 忽略文件
├── .golangci.yml               # Go Linter 配置
├── go.mod                       # Go 模块定义
├── go.sum                       # Go 依赖锁定
├── Makefile                     # 构建命令
├── LICENSE                      # 许可证
└── README.md                    # 项目说明
```

## 目录说明

### cmd/ - 应用入口
存放应用程序的入口点。每个子目录对应一个可执行程序:
- `server/`: Web 服务器,提供 API 和 UI
- `cli/`: 命令行工具,用于任务管理和执行

### internal/ - 私有代码
应用程序的私有代码,不能被外部项目导入:
- `api/`: HTTP API 层,处理 Web 请求
- `core/`: 核心业务逻辑,包括管道引擎、执行器等
- `scheduler/`: 任务调度器
- `plugin/`: 插件管理系统
- `storage/`: 数据存储层,包括数据库模型和访问层

### pkg/ - 公共库
可以被外部项目导入的公共库:
- `connector/`: 数据源连接器实现
- `reader/`: 数据读取器实现
- `writer/`: 数据写入器实现
- `transformer/`: 数据转换器实现
- `record/`: 数据记录定义
- `config/`: 配置管理
- `monitor/`: 监控、日志、追踪

### web/ - 前端代码
React + TypeScript 前端应用:
- `src/`: 源代码
- `public/`: 静态资源
- `dist/`: 构建产物(被 Go embed)

### configs/ - 配置文件
系统配置和任务配置:
- `system.yaml`: 系统级配置
- `jobs/`: 任务配置文件
- `connections/`: 数据源连接配置
- `plugins/`: 插件配置

### deployments/ - 部署配置
部署相关的配置文件:
- `docker/`: Docker 和 Docker Compose 配置
- `k8s/`: Kubernetes 部署配置

### docs/ - 文档
项目文档:
- 架构设计文档
- API 参考文档
- 开发指南
- 部署指南

### scripts/ - 脚本
构建、测试、发布脚本

## 代码组织原则

1. **清晰的分层**: API → Service → Repository → Database
2. **依赖倒置**: 高层模块依赖接口,不依赖具体实现
3. **单一职责**: 每个包只负责一个明确的功能
4. **包的独立性**: 包之间通过接口交互,减少耦合
5. **测试友好**: 便于编写单元测试和集成测试

## 包导入规则

```go
// 标准库
import (
    "context"
    "fmt"
)

// 第三方库
import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

// 项目内部包
import (
    "github.com/fustgo/fustgo2/internal/core/pipeline"
    "github.com/fustgo/fustgo2/pkg/reader"
)
```

## 命名约定

- **包名**: 小写,单数,简短(如 `reader`, `writer`)
- **文件名**: 小写,下划线分隔(如 `mysql_reader.go`)
- **接口名**: 大写开头,名词(如 `Reader`, `Connector`)
- **实现名**: 大写开头,描述性(如 `MySQLReader`, `PostgreSQLConnector`)

## 下一步

参考 `docs/roadmap.md` 开始 v0.1 的开发工作。
