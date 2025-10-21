.PHONY: help build run test clean docker

# 默认目标
help:
	@echo "FustGo 构建命令:"
	@echo "  make build       - 构建后端二进制文件"
	@echo "  make build-web   - 构建前端"
	@echo "  make build-all   - 构建前端和后端"
	@echo "  make run         - 运行服务器"
	@echo "  make run-dev     - 运行开发环境服务器"
	@echo "  make test        - 运行测试"
	@echo "  make lint        - 代码检查"
	@echo "  make clean       - 清理构建产物"
	@echo "  make docker      - 构建 Docker 镜像"
	@echo "  make dev-setup   - 初始化开发环境"
	@echo "  make dev         - 启动开发环境 (Docker)"
	@echo "  make prod        - 启动生产环境 (Docker)"

# 变量
APP_NAME=fustgo
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

# 构建后端
build:
	@echo "Building backend..."
	CGO_ENABLED=0 go build ${LDFLAGS} -o bin/${APP_NAME} ./cmd/server

# 构建前端
build-web:
	@echo "Building frontend..."
	cd web && pnpm install && pnpm build

# 构建前后端
build-all: build-web build

# 运行服务器
run:
	@echo "Running server..."
	go run ./cmd/server

# 运行开发环境服务器
run-dev:
	@echo "Running development server..."
	go run ./cmd/server --config configs/system.yaml

# 运行测试
test:
	@echo "Running tests..."
	go test -v -cover ./...

# 代码检查
lint:
	@echo "Running linter..."
	golangci-lint run

# 格式化代码
fmt:
	@echo "Formatting code..."
	go fmt ./...
	goimports -w .

# 清理构建产物
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -rf web/dist/
	rm -rf dist/

# 构建 Docker 镜像
docker:
	@echo "Building Docker image..."
	docker build -t ${APP_NAME}:${VERSION} -f deployments/docker/Dockerfile .

# 多平台构建
build-all-platforms:
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ${LDFLAGS} -o dist/${APP_NAME}-linux-amd64 ./cmd/server
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build ${LDFLAGS} -o dist/${APP_NAME}-linux-arm64 ./cmd/server
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build ${LDFLAGS} -o dist/${APP_NAME}-darwin-amd64 ./cmd/server
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build ${LDFLAGS} -o dist/${APP_NAME}-darwin-arm64 ./cmd/server
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build ${LDFLAGS} -o dist/${APP_NAME}-windows-amd64.exe ./cmd/server

# 初始化开发环境
dev-setup:
	@echo "Setting up development environment..."
	go mod download
	cd web && pnpm install
	@echo "Development environment ready!"

# 启动开发环境 (Docker)
dev:
	@echo "Starting development environment..."
	docker-compose -f deployments/docker/docker-compose.dev.yml up -d
	@echo "Development environment started! Access at http://localhost:8080"

# 启动生产环境 (Docker)
prod:
	@echo "Starting production environment..."
	docker-compose up -d
	@echo "Production environment started! Access at http://localhost:8080"

# 数据库迁移
migrate:
	@echo "Running database migrations..."
	go run ./cmd/server migrate

# 生成 API 文档
docs:
	@echo "Generating API documentation..."
	swag init -g cmd/server/main.go
