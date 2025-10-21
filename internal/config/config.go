package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Config 系统配置
type Config struct {
	// Server 服务器配置
	Server ServerConfig `mapstructure:"server"`
	
	// Database 数据库配置
	Database DatabaseConfig `mapstructure:"database"`
	
	// Cache 缓存配置
	Cache CacheConfig `mapstructure:"cache"`
	
	// Logging 日志配置
	Logging LoggingConfig `mapstructure:"logging"`
	
	// Executor 执行器配置
	Executor ExecutorConfig `mapstructure:"executor"`
	
	// Scheduler 调度器配置
	Scheduler SchedulerConfig `mapstructure:"scheduler"`
	
	// Monitoring 监控配置
	Monitoring MonitoringConfig `mapstructure:"monitoring"`
	
	// Security 安全配置
	Security SecurityConfig `mapstructure:"security"`
	
	// Plugins 插件配置
	Plugins PluginsConfig `mapstructure:"plugins"`
	
	// Storage 存储配置
	Storage StorageConfig `mapstructure:"storage"`
	
	// MessageQueue 消息队列配置
	MessageQueue MessageQueueConfig `mapstructure:"message_queue"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	// Port HTTP 端口
	Port int `mapstructure:"port"`
	
	// Mode 运行模式: debug, release, test
	Mode string `mapstructure:"mode"`
	
	// ReadTimeout 读取超时(秒)
	ReadTimeout int `mapstructure:"read_timeout"`
	
	// WriteTimeout 写入超时(秒)
	WriteTimeout int `mapstructure:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	// Type 数据库类型: sqlite, postgresql
	Type string `mapstructure:"type"`
	
	// SQLite SQLite 配置
	SQLite SQLiteConfig `mapstructure:"sqlite"`
	
	// PostgreSQL PostgreSQL 配置
	PostgreSQL PostgreSQLConfig `mapstructure:"postgresql"`
}

// SQLiteConfig SQLite 配置
type SQLiteConfig struct {
	// Path 数据库文件路径
	Path string `mapstructure:"path"`
}

// PostgreSQLConfig PostgreSQL 配置
type PostgreSQLConfig struct {
	// Host 主机地址
	Host string `mapstructure:"host"`
	
	// Port 端口
	Port int `mapstructure:"port"`
	
	// Database 数据库名
	Database string `mapstructure:"database"`
	
	// Username 用户名
	Username string `mapstructure:"username"`
	
	// Password 密码
	Password string `mapstructure:"password"`
	
	// SSLMode SSL 模式
	SSLMode string `mapstructure:"ssl_mode"`
	
	// MaxOpenConns 最大打开连接数
	MaxOpenConns int `mapstructure:"max_open_conns"`
	
	// MaxIdleConns 最大空闲连接数
	MaxIdleConns int `mapstructure:"max_idle_conns"`
	
	// ConnMaxLifetime 连接最大生命周期(秒)
	ConnMaxLifetime int `mapstructure:"conn_max_lifetime"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	// Type 缓存类型: memory, redis, file
	Type string `mapstructure:"type"`
	
	// Memory 内存缓存配置
	Memory MemoryCacheConfig `mapstructure:"memory"`
	
	// Redis Redis 配置
	Redis RedisCacheConfig `mapstructure:"redis"`
	
	// File 文件缓存配置
	File FileCacheConfig `mapstructure:"file"`
}

// MemoryCacheConfig 内存缓存配置
type MemoryCacheConfig struct {
	// MaxSizeMB 最大大小(MB)
	MaxSizeMB int `mapstructure:"max_size_mb"`
	
	// CleanupInterval 清理间隔(秒)
	CleanupInterval int `mapstructure:"cleanup_interval"`
}

// RedisCacheConfig Redis 缓存配置
type RedisCacheConfig struct {
	// Host 主机地址
	Host string `mapstructure:"host"`
	
	// Port 端口
	Port int `mapstructure:"port"`
	
	// Password 密码
	Password string `mapstructure:"password"`
	
	// DB 数据库索引
	DB int `mapstructure:"db"`
	
	// MaxRetries 最大重试次数
	MaxRetries int `mapstructure:"max_retries"`
}

// FileCacheConfig 文件缓存配置
type FileCacheConfig struct {
	// Path 缓存目录
	Path string `mapstructure:"path"`
	
	// MaxSizeMB 最大大小(MB)
	MaxSizeMB int `mapstructure:"max_size_mb"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	// Level 日志级别: debug, info, warn, error
	Level string `mapstructure:"level"`
	
	// Format 日志格式: json, console
	Format string `mapstructure:"format"`
	
	// Outputs 输出配置
	Outputs []LogOutputConfig `mapstructure:"outputs"`
}

// LogOutputConfig 日志输出配置
type LogOutputConfig struct {
	// Type 输出类型: stdout, file
	Type string `mapstructure:"type"`
	
	// Path 文件路径(type=file 时)
	Path string `mapstructure:"path"`
	
	// MaxSize 最大大小(MB)
	MaxSize int `mapstructure:"max_size"`
	
	// MaxBackups 最大备份数
	MaxBackups int `mapstructure:"max_backups"`
	
	// MaxAge 最大保留天数
	MaxAge int `mapstructure:"max_age"`
	
	// Compress 是否压缩
	Compress bool `mapstructure:"compress"`
}

