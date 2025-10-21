module github.com/fustgo/fustgo2

go 1.21

require (
	// Web Framework
	github.com/gin-gonic/gin v1.10.0
	
	// Database
	gorm.io/gorm v1.25.5
	gorm.io/driver/sqlite v1.5.4
	gorm.io/driver/postgres v1.5.4
	gorm.io/driver/mysql v1.5.2
	
	// Configuration
	github.com/spf13/viper v1.18.2
	
	// Logging
	go.uber.org/zap v1.26.0
	
	// Scheduling
	github.com/robfig/cron/v3 v3.0.1
	
	// UUID
	github.com/google/uuid v1.5.0
	
	// Validation
	github.com/go-playground/validator/v10 v10.16.0
)
