package mongodb

import (
	"context"
	"fmt"

	"github.com/fustgo/fustgo2/pkg/plugin"
	"github.com/fustgo/fustgo2/pkg/record"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBReader MongoDB读取器
type MongoDBReader struct {
	plugin.BasePlugin
	client *mongo.Client
	config *Config
}

// Config MongoDB读取器配置
type Config struct {
	URI      string `json:"uri" validate:"required"`
	Database string `json:"database" validate:"required"`
	Collection string `json:"collection" validate:"required"`
	Query    string `json:"query"` // 查询条件（JSON格式）
}

// NewMongoDBReader 创建MongoDB读取器实例
func NewMongoDBReader() plugin.Reader {
	return &MongoDBReader{
		BasePlugin: plugin.BasePlugin{
			PluginInfo: plugin.Info{
				Name:        "mongodb-reader",
				Type:        plugin.TypeReader,
				Version:     "1.0.0",
				Author:      "FustGo Team",
				Description: "MongoDB数据库读取器插件",
			},
		},
	}
}

// Init 初始化插件
func (r *MongoDBReader) Init(config map[string]interface{}) error {
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

	// 连接MongoDB
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	r.client = client
	return nil
}

// validateConfig 验证配置
func (r *MongoDBReader) validateConfig(cfg *Config) error {
	if cfg.URI == "" {
		return fmt.Errorf("uri is required")
	}
	if cfg.Database == "" {
		return fmt.Errorf("database is required")
	}
	if cfg.Collection == "" {
		return fmt.Errorf("collection is required")
	}
	return nil
}

// Read 从MongoDB读取数据
func (r *MongoDBReader) Read(ctx context.Context, output chan<- *record.Record) error {
	defer close(output)

	// 获取集合
	collection := r.client.Database(r.config.Database).Collection(r.config.Collection)

	// 解析查询条件
	var filter interface{}
	if r.config.Query != "" {
		// 这里简化处理，实际应该解析JSON查询条件
		filter = r.config.Query
	}

	// 执行查询
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	defer cursor.Close(ctx)

	// 迭代结果
	for cursor.Next(ctx) {
		// 检查上下文是否取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// 解析文档
		var doc map[string]interface{}
		if err := cursor.Decode(&doc); err != nil {
			r.GetMetrics().ErrorCount++
			continue
		}

		// 创建记录
		rec := record.NewRecord(doc)

		// 发送到输出通道
		select {
		case output <- rec:
			r.GetMetrics().RecordsProcessed++
			r.GetMetrics().SuccessCount++
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	// 检查迭代错误
	if err := cursor.Err(); err != nil {
		return fmt.Errorf("cursor error: %w", err)
	}

	return nil
}

// GetProgress 获取读取进度
func (r *MongoDBReader) GetProgress() *plugin.Progress {
	return &plugin.Progress{
		ProcessedRecords: r.GetMetrics().RecordsProcessed,
		BytesProcessed:   r.GetMetrics().BytesProcessed,
	}
}

// Close 关闭插件
func (r *MongoDBReader) Close() error {
	if r.client != nil {
		return r.client.Disconnect(context.Background())
	}
	return nil
}

// Validate 验证配置
func (r *MongoDBReader) Validate() error {
	if r.config == nil {
		return fmt.Errorf("config is nil")
	}
	return r.validateConfig(r.config)
}

// GetInfo 获取插件信息
func (r *MongoDBReader) GetInfo() *plugin.Info {
	return &r.PluginInfo
}

// GetName 获取插件名称
func (r *MongoDBReader) GetName() string {
	return r.PluginInfo.Name
}

// GetType 获取插件类型
func (r *MongoDBReader) GetType() plugin.Type {
	return r.PluginInfo.Type
}