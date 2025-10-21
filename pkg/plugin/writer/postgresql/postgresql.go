package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/fustgo/fustgo2/pkg/plugin"
	"github.com/fustgo/fustgo2/pkg/record"
	_ "github.com/lib/pq"
)

// PostgreSQLWriter PostgreSQL写入器
type PostgreSQLWriter struct {
	plugin.BasePlugin
	db     *sql.DB
	config *Config
}

// Config PostgreSQL写入器配置
type Config struct {
	Host     string `json:"host" validate:"required"`
	Port     int    `json:"port" validate:"required,min=1,max=65535"`
	Database string `json:"database" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Table    string `json:"table" validate:"required"`
	Mode     string `json:"mode"` // insert, upsert, update
}

// NewPostgreSQLWriter 创建PostgreSQL写入器实例
func NewPostgreSQLWriter() plugin.Writer {
	return &PostgreSQLWriter{
		BasePlugin: plugin.BasePlugin{
			PluginInfo: plugin.Info{
				Name:        "postgresql-writer",
				Type:        plugin.TypeWriter,
				Version:     "1.0.0",
				Author:      "FustGo Team",
				Description: "PostgreSQL数据库写入器插件",
			},
		},
	}
}

// Init 初始化插件
func (w *PostgreSQLWriter) Init(config map[string]interface{}) error {
	// 解析配置
	cfg := &Config{}
	if err := plugin.MapToStruct(config, cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// 验证配置
	if err := w.validateConfig(cfg); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// 设置默认模式
	if cfg.Mode == "" {
		cfg.Mode = "insert"
	}

	w.config = cfg

	// 构建连接字符串
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)

	// 打开数据库连接
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// 设置连接池参数
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	w.db = db
	return nil
}

// validateConfig 验证配置
func (w *PostgreSQLWriter) validateConfig(cfg *Config) error {
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
	if cfg.Mode != "" && cfg.Mode != "insert" && cfg.Mode != "upsert" && cfg.Mode != "update" {
		return fmt.Errorf("invalid mode: %s", cfg.Mode)
	}
	return nil
}

// Write 将数据写入PostgreSQL
func (w *PostgreSQLWriter) Write(ctx context.Context, input <-chan *record.Record) error {
	for {
		select {
		case rec, ok := <-input:
			if !ok {
				// 输入通道已关闭
				return nil
			}

			// 检查上下文是否取消
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			// 写入记录
			if err := w.writeRecord(ctx, rec); err != nil {
				w.GetMetrics().ErrorCount++
				// 根据错误策略决定是否继续
				if plugin.IsRetryable(err) {
					return err
				}
				// 非重试错误继续处理下一条记录
				continue
			}

			w.GetMetrics().RecordsProcessed++
			w.GetMetrics().SuccessCount++

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// writeRecord 写入单条记录
func (w *PostgreSQLWriter) writeRecord(ctx context.Context, rec *record.Record) error {
	// 获取字段信息
	fields := rec.GetFields()
	if len(fields) == 0 {
		return fmt.Errorf("record has no fields")
	}

	// 构建SQL语句
	var query string
	var args []interface{}

	switch w.config.Mode {
	case "insert":
		query, args = w.buildInsertSQL(fields)
	case "upsert":
		query, args = w.buildUpsertSQL(fields)
	case "update":
		query, args = w.buildUpdateSQL(fields)
	default:
		query, args = w.buildInsertSQL(fields)
	}

	// 执行SQL
	_, err := w.db.ExecContext(ctx, query, args...)
	if err != nil {
		return plugin.NewError(w.GetName(), "write", err, isRetryableError(err))
	}

	return nil
}

// buildInsertSQL 构建INSERT SQL
func (w *PostgreSQLWriter) buildInsertSQL(fields map[string]interface{}) (string, []interface{}) {
	columns := make([]string, 0, len(fields))
	placeholders := make([]string, 0, len(fields))
	args := make([]interface{}, 0, len(fields))

	i := 1
	for col, val := range fields {
		columns = append(columns, col)
		placeholders = append(placeholders, fmt.Sprintf("$%d", i))
		args = append(args, val)
		i++
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		w.config.Table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	return query, args
}

// buildUpsertSQL 构建UPSERT SQL
func (w *PostgreSQLWriter) buildUpsertSQL(fields map[string]interface{}) (string, []interface{}) {
	// 简化实现，实际应该根据主键进行UPSERT
	insertQuery, args := w.buildInsertSQL(fields)
	
	// PostgreSQL使用ON CONFLICT实现UPSERT
	// 这里简化处理，实际需要根据表的主键或唯一约束来构建
	query := insertQuery + " ON CONFLICT DO NOTHING"
	
	return query, args
}

// buildUpdateSQL 构建UPDATE SQL
func (w *PostgreSQLWriter) buildUpdateSQL(fields map[string]interface{}) (string, []interface{}) {
	// 简化实现，实际应该根据条件进行UPDATE
	setParts := make([]string, 0, len(fields))
	args := make([]interface{}, 0, len(fields))

	i := 1
	for col, val := range fields {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", col, i))
		args = append(args, val)
		i++
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d",
		w.config.Table,
		strings.Join(setParts, ", "),
		i)

	// 假设记录中有id字段用于WHERE条件
	if id, ok := fields["id"]; ok {
		args = append(args, id)
	} else {
		// 如果没有id字段，添加一个默认值
		args = append(args, 0)
	}

	return query, args
}

// isRetryableError 判断错误是否可重试
func isRetryableError(err error) bool {
	// 连接错误、超时错误等通常可重试
	if strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "timeout") ||
		strings.Contains(err.Error(), "deadlock") {
		return true
	}
	return false
}

// Flush 刷新缓冲区
func (w *PostgreSQLWriter) Flush() error {
	// PostgreSQL写入器不需要显式刷新
	return nil
}

// Close 关闭插件
func (w *PostgreSQLWriter) Close() error {
	if w.db != nil {
		return w.db.Close()
	}
	return nil
}

// Validate 验证配置
func (w *PostgreSQLWriter) Validate() error {
	if w.config == nil {
		return fmt.Errorf("config is nil")
	}
	return w.validateConfig(w.config)
}

// GetInfo 获取插件信息
func (w *PostgreSQLWriter) GetInfo() *plugin.Info {
	return &w.PluginInfo
}

// GetName 获取插件名称
func (w *PostgreSQLWriter) GetName() string {
	return w.PluginInfo.Name
}

// GetType 获取插件类型
func (w *PostgreSQLWriter) GetType() plugin.Type {
	return w.PluginInfo.Type
}