// ExecutorConfig 执行器配置
type ExecutorConfig struct {
	// MaxConcurrentJobs 最大并发任务数
	MaxConcurrentJobs int `mapstructure:"max_concurrent_jobs"`
	
	// MaxGoroutinesPerJob 单任务最大协程数
	MaxGoroutinesPerJob int `mapstructure:"max_goroutines_per_job"`
	
	// DefaultTimeout 默认超时时间(秒)
	DefaultTimeout int `mapstructure:"default_timeout"`
	
	// WorkerPool Worker Pool 配置
	WorkerPool WorkerPoolConfig `mapstructure:"worker_pool"`
}

// WorkerPoolConfig Worker Pool 配置
type WorkerPoolConfig struct {
	// MinWorkers 最小工作协程数
	MinWorkers int `mapstructure:"min_workers"`
	
	// MaxWorkers 最大工作协程数
	MaxWorkers int `mapstructure:"max_workers"`
	
	// QueueSize 队列大小
	QueueSize int `mapstructure:"queue_size"`
	
	// IdleTimeout 空闲超时(秒)
	IdleTimeout int `mapstructure:"idle_timeout"`
}

// SchedulerConfig 调度器配置
type SchedulerConfig struct {
	// Enabled 是否启用
	Enabled bool `mapstructure:"enabled"`
	
	// Timezone 时区
	Timezone string `mapstructure:"timezone"`
	
	// CheckInterval 检查间隔(秒)
	CheckInterval int `mapstructure:"check_interval"`
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
	// Enabled 是否启用
	Enabled bool `mapstructure:"enabled"`
	
	// Metrics 指标配置
	Metrics MetricsConfig `mapstructure:"metrics"`
	
	// OpenObserve OpenObserve 配置
	OpenObserve OpenObserveConfig `mapstructure:"openobserve"`
	
	// Tracing 追踪配置
	Tracing TracingConfig `mapstructure:"tracing"`
}

// MetricsConfig 指标配置
type MetricsConfig struct {
	// Enabled 是否启用
	Enabled bool `mapstructure:"enabled"`
	
	// Interval 采集间隔(秒)
	Interval int `mapstructure:"interval"`
}

// OpenObserveConfig OpenObserve 配置
type OpenObserveConfig struct {
	// Enabled 是否启用
	Enabled bool `mapstructure:"enabled"`
	
	// Endpoint 端点地址
	Endpoint string `mapstructure:"endpoint"`
	
	// Organization 组织
	Organization string `mapstructure:"organization"`
	
	// Stream 流名称
	Stream string `mapstructure:"stream"`
	
	// Username 用户名
	Username string `mapstructure:"username"`
	
	// Password 密码
	Password string `mapstructure:"password"`
}

