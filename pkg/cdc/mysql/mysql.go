package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fustgo/fustgo2/pkg/record"
	_ "github.com/go-sql-driver/mysql"
)

// MySQLCDC MySQL CDC监听器
type MySQLCDC struct {
	db        *sql.DB
	config    *Config
	lastCheck time.Time
}

// Config MySQL CDC配置
type Config struct {
	Host     string `json:"host" validate:"required"`
	Port     int    `json:"port" validate:"required,min=1,max=65535"`
	Database string `json:"database" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Table    string `json:"table" validate:"required"`
	// Binlog配置
	BinlogFile string `json:"binlog_file"`
	BinlogPos  int64  `json:"binlog_pos"`
	// 轮询配置
	PollInterval time.Duration `json:"poll_interval"` // 轮询间隔
}

// ChangeType 变更类型
type ChangeType string

const (
	ChangeTypeInsert ChangeType = "insert"
	ChangeTypeUpdate ChangeType = "update"
	ChangeTypeDelete ChangeType = "delete"
)

// ChangeRecord 变更记录
type ChangeRecord struct {
	Type      ChangeType     `json:"type"`
	Table     string         `json:"table"`
	Data      *record.Record `json:"data"`
	Timestamp time.Time      `json:"timestamp"`
}

// NewMySQLCDC 创建MySQL CDC监听器实例
func NewMySQLCDC(config *Config) (*MySQLCDC, error) {
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// 设置默认轮询间隔
	if config.PollInterval <= 0 {
		config.PollInterval = 1 * time.Second
	}

	// 构建连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", 
		config.Username, config.Password, config.Host, config.Port, config.Database)

	// 打开数据库连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	return &MySQLCDC{
		db:     db,
		config: config,
	}, nil
}

// validateConfig 验证配置
func validateConfig(cfg *Config) error {
	if cfg.Host == "" {
		return fmt.Errorf("host is required")
	}
	if cfg.Port <= 0 || cfg.Port > 65535 {
		return fmt.Errorf("invalid port: %d", cfg.Port)
	}
	if cfg.Database == "" {
		return fmt.Errorf("database is required")
	}
	if cfg.Username == "" {
		return fmt.Errorf("username is required")
	}
	if cfg.Table == "" {
		return fmt.Errorf("table is required")
	}
	return nil
}

// Start 启动CDC监听
func (c *MySQLCDC) Start(ctx context.Context, output chan<- *ChangeRecord) error {
	defer close(output)

	// 如果配置了Binlog位置，使用Binlog监听
	if c.config.BinlogFile != "" {
		return c.startBinlogListener(ctx, output)
	}

	// 否则使用轮询方式
	return c.startPollingListener(ctx, output)
}

// startBinlogListener 启动Binlog监听器
func (c *MySQLCDC) startBinlogListener(ctx context.Context, output chan<- *ChangeRecord) error {
	// TODO: 实现Binlog监听
	// 这需要使用第三方库如go-mysql来解析Binlog
	return fmt.Errorf("binlog listener not implemented yet")
}

// startPollingListener 启动轮询监听器
func (c *MySQLCDC) startPollingListener(ctx context.Context, output chan<- *ChangeRecord) error {
	ticker := time.NewTicker(c.config.PollInterval)
	defer ticker.Stop()

	// 初始化检查时间
	c.lastCheck = time.Now().Add(-c.config.PollInterval)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := c.pollChanges(ctx, output); err != nil {
				return fmt.Errorf("failed to poll changes: %w", err)
			}
		}
	}
}

// pollChanges 轮询变更
func (c *MySQLCDC) pollChanges(ctx context.Context, output chan<- *ChangeRecord) error {
	// 查询自上次检查以来的变更
	// 这里简化处理，实际应该使用时间戳字段或版本字段来跟踪变更
	query := fmt.Sprintf("SELECT * FROM %s WHERE updated_at > ?", c.config.Table)
	rows, err := c.db.QueryContext(ctx, query, c.lastCheck)
	if err != nil {
		return fmt.Errorf("failed to query changes: %w", err)
	}
	defer rows.Close()

	// 处理变更记录
	for rows.Next() {
		// 检查上下文是否取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// 获取列信息
		columns, err := rows.Columns()
		if err != nil {
			return fmt.Errorf("failed to get columns: %w", err)
		}

		// 创建列值切片
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// 扫描行数据
		if err := rows.Scan(valuePtrs...); err != nil {
			continue
		}

		// 构建记录
		fields := make(map[string]interface{})
		for i, col := range columns {
			fields[col] = values[i]
		}

		rec := &ChangeRecord{
			Type:      ChangeTypeUpdate, // 简化处理，假设都是更新
			Table:     c.config.Table,
			Data:      record.NewRecord(fields),
			Timestamp: time.Now(),
		}

		// 发送到输出通道
		select {
		case output <- rec:
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	// 更新检查时间
	c.lastCheck = time.Now()

	// 检查迭代错误
	if err := rows.Err(); err != nil {
		return fmt.Errorf("row iteration error: %w", err)
	}

	return nil
}

// Close 关闭CDC监听器
func (c *MySQLCDC) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}