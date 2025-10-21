package api

import (
	"net/http"

	"github.com/fustgo/fustgo2/internal/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SetupRouter 设置路由
func SetupRouter(cfg *config.Config, db *gorm.DB, logger *zap.Logger) *gin.Engine {
	r := gin.New()
	
	// 中间件
	r.Use(gin.Recovery())
	r.Use(LoggerMiddleware(logger))
	
	// CORS
	if cfg.Security.CORS.Enabled {
		r.Use(CORSMiddleware(cfg.Security.CORS))
	}
	
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
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
			jobs.GET("", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "List jobs",
				})
			})
			jobs.POST("", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Create job",
				})
			})
			jobs.GET("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Get job",
					"id":      c.Param("id"),
				})
			})
			jobs.PUT("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Update job",
					"id":      c.Param("id"),
				})
			})
			jobs.DELETE("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Delete job",
					"id":      c.Param("id"),
				})
			})
		}
		
		// Executions API
		executions := v1.Group("/executions")
		{
			executions.GET("", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "List executions",
				})
			})
			executions.GET("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Get execution",
					"id":      c.Param("id"),
				})
			})
		}
		
		// Connections API
		connections := v1.Group("/connections")
		{
			connections.GET("", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "List connections",
				})
			})
			connections.POST("", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Create connection",
				})
			})
			connections.GET("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Get connection",
					"id":      c.Param("id"),
				})
			})
			connections.PUT("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Update connection",
					"id":      c.Param("id"),
				})
			})
			connections.DELETE("/:id", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Delete connection",
					"id":      c.Param("id"),
				})
			})
			connections.POST("/:id/test", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Test connection",
					"id":      c.Param("id"),
				})
			})
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
