# 多阶段构建 - 前端
FROM node:18-alpine AS frontend-builder

WORKDIR /app/web

# 安装 pnpm
RUN npm install -g pnpm

# 复制前端依赖文件
COPY web/package.json web/pnpm-lock.yaml* ./

# 安装依赖
RUN pnpm install --frozen-lockfile

# 复制前端源码
COPY web/ ./

# 构建前端
RUN pnpm build

# 多阶段构建 - 后端
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app

# 安装必要的工具
RUN apk add --no-cache git make

# 复制 go mod 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源码
COPY . .

# 复制前端构建产物
COPY --from=frontend-builder /app/web/dist ./web/dist

# 构建后端
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o fustgo \
    ./cmd/server

# 最终镜像
FROM alpine:latest

# 安装 CA 证书和时区数据
RUN apk --no-cache add ca-certificates tzdata

# 创建应用用户
RUN addgroup -g 1000 fustgo && \
    adduser -D -u 1000 -G fustgo fustgo

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=backend-builder /app/fustgo /usr/local/bin/fustgo

# 创建必要的目录
RUN mkdir -p /app/configs /app/data /app/logs && \
    chown -R fustgo:fustgo /app

# 切换到应用用户
USER fustgo

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动命令
ENTRYPOINT ["fustgo"]
CMD ["server"]
