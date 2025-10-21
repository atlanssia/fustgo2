package pipeline

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/fustgo/fustgo2/pkg/plugin"
	"github.com/fustgo/fustgo2/pkg/record"
)

// Executor 管道执行器
type Executor struct {
	pipeline *Pipeline
	
	// 控制
	mu     sync.RWMutex
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
	
	// 状态
	started bool
	stopped bool
	paused  bool
	
	// 错误处理
	errorHandler *ErrorHandler
	
	// 状态管理
	stateManager *StateManager
}

// NewExecutor 创建新的执行器
func NewExecutor(pipeline *Pipeline) *Executor {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &Executor{
		pipeline:     pipeline,
		ctx:          ctx,
		cancel:       cancel,
		errorHandler: NewErrorHandler(pipeline.Config),
		stateManager: NewStateManager(pipeline.Config),
	}
}

// Start 启动管道执行
func (e *Executor) Start() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	if e.started {
		return fmt.Errorf("pipeline already started")
	}
	
	if e.stopped {
		return fmt.Errorf("pipeline already stopped, cannot restart")
	}
	
	// 验证管道
	if err := e.pipeline.Validate(); err != nil {
		return fmt.Errorf("pipeline validation failed: %w", err)
	}
	
	// 更新状态
	e.pipeline.State.SetStatus(StatusRunning)
	e.pipeline.State.SetStartTime(time.Now())
	e.started = true
	
	// 启动执行
	go e.execute()
	
	return nil
}

// Stop 停止管道执行
func (e *Executor) Stop() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	if !e.started {
		return fmt.Errorf("pipeline not started")
	}
	
	if e.stopped {
		return fmt.Errorf("pipeline already stopped")
	}
	
	// 取消上下文
	e.cancel()
	e.stopped = true
	e.pipeline.State.SetStatus(StatusStopped)
	e.pipeline.State.SetEndTime(time.Now())
	
	// 等待所有协程完成
	e.wg.Wait()
	
	return nil
}

// Pause 暂停管道执行
func (e *Executor) Pause() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	if !e.started {
		return fmt.Errorf("pipeline not started")
	}
	
	if e.stopped {
		return fmt.Errorf("pipeline already stopped")
	}
	
	if e.paused {
		return fmt.Errorf("pipeline already paused")
	}
	
	e.paused = true
	e.pipeline.State.SetStatus(StatusPaused)
	
	return nil
}

// Resume 恢复管道执行
func (e *Executor) Resume() error {
	e.mu.Lock()
	defer e.mu.Unlock()
	
	if !e.started {
		return fmt.Errorf("pipeline not started")
	}
	
	if e.stopped {
		return fmt.Errorf("pipeline already stopped")
	}
	
	if !e.paused {
		return fmt.Errorf("pipeline not paused")
	}
	
	e.paused = false
	e.pipeline.State.SetStatus(StatusRunning)
	
	return nil
}

// Wait 等待管道执行完成
func (e *Executor) Wait() error {
	// 等待所有协程完成
	e.wg.Wait()
	return nil
}

// IsRunning 检查管道是否正在运行
func (e *Executor) IsRunning() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.started && !e.stopped && !e.paused
}

// execute 执行管道
func (e *Executor) execute() {
	// 创建通道
	// 我们需要为每个连接创建通道：Reader -> Processor1 -> Processor2 -> ... -> Writer
	channelCount := len(e.pipeline.Processors) + 1
	channels := make([]chan *record.Record, channelCount)
	for i := range channels {
		channels[i] = make(chan *record.Record, e.pipeline.Config.ChannelSize)
	}
	e.pipeline.channels = channels
	
	// 启动 Writer
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		e.executeWriter()
	}()
	
	// 启动 Processors (逆序启动以确保依赖正确)
	for i := len(e.pipeline.Processors) - 1; i >= 0; i-- {
		processor := e.pipeline.Processors[i]
		e.wg.Add(1)
		go func(idx int, p plugin.Processor) {
			defer e.wg.Done()
			e.executeProcessor(idx, p)
		}(i, processor)
	}
	
	// 启动 Reader
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		e.executeReader()
	}()
	
	// 等待所有协程完成
	e.wg.Wait()
	
	// 确保在函数退出时更新状态
	e.mu.Lock()
	if e.pipeline.State.GetStatus() == StatusRunning {
		e.pipeline.State.SetStatus(StatusComplete)
		e.pipeline.State.SetEndTime(time.Now())
	}
	e.mu.Unlock()
	
	// 取消上下文
	e.cancel()
}

// executeReader 执行 Reader
func (e *Executor) executeReader() {
	err := e.pipeline.Reader.Read(e.ctx, e.pipeline.channels[0])
	if err != nil {
		e.handleError(fmt.Errorf("reader error: %w", err))
		// Reader 应该自己关闭通道，但如果没关闭，我们需要确保关闭
		select {
		case <-e.ctx.Done():
			// 上下文已取消，不需要处理
		default:
			// 尝试关闭通道（如果还未关闭）
			close(e.pipeline.channels[0])
		}
		return
	}
	
	// 更新指标
	progress := e.pipeline.Reader.GetProgress()
	e.pipeline.Metrics.IncrRecordsRead(progress.ProcessedRecords)
}

// executeProcessor 执行 Processor
func (e *Executor) executeProcessor(index int, processor plugin.Processor) {
	input := e.pipeline.channels[index]
	output := e.pipeline.channels[index+1]
	
	err := processor.Process(e.ctx, input, output)
	if err != nil {
		e.handleError(fmt.Errorf("processor %d error: %w", index, err))
		return
	}
	
	// 更新指标
	metrics := processor.GetMetrics()
	e.pipeline.Metrics.IncrRecordsProcessed(metrics.RecordsProcessed)
	e.pipeline.Metrics.IncrRecordsFiltered(metrics.FilteredCount)
	e.pipeline.Metrics.IncrRecordsError(metrics.ErrorCount)
}

// executeWriter 执行 Writer
func (e *Executor) executeWriter() {
	// Writer 的输入通道取决于是否有 Processor
	var input <-chan *record.Record
	if len(e.pipeline.Processors) > 0 {
		// 有 Processor，使用最后一个 Processor 的输出通道
		input = e.pipeline.channels[len(e.pipeline.channels)-1]
	} else {
		// 没有 Processor，直接连接 Reader
		input = e.pipeline.channels[0]
	}
	
	err := e.pipeline.Writer.Write(e.ctx, input)
	if err != nil {
		e.handleError(fmt.Errorf("writer error: %w", err))
		return
	}
	
	// 刷新 Writer
	if err := e.pipeline.Writer.Flush(); err != nil {
		e.handleError(fmt.Errorf("writer flush error: %w", err))
		return
	}
	
	// 更新指标
	metrics := e.pipeline.Writer.GetMetrics()
	e.pipeline.Metrics.IncrRecordsWritten(metrics.RecordsProcessed)
	e.pipeline.Metrics.IncrRecordsError(metrics.ErrorCount)
}

// handleError 处理错误
func (e *Executor) handleError(err error) {
	e.pipeline.State.SetError(err)
	
	// 根据错误策略处理
	switch e.pipeline.Config.ErrorStrategy {
	case StrategyFailFast:
		// 快速失败，停止管道
		e.Stop()
	case StrategyContinue:
		// 继续执行，记录错误
		e.pipeline.Metrics.IncrRecordsError(1)
	case StrategyRetry:
		// 重试逻辑（简化处理）
		e.pipeline.Metrics.IncrRecordsError(1)
	}
}
