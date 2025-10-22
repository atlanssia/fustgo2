package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fustgo/fustgo2/internal/api"
	"github.com/fustgo/fustgo2/internal/config"
	"github.com/fustgo/fustgo2/internal/monitor"
	"github.com/fustgo/fustgo2/internal/plugin"
	"github.com/fustgo/fustgo2/internal/repository"
	"github.com/fustgo/fustgo2/internal/storage"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Version 版本号
	Version = "dev"
	
	// BuildTime 构建时间
	BuildTime = "unknown"
	
	// configFile 配置文件路径
	configFile = flag.String("config", "", "配置文件路径")
	
	// showVersion 显示版本信息
	showVersion = flag.Bool("version", false, "显示版本信息")
)

func main() {
	flag.Parse()
	
	// 显示版本信息
	if *showVersion {
		fmt.Printf("FustGo Version: %s\n", Version)
		fmt.Printf("Build Time: %s\n", BuildTime)
		os.Exit(0)
	}
	
	// 加载配置
	cfg, err := config.Load(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	
	// 初始化日志
	logger := initLogger(cfg)
	defer logger.Sync()
	
	logger.Info("Starting FustGo server",
		zap.String("version", Version),
		zap.String("build_time", BuildTime),
	)
	
	// 初始化数据库
	db, err := storage.InitDB(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize database", zap.Error(err))
	}
	
	logger.Info("Database initialized successfully")
	
	// 注册内置插件
	if err := plugin.RegisterPlugins(); err != nil {
		logger.Error("Failed to register plugins", zap.Error(err))
	} else {
		logger.Info("Plugins registered successfully")
	}
	
	// 初始化仓库
	jobRepo := repository.NewJobRepository(db)
	executionRepo := repository.NewExecutionRepository(db)
	pipelineRepo := repository.NewPipelineRepository(db)
	
	// 初始化服务
	jobService := service.NewJobService(jobRepo)
	
	// 初始化监控器
	monitor := monitor.NewMonitor(jobRepo, executionRepo, pipelineRepo, logger)
	if err := monitor.Start(); err != nil {
		logger.Error("Failed to start monitor", zap.Error(err))
	} else {
		logger.Info("Monitor started successfully")
	}
	defer monitor.Stop()
	
	// 初始化调度器
	// scheduler := scheduler.NewScheduler(jobRepo, jobService, pipelineRepo, logger)
	// if err := scheduler.Start(); err != nil {
	// 	logger.Error("Failed to start scheduler", zap.Error(err))
	// } else {
	// 	logger.Info("Scheduler started successfully")
	// }
	// defer scheduler.Stop()
	
	// 初始化优化器
	// optimizer := optimizer.NewOptimizer(jobRepo, executionRepo, pipelineRepo, logger)
	// if err := optimizer.Start(); err != nil {
	// 	logger.Error("Failed to start optimizer", zap.Error(err))
	// } else {
	// 	logger.Info("Optimizer started successfully")
	// }
	// defer optimizer.Stop()
	
	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)
	
	// 创建路由
	router := api.SetupRouter(cfg, db, logger)
	
	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
	}
	
	// 启动服务器(goroutine)
	go func() {
		logger.Info("Starting HTTP server",
			zap.Int("port", cfg.Server.Port),
			zap.String("mode", cfg.Server.Mode),
		)
		
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()
	
	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	logger.Info("Shutting down server...")
	
	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}
	
	logger.Info("Server exited")
}

// initLogger 初始化日志
func initLogger(cfg *config.Config) *zap.Logger {
	// 日志级别
	level := zapcore.InfoLevel
	switch cfg.Logging.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	}
	
	// 日志配置
	zapConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(level),
		Development:      false,
		Encoding:         cfg.Logging.Format,
		EncoderConfig:    getEncoderConfig(cfg.Logging.Format),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	
	// 添加文件输出
	for _, output := range cfg.Logging.Outputs {
		if output.Type == "file" && output.Path != "" {
			zapConfig.OutputPaths = append(zapConfig.OutputPaths, output.Path)
		}
	}
	
	logger, err := zapConfig.Build()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	
	return logger
}

// getEncoderConfig 获取编码配置
func getEncoderConfig(format string) zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	
	if format == "console" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	
	return encoderConfig
}
