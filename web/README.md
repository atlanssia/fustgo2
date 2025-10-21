# FustGo Web UI

现代化数据同步平台的前端界面。

## 技术栈

- **React 18** - UI 框架
- **TypeScript** - 类型安全
- **Vite** - 构建工具
- **Ant Design** - UI 组件库
- **Tailwind CSS** - 样式工具
- **React Router** - 路由管理
- **React Query** - 数据请求
- **Zustand** - 状态管理
- **Axios** - HTTP 客户端

## 项目结构

```
web/
├── src/
│   ├── layouts/          # 布局组件
│   │   └── MainLayout.tsx
│   ├── pages/            # 页面组件
│   │   ├── Dashboard.tsx
│   │   ├── Jobs.tsx
│   │   ├── Executions.tsx
│   │   └── Connections.tsx
│   ├── services/         # API 服务
│   │   ├── api.ts
│   │   ├── jobs.ts
│   │   ├── executions.ts
│   │   └── connections.ts
│   ├── App.tsx           # 根组件
│   ├── main.tsx          # 入口文件
│   └── index.css         # 全局样式
├── public/               # 静态资源
├── index.html            # HTML 模板
├── package.json          # 依赖配置
├── tsconfig.json         # TypeScript 配置
├── vite.config.ts        # Vite 配置
└── tailwind.config.js    # Tailwind 配置
```

## 快速开始

### 安装依赖

```bash
# 使用 pnpm (推荐)
pnpm install

# 或使用 npm
npm install

# 或使用 yarn
yarn install
```

### 开发服务器

```bash
pnpm dev
```

访问 http://localhost:3000

### 构建生产版本

```bash
pnpm build
```

构建产物在 `dist/` 目录。

### 预览生产版本

```bash
pnpm preview
```

### 代码检查

```bash
# ESLint 检查
pnpm lint

# 代码格式化
pnpm format
```

## 环境变量

复制 `.env.example` 为 `.env` 并配置:

```bash
cp .env.example .env
```

配置项:
- `VITE_API_BASE_URL`: API 基础地址 (默认: http://localhost:8080/api/v1)
- `VITE_APP_TITLE`: 应用标题

## 页面说明

### 仪表板 (Dashboard)
- 显示系统概览统计
- 任务运行状态汇总
- 快速操作入口

### 任务管理 (Jobs)
- 数据同步任务列表
- 创建/编辑/删除任务
- 立即执行任务
- 启用/禁用任务

### 执行历史 (Executions)
- 任务执行记录
- 执行状态和指标
- 错误日志查看

### 连接配置 (Connections)
- 数据源连接管理
- 连接测试
- 连接信息加密存储

## API 集成

所有 API 请求通过 `src/services/` 目录下的服务类进行:

```typescript
import { jobService } from '@/services/jobs'

// 获取任务列表
const jobs = await jobService.list()

// 创建任务
const newJob = await jobService.create({
  name: '任务名称',
  description: '任务描述',
  pipeline_config: '...',
})
```

## 开发规范

### 组件命名
- 使用 PascalCase: `MyComponent.tsx`
- 默认导出组件

### 类型定义
- 接口使用 `interface`
- 类型别名使用 `type`
- 导出公共类型

### 样式
- 优先使用 Tailwind CSS
- 复杂样式使用 CSS Modules
- 避免内联样式

### 状态管理
- 局部状态使用 `useState`
- 全局状态使用 Zustand
- 服务端状态使用 React Query

## 浏览器支持

- Chrome >= 90
- Firefox >= 88
- Safari >= 14
- Edge >= 90

## 构建优化

- 代码分割: React vendor, Ant Design vendor
- Tree Shaking: 自动移除未使用代码
- 资源压缩: Gzip/Brotli
- 懒加载: 路由级代码分割

## 常见问题

### 1. API 请求 CORS 错误?
开发环境已配置代理,确保后端服务运行在 `http://localhost:8080`

### 2. 依赖安装失败?
尝试清除缓存: `pnpm store prune` 或 `npm cache clean --force`

### 3. 热更新不生效?
重启开发服务器: `pnpm dev`

## License

Apache 2.0
