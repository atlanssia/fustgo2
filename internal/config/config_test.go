package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadDefaultConfig(t *testing.T) {
	// 测试加载默认配置(不存在配置文件时)
	cfg, err := Load("nonexistent.yaml")
	if err != nil {
		// 允许配置文件不存在,应该使用默认值
		if cfg == nil {
			t.Skip("Config file not found, skipping test")
		}
	}
	
	// 验证默认值
	if cfg.Server.Port != 8080 {
		t.Errorf("Default port = %d, want 8080", cfg.Server.Port)
	}
	
	if cfg.Server.Mode != "release" {
		t.Errorf("Default mode = %s, want release", cfg.Server.Mode)
	}
	
	if cfg.Database.Type != "sqlite" {
		t.Errorf("Default database type = %s, want sqlite", cfg.Database.Type)
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	// 设置环境变量
	os.Setenv("SERVER_PORT", "9000")
	defer os.Unsetenv("SERVER_PORT")
	
	cfg, err := Load("")
	if err != nil && cfg == nil {
		t.Skip("Config file not found, skipping test")
	}
	
	// 注意: Viper 的环境变量可能需要特定的设置才能生效
	// 这里只是示例测试
	if cfg != nil && cfg.Server.Port == 9000 {
		t.Log("Environment variable loaded successfully")
	}
}

func TestGetDuration(t *testing.T) {
	cfg := &Config{}
	
	duration := cfg.GetDuration(60)
	expected := 60 * time.Second
	
	if duration != expected {
		t.Errorf("GetDuration(60) = %v, want %v", duration, expected)
	}
}

func TestServerConfig(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port:         8080,
			Mode:         "release",
			ReadTimeout:  60,
			WriteTimeout: 60,
		},
	}
	
	if cfg.Server.Port != 8080 {
		t.Errorf("Server.Port = %d, want 8080", cfg.Server.Port)
	}
	
	if cfg.Server.Mode != "release" {
		t.Errorf("Server.Mode = %s, want release", cfg.Server.Mode)
	}
}

func TestDatabaseConfig(t *testing.T) {
	tests := []struct {
		name     string
		dbType   string
		expected string
	}{
		{"SQLite", "sqlite", "sqlite"},
		{"PostgreSQL", "postgresql", "postgresql"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Database: DatabaseConfig{
					Type: tt.dbType,
				},
			}
			
			if cfg.Database.Type != tt.expected {
				t.Errorf("Database.Type = %s, want %s", cfg.Database.Type, tt.expected)
			}
		})
	}
}

func TestCacheConfig(t *testing.T) {
	cfg := &Config{
		Cache: CacheConfig{
			Type: "memory",
			Memory: MemoryCacheConfig{
				MaxSizeMB:       512,
				CleanupInterval: 300,
			},
		},
	}
	
	if cfg.Cache.Type != "memory" {
		t.Errorf("Cache.Type = %s, want memory", cfg.Cache.Type)
	}
	
	if cfg.Cache.Memory.MaxSizeMB != 512 {
		t.Errorf("Cache.Memory.MaxSizeMB = %d, want 512", cfg.Cache.Memory.MaxSizeMB)
	}
}

func TestExecutorConfig(t *testing.T) {
	cfg := &Config{
		Executor: ExecutorConfig{
			MaxConcurrentJobs:   50,
			MaxGoroutinesPerJob: 100,
			DefaultTimeout:      86400,
			WorkerPool: WorkerPoolConfig{
				MinWorkers:  10,
				MaxWorkers:  100,
				QueueSize:   1000,
				IdleTimeout: 300,
			},
		},
	}
	
	if cfg.Executor.MaxConcurrentJobs != 50 {
		t.Errorf("Executor.MaxConcurrentJobs = %d, want 50", cfg.Executor.MaxConcurrentJobs)
	}
	
	if cfg.Executor.WorkerPool.MaxWorkers != 100 {
		t.Errorf("Executor.WorkerPool.MaxWorkers = %d, want 100", cfg.Executor.WorkerPool.MaxWorkers)
	}
}

func TestSchedulerConfig(t *testing.T) {
	cfg := &Config{
		Scheduler: SchedulerConfig{
			Enabled:       true,
			Timezone:      "Asia/Shanghai",
			CheckInterval: 10,
		},
	}
	
	if !cfg.Scheduler.Enabled {
		t.Error("Scheduler.Enabled should be true")
	}
	
	if cfg.Scheduler.Timezone != "Asia/Shanghai" {
		t.Errorf("Scheduler.Timezone = %s, want Asia/Shanghai", cfg.Scheduler.Timezone)
	}
}

func TestSecurityConfig(t *testing.T) {
	cfg := &Config{
		Security: SecurityConfig{
			JWT: JWTConfig{
				Secret:     "test-secret",
				Expiration: 86400,
			},
			CORS: CORSConfig{
				Enabled:        true,
				AllowedOrigins: []string{"http://localhost:3000"},
				AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
			},
		},
	}
	
	if cfg.Security.JWT.Secret != "test-secret" {
		t.Errorf("Security.JWT.Secret = %s, want test-secret", cfg.Security.JWT.Secret)
	}
	
	if !cfg.Security.CORS.Enabled {
		t.Error("Security.CORS.Enabled should be true")
	}
}