// TracingConfig 追踪配置
type TracingConfig struct {
	// Enabled 是否启用
	Enabled bool `mapstructure:"enabled"`
	
	// Endpoint 端点地址
	Endpoint string `mapstructure:"endpoint"`
	
	// SampleRate 采样率
	SampleRate float64 `mapstructure:"sample_rate"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	// JWT JWT 配置
	JWT JWTConfig `mapstructure:"jwt"`
	
	// EncryptionKey 加密密钥
	EncryptionKey string `mapstructure:"encryption_key"`
	
	// CORS CORS 配置
	CORS CORSConfig `mapstructure:"cors"`
}

// JWTConfig JWT 配置
type JWTConfig struct {
	// Secret 密钥
	Secret string `mapstructure:"secret"`
	
	// Expiration 过期时间(秒)
	Expiration int `mapstructure:"expiration"`
}

// CORSConfig CORS 配置
type CORSConfig struct {
	// Enabled 是否启用
	Enabled bool `mapstructure:"enabled"`
	
	// AllowedOrigins 允许的源
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	
	// AllowedMethods 允许的方法
	AllowedMethods []string `mapstructure:"allowed_methods"`
	
	// AllowedHeaders 允许的头
	AllowedHeaders []string `mapstructure:"allowed_headers"`
}

// PluginsConfig 插件配置
type PluginsConfig struct {
	// ConfigDir 配置目录
	ConfigDir string `mapstructure:"config_dir"`
	
	// ConnectionPool 连接池默认配置
	ConnectionPool ConnectionPoolConfig `mapstructure:"connection_pool"`
}

// ConnectionPoolConfig 连接池配置
type ConnectionPoolConfig struct {
	// MaxOpenConns 最大打开连接数
	MaxOpenConns int `mapstructure:"max_open_conns"`
	
	// MaxIdleConns 最大空闲连接数
	MaxIdleConns int `mapstructure:"max_idle_conns"`
	
	// ConnMaxLifetime 连接最大生命周期(秒)
	ConnMaxLifetime int `mapstructure:"conn_max_lifetime"`
}

// StorageConfig 存储配置
type StorageConfig struct {
	// Type 类型: local, s3, minio
	Type string `mapstructure:"type"`
	
	// Local 本地存储配置
	Local LocalStorageConfig `mapstructure:"local"`
	
	// S3 S3 配置
	S3 S3Config `mapstructure:"s3"`
}

// LocalStorageConfig 本地存储配置
type LocalStorageConfig struct {
	// Path 存储路径
	Path string `mapstructure:"path"`
}

// S3Config S3 配置
type S3Config struct {
	// Endpoint 端点地址
	Endpoint string `mapstructure:"endpoint"`
	
	// AccessKey 访问密钥
	AccessKey string `mapstructure:"access_key"`
	
	// SecretKey 密钥
	SecretKey string `mapstructure:"secret_key"`
	
	// Bucket 存储桶
	Bucket string `mapstructure:"bucket"`
	
	// Region 区域
	Region string `mapstructure:"region"`
	
	// UseSSL 是否使用 SSL
	UseSSL bool `mapstructure:"use_ssl"`
}

// MessageQueueConfig 消息队列配置
type MessageQueueConfig struct {
	// Enabled 是否启用
	Enabled bool `mapstructure:"enabled"`
	
	// NATS NATS 配置
	NATS NATSConfig `mapstructure:"nats"`
}

// NATSConfig NATS 配置
type NATSConfig struct {
	// URL 连接地址
	URL string `mapstructure:"url"`
	
	// MaxReconnects 最大重连次数
	MaxReconnects int `mapstructure:"max_reconnects"`
	
	// ReconnectWait 重连等待时间(秒)
	ReconnectWait int `mapstructure:"reconnect_wait"`
}

// Load 加载配置
func Load(configPath string) (*Config, error) {
	v := viper.New()
	
	// 设置配置文件
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("system")
		v.SetConfigType("yaml")
		v.AddConfigPath("./configs")
		v.AddConfigPath(".")
	}
	
	// 启用环境变量
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	
	// 设置默认值
	setDefaults(v)
	
	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		// 如果配置文件不存在,使用默认值
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}
	
	// 解析配置
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	return &config, nil
}

// setDefaults 设置默认值
func setDefaults(v *viper.Viper) {
	// Server
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.mode", "release")
	v.SetDefault("server.read_timeout", 60)
	v.SetDefault("server.write_timeout", 60)
	
	// Database
	v.SetDefault("database.type", "sqlite")
	v.SetDefault("database.sqlite.path", "./data/fustgo.db")
	v.SetDefault("database.postgresql.port", 5432)
	v.SetDefault("database.postgresql.ssl_mode", "disable")
	v.SetDefault("database.postgresql.max_open_conns", 100)
	v.SetDefault("database.postgresql.max_idle_conns", 10)
	v.SetDefault("database.postgresql.conn_max_lifetime", 3600)
	
	// Cache
	v.SetDefault("cache.type", "memory")
	v.SetDefault("cache.memory.max_size_mb", 512)
	v.SetDefault("cache.memory.cleanup_interval", 300)
	v.SetDefault("cache.redis.port", 6379)
	v.SetDefault("cache.redis.db", 0)
	v.SetDefault("cache.redis.max_retries", 3)
	
	// Logging
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")
	
	// Executor
	v.SetDefault("executor.max_concurrent_jobs", 50)
	v.SetDefault("executor.max_goroutines_per_job", 100)
	v.SetDefault("executor.default_timeout", 86400)
	v.SetDefault("executor.worker_pool.min_workers", 10)
	v.SetDefault("executor.worker_pool.max_workers", 100)
	v.SetDefault("executor.worker_pool.queue_size", 1000)
	v.SetDefault("executor.worker_pool.idle_timeout", 300)
	
	// Scheduler
	v.SetDefault("scheduler.enabled", true)
	v.SetDefault("scheduler.timezone", "Asia/Shanghai")
	v.SetDefault("scheduler.check_interval", 10)
	
	// Monitoring
	v.SetDefault("monitoring.enabled", true)
	v.SetDefault("monitoring.metrics.enabled", true)
	v.SetDefault("monitoring.metrics.interval", 30)
	v.SetDefault("monitoring.tracing.sample_rate", 0.1)
	
	// Security
	v.SetDefault("security.jwt.expiration", 86400)
	v.SetDefault("security.cors.enabled", true)
	
	// Plugins
	v.SetDefault("plugins.config_dir", "./configs/plugins")
	v.SetDefault("plugins.connection_pool.max_open_conns", 50)
	v.SetDefault("plugins.connection_pool.max_idle_conns", 10)
	v.SetDefault("plugins.connection_pool.conn_max_lifetime", 3600)
	
	// Storage
	v.SetDefault("storage.type", "local")
	v.SetDefault("storage.local.path", "./data/files")
	
	// MessageQueue
	v.SetDefault("message_queue.enabled", false)
	v.SetDefault("message_queue.nats.max_reconnects", 10)
	v.SetDefault("message_queue.nats.reconnect_wait", 2)
}

// GetDuration 获取时间间隔配置
func (c *Config) GetDuration(seconds int) time.Duration {
	return time.Duration(seconds) * time.Second
}
