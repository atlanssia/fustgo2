package optimizer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/fustgo/fustgo2/internal/core/types"
	"github.com/fustgo/fustgo2/internal/repository"
	"github.com/fustgo/fustgo2/pkg/pipeline"
	"go.uber.org/zap"
)

// Optimizer 性能优化器
type Optimizer struct {
	jobRepo       *repository.JobRepository
	executionRepo *repository.ExecutionRepository
	pipelineRepo  *repository.PipelineRepository
	logger        *zap.Logger
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	interval      time.Duration
}

// NewOptimizer 创建性能优化器实例
func NewOptimizer(
	jobRepo *repository.JobRepository,
	executionRepo *repository.ExecutionRepository,
	pipelineRepo *repository.PipelineRepository,
	logger *zap.Logger,
) *Optimizer {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &Optimizer{
		jobRepo:       jobRepo,
		executionRepo: executionRepo,
		pipelineRepo:  pipelineRepo,
		logger:        logger,
		ctx:           ctx,
		cancel:        cancel,
		interval:      1 * time.Hour, // 默认1小时优化一次
	}
}

// Start 启动优化器
func (o *Optimizer) Start() error {
	o.logger.Info("Starting optimizer")

	// 启动优化goroutine
	o.wg.Add(1)
	go o.optimizeLoop()

	o.logger.Info("Optimizer started")
	return nil
}

// Stop 停止优化器
func (o *Optimizer) Stop() error {
	o.logger.Info("Stopping optimizer")

	// 取消上下文
	o.cancel()

	// 等待所有goroutine完成
	o.wg.Wait()

	o.logger.Info("Optimizer stopped")
	return nil
}

// optimizeLoop 优化循环
func (o *Optimizer) optimizeLoop() {
	defer o.wg.Done()

	ticker := time.NewTicker(o.interval)
	defer ticker.Stop()

	// 立即执行一次优化
	if err := o.optimize(); err != nil {
		o.logger.Error("Failed to perform initial optimization", zap.Error(err))
	}

	for {
		select {
		case <-o.ctx.Done():
			return
		case <-ticker.C:
			if err := o.optimize(); err != nil {
				o.logger.Error("Failed to perform optimization", zap.Error(err))
			}
		}
	}
}

// optimize 执行优化
func (o *Optimizer) optimize() error {
	o.logger.Info("Performing optimization")

	// 优化任务配置
	if err := o.optimizeJobs(); err != nil {
		return fmt.Errorf("failed to optimize jobs: %w", err)
	}

	// 优化执行记录
	if err := o.optimizeExecutions(); err != nil {
		return fmt.Errorf("failed to optimize executions: %w", err)
	}

	// 优化管道配置
	if err := o.optimizePipelines(); err != nil {
		return fmt.Errorf("failed to optimize pipelines: %w", err)
	}

	o.logger.Info("Optimization completed")
	return nil
}

// optimizeJobs 优化任务配置
func (o *Optimizer) optimizeJobs() error {
	o.logger.Debug("Optimizing jobs")

	// 获取所有任务
	jobs, err := o.jobRepo.List(0, 1000)
	if err != nil {
		return fmt.Errorf("failed to list jobs: %w", err)
	}

	// 分析任务执行历史，优化配置
	for _, job := range jobs {
		if err := o.optimizeJob(job); err != nil {
			o.logger.Error("Failed to optimize job", 
				zap.String("job_id", job.ID), 
				zap.String("job_name", job.Name), 
				zap.Error(err))
			continue
		}
	}

	return nil
}

// optimizeJob 优化单个任务
func (o *Optimizer) optimizeJob(job *types.Job) error {
	// 获取任务的执行历史
	executions, err := o.executionRepo.ListByJobID(job.ID, 0, 100)
	if err != nil {
		return fmt.Errorf("failed to list executions for job %s: %w", job.ID, err)
	}

	if len(executions) == 0 {
		// 没有执行历史，跳过优化
		return nil
	}

	// 分析执行历史
	stats := o.analyzeExecutionHistory(executions)

	// 根据分析结果优化配置
	// 这里简化处理，实际应该根据具体的性能指标进行优化
	o.logger.Debug("Job execution statistics",
		zap.String("job_id", job.ID),
		zap.String("job_name", job.Name),
		zap.Any("stats", stats))

	return nil
}

// analyzeExecutionHistory 分析执行历史
func (o *Optimizer) analyzeExecutionHistory(executions []*types.Execution) *ExecutionStats {
	stats := &ExecutionStats{
		TotalExecutions: len(executions),
	}

	if len(executions) == 0 {
		return stats
	}

	// 计算平均执行时间
	var totalDuration int64
	for _, exec := range executions {
		totalDuration += exec.Duration
	}
	stats.AvgDuration = totalDuration / int64(len(executions))

	// 计算成功率
	var successCount int64
	for _, exec := range executions {
		if exec.Status == types.ExecutionStatusSucceeded {
			successCount++
		}
	}
	stats.SuccessRate = float64(successCount) / float64(len(executions))

	// 计算平均记录数
	var totalRecords int64
	for _, exec := range executions {
		totalRecords += exec.RecordsRead
	}
	stats.AvgRecords = totalRecords / int64(len(executions))

	return stats
}

// optimizeExecutions 优化执行记录
func (o *Optimizer) optimizeExecutions() error {
	o.logger.Debug("Optimizing executions")

	// 清理过期的执行记录
	// 这里简化处理，实际应该根据配置的保留策略进行清理

	return nil
}

// optimizePipelines 优化管道配置
func (o *Optimizer) optimizePipelines() error {
	o.logger.Debug("Optimizing pipelines")

	// 获取所有管道
	pipelines, err := o.pipelineRepo.List(0, 1000)
	if err != nil {
		return fmt.Errorf("failed to list pipelines: %w", err)
	}

	// 分析管道执行历史，优化配置
	for _, pipeline := range pipelines {
		if err := o.optimizePipeline(pipeline); err != nil {
			o.logger.Error("Failed to optimize pipeline", 
				zap.String("pipeline_id", pipeline.ID), 
				zap.String("pipeline_name", pipeline.Name), 
				zap.Error(err))
			continue
		}
	}

	return nil
}

// optimizePipeline 优化单个管道
func (o *Optimizer) optimizePipeline(p *types.Pipeline) error {
	// 获取管道的执行历史
	// 这里简化处理，实际应该分析管道的性能指标并优化配置

	o.logger.Debug("Optimizing pipeline",
		zap.String("pipeline_id", p.ID),
		zap.String("pipeline_name", p.Name))

	return nil
}

// SetInterval 设置优化间隔
func (o *Optimizer) SetInterval(interval time.Duration) {
	o.interval = interval
}

// ExecutionStats 执行统计信息
type ExecutionStats struct {
	TotalExecutions int     `json:"total_executions"`
	AvgDuration     int64   `json:"avg_duration"`
	SuccessRate     float64 `json:"success_rate"`
	AvgRecords      int64   `json:"avg_records"`
}