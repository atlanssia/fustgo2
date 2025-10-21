package elasticsearch

import (
	"context"
	"fmt"
	"strings"

	"github.com/fustgo/fustgo2/pkg/plugin"
	"github.com/fustgo/fustgo2/pkg/record"
	elastic "github.com/olivere/elastic/v7"
)

// ElasticsearchWriter Elasticsearch写入器
type ElasticsearchWriter struct {
	plugin.BasePlugin
	client *elastic.Client
	config *Config
}

// Config Elasticsearch写入器配置
type Config struct {
	Addresses []string `json:"addresses" validate:"required"`
	Index     string   `json:"index" validate:"required"`
	Username  string   `json:"username"`
	Password  string   `json:"password"`
	BulkSize  int      `json:"bulk_size"` // 批量大小
}

// NewElasticsearchWriter 创建Elasticsearch写入器实例
func NewElasticsearchWriter() plugin.Writer {
	return &ElasticsearchWriter{
		BasePlugin: plugin.BasePlugin{
			PluginInfo: plugin.Info{
				Name:        "elasticsearch-writer",
				Type:        plugin.TypeWriter,
				Version:     "1.0.0",
				Author:      "FustGo Team",
				Description: "Elasticsearch搜索引擎写入器插件",
			},
		},
	}
}

// Init 初始化插件
func (w *ElasticsearchWriter) Init(config map[string]interface{}) error {
	// 解析配置
	cfg := &Config{}
	if err := plugin.MapToStruct(config, cfg); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// 验证配置
	if err := w.validateConfig(cfg); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	// 设置默认批量大小
	if cfg.BulkSize <= 0 {
		cfg.BulkSize = 1000
	}

	w.config = cfg

	// 连接Elasticsearch
	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(cfg.Addresses...),
		elastic.SetSniff(false),
	}

	// 如果提供了认证信息
	if cfg.Username != "" && cfg.Password != "" {
		opts = append(opts, elastic.SetBasicAuth(cfg.Username, cfg.Password))
	}

	client, err := elastic.NewClient(opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to Elasticsearch: %w", err)
	}

	w.client = client
	return nil
}

// validateConfig 验证配置
func (w *ElasticsearchWriter) validateConfig(cfg *Config) error {
	if len(cfg.Addresses) == 0 {
		return fmt.Errorf("addresses is required")
	}
	if cfg.Index == "" {
		return fmt.Errorf("index is required")
	}
	return nil
}

// Write 将数据写入Elasticsearch
func (w *ElasticsearchWriter) Write(ctx context.Context, input <-chan *record.Record) error {
	// 创建批量请求
	bulkRequest := w.client.Bulk()

	for {
		select {
		case rec, ok := <-input:
			if !ok {
				// 输入通道已关闭，刷新剩余数据
				if bulkRequest.NumberOfActions() > 0 {
					if err := w.flushBulk(ctx, bulkRequest); err != nil {
						return err
					}
				}
				return nil
			}

			// 检查上下文是否取消
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			// 添加到批量请求
			doc := rec.GetFields()
			req := elastic.NewBulkIndexRequest().
				Index(w.config.Index).
				Doc(doc)
			bulkRequest.Add(req)

			// 如果达到批量大小，执行批量请求
			if bulkRequest.NumberOfActions() >= w.config.BulkSize {
				if err := w.flushBulk(ctx, bulkRequest); err != nil {
					return err
				}
			}

			w.GetMetrics().RecordsProcessed++
			w.GetMetrics().SuccessCount++

		case <-ctx.Done():
			// 刷新剩余数据
			if bulkRequest.NumberOfActions() > 0 {
				if err := w.flushBulk(ctx, bulkRequest); err != nil {
					return err
				}
			}
			return ctx.Err()
		}
	}
}

// flushBulk 刷新批量请求
func (w *ElasticsearchWriter) flushBulk(ctx context.Context, bulkRequest *elastic.BulkService) error {
	// 执行批量请求
	result, err := bulkRequest.Do(ctx)
	if err != nil {
		w.GetMetrics().ErrorCount++
		return plugin.NewError(w.GetName(), "bulk_write", err, isRetryableError(err))
	}

	// 检查是否有错误
	if result.Errors {
		for _, item := range result.Items {
			for _, res := range item {
				if res.Error != nil {
					w.GetMetrics().ErrorCount++
					// 记录错误但继续处理
					continue
				}
			}
		}
	}

	return nil
}

// isRetryableError 判断错误是否可重试
func isRetryableError(err error) bool {
	// 连接错误、超时错误等通常可重试
	if strings.Contains(err.Error(), "connection refused") ||
		strings.Contains(err.Error(), "timeout") ||
		strings.Contains(err.Error(), "no such host") {
		return true
	}
	return false
}

// Flush 刷新缓冲区
func (w *ElasticsearchWriter) Flush() error {
	// Elasticsearch写入器在Write方法中已经处理了刷新逻辑
	return nil
}

// Close 关闭插件
func (w *ElasticsearchWriter) Close() error {
	if w.client != nil {
		w.client.Stop()
	}
	return nil
}

// Validate 验证配置
func (w *ElasticsearchWriter) Validate() error {
	if w.config == nil {
		return fmt.Errorf("config is nil")
	}
	return w.validateConfig(w.config)
}

// GetInfo 获取插件信息
func (w *ElasticsearchWriter) GetInfo() *plugin.Info {
	return &w.PluginInfo
}

// GetName 获取插件名称
func (w *ElasticsearchWriter) GetName() string {
	return w.PluginInfo.Name
}

// GetType 获取插件类型
func (w *ElasticsearchWriter) GetType() plugin.Type {
	return w.PluginInfo.Type
}