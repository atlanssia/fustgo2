package pipeline

import (
	"fmt"
	"math/rand"
	"time"
)

// ErrorHandler 错误处理器
type ErrorHandler struct {
	config *Config
}

// NewErrorHandler 创建新的错误处理器
func NewErrorHandler(config *Config) *ErrorHandler {
	return &ErrorHandler{
		config: config,
	}
}

// HandleError 处理错误
func (h *ErrorHandler) HandleError(err error) error {
	// 根据配置的错误策略处理
	switch h.config.ErrorStrategy {
	case StrategyFailFast:
		return fmt.Errorf("fail fast: %w", err)
	case StrategyContinue:
		// 记录错误并继续
		return nil
	case StrategyRetry:
		// 重试逻辑在执行器中处理
		return err
	default:
		return fmt.Errorf("unknown error strategy: %v, error: %w", h.config.ErrorStrategy, err)
	}
}

// ShouldRetry 判断是否应该重试
func (h *ErrorHandler) ShouldRetry(err error, attempt int) bool {
	if h.config.ErrorStrategy != StrategyRetry {
		return false
	}
	
	if attempt >= h.config.MaxRetries {
		return false
	}
	
	// 检查错误是否可重试
	if retryableErr, ok := err.(interface{ IsRetryable() bool }); ok {
		return retryableErr.IsRetryable()
	}
	
	// 默认认为所有错误都可重试
	return true
}

// CalculateRetryDelay 计算重试延迟
func (h *ErrorHandler) CalculateRetryDelay(attempt int) time.Duration {
	// 基础延迟
	delay := h.config.RetryInterval
	
	// 指数退避
	for i := 0; i < attempt; i++ {
		delay *= 2
	}
	
	// 添加随机抖动 (±25%)
	jitter := time.Duration(rand.Int63n(int64(delay/2)) - int64(delay/4))
	delay += jitter
	
	// 限制最大延迟
	if delay > time.Minute {
		delay = time.Minute
	}
	
	return delay
}

// RetryableError 可重试错误
type RetryableError struct {
	Err       error
	Retryable bool
}

// NewRetryableError 创建可重试错误
func NewRetryableError(err error, retryable bool) *RetryableError {
	return &RetryableError{
		Err:       err,
		Retryable: retryable,
	}
}

// Error 实现 error 接口
func (e *RetryableError) Error() string {
	return e.Err.Error()
}

// Unwrap 实现 unwrap 接口
func (e *RetryableError) Unwrap() error {
	return e.Err
}

// IsRetryable 判断是否可重试
func (e *RetryableError) IsRetryable() bool {
	return e.Retryable
}
