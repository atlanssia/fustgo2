package api

import (
	"github.com/fustgo/fustgo2/internal/api/handlers"
	"github.com/fustgo/fustgo2/internal/config"
	"github.com/fustgo/fustgo2/internal/repository"
	"github.com/fustgo/fustgo2/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SetupRouter 设置路由
func SetupRouter(cfg *config.Config, db *gorm.DB, logger *zap.Logger) *gin.Engine {
	r := gin.New()
	
	// 中间件
	r.Use(gin.Recovery())
	r.Use(handlers.LoggerMiddleware(logger))
	
	// CORS
	if cfg.Security.CORS.Enabled {
		r.Use(handlers.CORSMiddleware(cfg.Security.CORS))
	}
	
	// 初始化仓库
	jobRepo := repository.NewJobRepository(db)
	executionRepo := repository.NewExecutionRepository(db)
	connectionRepo := repository.NewConnectionRepository(db)
	pluginRepo := repository.NewPluginRepository(db)
	pluginInstanceRepo := repository.NewPluginInstanceRepository(db)
	pipelineRepo := repository.NewPipelineRepository(db)
	
	// 初始化服务
	jobService := service.NewJobService(jobRepo)
	executionService := service.NewExecutionService(executionRepo)
	connectionService := service.NewConnectionService(connectionRepo)
	pluginService := service.NewPluginService(pluginRepo)
	pluginInstanceService := service.NewPluginInstanceService(pluginInstanceRepo)
	pipelineService := service.NewPipelineService(pipelineRepo)
	
	// 初始化处理程序
	jobHandler := handlers.NewJobHandler(jobService, logger)
	executionHandler := handlers.NewExecutionHandler(executionService, logger)
	connectionHandler := handlers.NewConnectionHandler(connectionService, logger)
	pluginHandler := handlers.NewPluginHandler(pluginService, pluginInstanceService, pluginRepo, pluginInstanceRepo, logger)
	pipelineHandler := handlers.NewPipelineHandler(pipelineService, pipelineRepo, pluginInstanceRepo, logger)
	
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"version": "dev",
		})
	})
	
	// API 路由组
	v1 := r.Group("/api/v1")
	{
		// Jobs API
		jobs := v1.Group("/jobs")
		{
			jobs.GET("", jobHandler.List)
			jobs.POST("", jobHandler.Create)
			jobs.GET("/:id", jobHandler.GetByID)
			jobs.PUT("/:id", jobHandler.Update)
			jobs.DELETE("/:id", jobHandler.Delete)
		}
		
		// Executions API
		executions := v1.Group("/executions")
		{
			executions.GET("", executionHandler.List)
			executions.POST("", executionHandler.Create)
			executions.GET("/:id", executionHandler.GetByID)
			executions.PUT("/:id", executionHandler.Update)
			executions.GET("/job/:job_id", executionHandler.ListByJobID)
		}
		
		// Connections API
		connections := v1.Group("/connections")
		{
			connections.GET("", connectionHandler.List)
			connections.POST("", connectionHandler.Create)
			connections.GET("/:id", connectionHandler.GetByID)
			connections.GET("/name/:name", connectionHandler.GetByName)
			connections.PUT("/:id", connectionHandler.Update)
			connections.DELETE("/:id", connectionHandler.Delete)
			connections.POST("/:id/test", connectionHandler.Test)
		}
		
		// Plugins API
		plugins := v1.Group("/plugins")
		{
			plugins.GET("", pluginHandler.ListPlugins)
			plugins.GET("/:id", pluginHandler.GetPluginByID)
		}
		
		// Plugin Instances API
		pluginInstances := v1.Group("/plugin-instances")
		{
			pluginInstances.GET("", pluginHandler.ListPluginInstances)
			pluginInstances.POST("", pluginHandler.CreatePluginInstance)
		}
		
		// Pipelines API
		pipelines := v1.Group("/pipelines")
		{
			pipelines.GET("", pipelineHandler.List)
			pipelines.POST("", pipelineHandler.Create)
			pipelines.GET("/:id", pipelineHandler.GetByID)
			pipelines.PUT("/:id", pipelineHandler.Update)
			pipelines.DELETE("/:id", pipelineHandler.Delete)
			pipelines.POST("/:id/execute", pipelineHandler.Execute)
		}
	}
	
	return r
}

// LoggerMiddleware 日志中间件
func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := c.Request.Context().Value("start")
		if start == nil {
			start = c.GetTime("start")
		}
		
		c.Next()
		
		logger.Info("HTTP Request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.String("ip", c.ClientIP()),
		)
	}
}

// CORSMiddleware CORS 中间件
func CORSMiddleware(cfg config.CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
