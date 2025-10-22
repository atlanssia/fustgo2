package storage

import (
	"fmt"
	"time"

	"github.com/fustgo/fustgo2/internal/config"
	"github.com/fustgo/fustgo2/internal/core/types"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitDB 初始化数据库连接
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector
	
	// 根据配置选择数据库
	switch cfg.Database.Type {
	case "sqlite":
		dialector = sqlite.Open(cfg.Database.SQLite.Path)
	case "postgresql":
		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Database.PostgreSQL.Host,
			cfg.Database.PostgreSQL.Port,
			cfg.Database.PostgreSQL.Username,
			cfg.Database.PostgreSQL.Password,
			cfg.Database.PostgreSQL.Database,
			cfg.Database.PostgreSQL.SSLMode,
		)
		dialector = postgres.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
	}
	
	// GORM 配置
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}
	
	// 打开数据库连接
	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	
	// 获取底层 SQL DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}
	
	// 配置连接池(仅 PostgreSQL)
	if cfg.Database.Type == "postgresql" {
		sqlDB.SetMaxOpenConns(cfg.Database.PostgreSQL.MaxOpenConns)
		sqlDB.SetMaxIdleConns(cfg.Database.PostgreSQL.MaxIdleConns)
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.PostgreSQL.ConnMaxLifetime) * time.Second)
	}
	
	// 自动迁移数据库表
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}
	
	return db, nil
}

// autoMigrate 自动迁移数据库表
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&types.Job{},
		&types.Execution{},
		&types.Connection{},
		&types.Plugin{},
		&types.PluginInstance{},
		&types.Pipeline{},
	)
}
