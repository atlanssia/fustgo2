package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/fustgo/fustgo2/pkg/plugin"
	"github.com/fustgo/fustgo2/pkg/record"
	_ "github.com/go-sql-driver/mysql"
)

// MySQLReader MySQL读取器
type MySQLReader struct {
	plugin.BasePlugin
	db     *sql.DB
	config *Config
}

// Config MySQL读取器配置
type Config struct {
	Host     string `json:"host" validate:"required"`
	Port     int    `json:"port" validate:"required,min=1,max=65535"`
	Database string `json:"database" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Table    string `json:"table" validate:"required"`
	Query    string `json:"query"` // 自定义查询SQL
}

// NewMySQLReader 创建MySQL读取器实例
func NewMySQLReader() plugin.Reader {
	return &MySQLReader{
		BasePlugin: plugin.BasePlugin{
			PluginInfo: plugin.Info{
				Name:        "mysql-reader",
				Type:        plugin.TypeReader,
				Version:     "1.0.0",
				Author:      "FustGo Team",
				Description: "MySQL数据库读取器插件",
			},
		},
	}
}

// Init 初始化插件
func (r *MySQLReader) Init(config map[string]interface{}) error {
	// 解析配置
	cfg := &Config{}
	if err := plugin.MapToStruct(config, cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// 验证配置
	if err := r.validateConfig(cfg); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	r.config = cfg

	// 构建连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", 
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	// 打开数据库连接
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	r.db = db
	return nil
}

// validateConfig 验证配置
func (r *MySQLReader) validateConfig(cfg *Config) error {
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
	if cfg.Table == "" && cfg.Query == "" {
		return fmt.Errorf("either table or query must be specified")
	}
	return nil
}

// Read 从MySQL读取数据
func (r *MySQLReader) Read(ctx context.Context, output chan<- *record.Record) error {
	defer close(output)

	// 构建查询SQL
	var query string
	if r.config.Query != "" {
		query = r.config.Query
	} else {
		query = fmt.Sprintf("SELECT * FROM %s", r.config.Table)
	}

	// 执行查询
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// 获取列信息
	columns, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("failed to get columns: %w", err)
	}

	// 读取数据
	var wg sync.WaitGroup
	for rows.Next() {
		// 检查上下文是否取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// 创建列值切片
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// 扫描行数据
		if err := rows.Scan(valuePtrs...); err != nil {
			r.GetMetrics().ErrorCount++
			continue
		}

		// 构建记录
		fields := make(map[string]interface{})
		for i, col := range columns {
			fields[col] = values[i]
		}

		rec := record.NewRecord(fields)
		
		// 发送到输出通道
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case output <- rec:
				r.GetMetrics().RecordsProcessed++
			case <-ctx.Done():
			}
		}()
	}

	// 等待所有记录发送完成
	wg.Wait()

	// 检查迭代错误
	if err := rows.Err(); err != nil {
		return fmt.Errorf("row iteration error: %w", err)
	}

	return nil
}

// GetProgress 获取读取进度
func (r *MySQLReader) GetProgress() *plugin.Progress {
	return &plugin.Progress{
		ProcessedRecords: r.GetMetrics().RecordsProcessed,
		BytesProcessed:   r.GetMetrics().BytesProcessed,
	}
}

// Close 关闭插件
func (r *MySQLReader) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}

// Validate 验证配置
func (r *MySQLReader) Validate() error {
	if r.config == nil {
		return fmt.Errorf("config is nil")
	}
	return r.validateConfig(r.config)
}

// GetInfo 获取插件信息
func (r *MySQLReader) GetInfo() *plugin.Info {
	return &r.PluginInfo
}

// GetName 获取插件名称
func (r *MySQLReader) GetName() string {
	return r.PluginInfo.Name
}

// GetType 获取插件类型
func (r *MySQLReader) GetType() plugin.Type {
	return r.PluginInfo.Type
}