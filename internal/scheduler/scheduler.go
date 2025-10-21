package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/fustgo/fustgo2/internal/core/types"
	"github.com/fustgo/fustgo2/internal/repository"
	"github.com/fustgo/fustgo2/internal/service"
	"github.com/fustgo/fustgo2/pkg/pipeline"
	"github.com/fustgo/fustgo2/pkg/plugin"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// Scheduler 任务调度器
type Scheduler struct {
	jobRepo       *repository.JobRepository
	jobService    *service.JobService
	pipelineRepo  *repository.PipelineRepository
	executor      *pipeline.Executor
	logger        *zap.Logger
	cron          *cron.Cron
	entries       map[string]cron.EntryID
	mu            sync.RWMutex
	ctx           context.Context
	cancel        context.CancelFunc
}

// NewScheduler 创建调度器实例
func NewScheduler(
	jobRepo *repository.JobRepository,
	jobService *service.JobService,
	pipelineRepo *repository.PipelineRepository,
	logger *zap.Logger,
) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	
	// 创建cron调度器
	c := cron.New(cron.WithSeconds())

	return &Scheduler{
		jobRepo:      jobRepo,
		jobService:   jobService,
		pipelineRepo: pipelineRepo,
		logger:       logger,
		cron:         c,
		entries:      make(map[string]cron.EntryID),
		ctx:          ctx,
		cancel:       cancel,
	}
}

// Start 启动调度器
func (s *Scheduler) Start() error {
	s.logger.Info("Starting scheduler")

	// 启动cron调度器
	s.cron.Start()

	// 加载所有启用的任务
	jobs, err := s.loadEnabledJobs()
	if err != nil {
		return fmt.Errorf("failed to load enabled jobs: %w", err)
	}

	// 注册任务到调度器
	for _, job := range jobs {
		if err := s.scheduleJob(job); err != nil {
			s.logger.Error("Failed to schedule job", 
				zap.String("job_id", job.ID), 
				zap.String("job_name", job.Name), 
				zap.Error(err))
			continue
		}
	}

	s.logger.Info("Scheduler started", zap.Int("scheduled_jobs", len(jobs)))
	return nil
}

// Stop 停止调度器
func (s *Scheduler) Stop() error {
	s.logger.Info("Stopping scheduler")

	// 取消上下文
	s.cancel()

	// 停止cron调度器
	s.cron.Stop()

	s.logger.Info("Scheduler stopped")
	return nil
}

// loadEnabledJobs 加载所有启用的任务
func (s *Scheduler) loadEnabledJobs() ([]*types.Job, error) {
	// 查询所有启用的任务
	// 这里简化处理，实际应该分页查询
	jobs, err := s.jobRepo.List(0, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to list jobs: %w", err)
	}
	
	// 过滤启用的任务
	var enabledJobs []*types.Job
	for _, job := range jobs {
		if job.Enabled {
			enabledJobs = append(enabledJobs, job)
		}
	}
	
	return enabledJobs, nil
}

// scheduleJob 注册任务到调度器
func (s *Scheduler) scheduleJob(job *types.Job) error {
	// 解析调度配置
	scheduleConfig := &types.ScheduleConfig{}
	if err := parseJSON(job.ScheduleConfig, scheduleConfig); err != nil {
		return fmt.Errorf("failed to parse schedule config: %w", err)
	}

	// 根据调度类型注册任务
	switch scheduleConfig.Type {
	case types.ScheduleTypeCron:
		return s.scheduleCronJob(job, scheduleConfig)
	case types.ScheduleTypeManual:
		// 手动触发任务，不需要注册到调度器
		return nil
	case types.ScheduleTypeEvent:
		// 事件触发任务，需要注册事件监听器
		return s.scheduleEventJob(job, scheduleConfig)
	case types.ScheduleTypeDependency:
		// 依赖调度任务，需要注册依赖监听器
		return s.scheduleDependencyJob(job, scheduleConfig)
	default:
		return fmt.Errorf("unsupported schedule type: %s", scheduleConfig.Type)
	}
}

// scheduleCronJob 注册Cron任务
func (s *Scheduler) scheduleCronJob(job *types.Job, scheduleConfig *types.ScheduleConfig) error {
	if scheduleConfig.Cron == "" {
		return fmt.Errorf("cron expression is required")
	}

	// 添加到cron调度器
	entryID, err := s.cron.AddFunc(scheduleConfig.Cron, func() {
		s.executeJob(job)
	})
	if err != nil {
		return fmt.Errorf("failed to add cron job: %w", err)
	}

	// 记录entry ID
	s.mu.Lock()
	s.entries[job.ID] = entryID
	s.mu.Unlock()

	s.logger.Info("Scheduled cron job", 
		zap.String("job_id", job.ID), 
		zap.String("job_name", job.Name), 
		zap.String("cron", scheduleConfig.Cron))

	return nil
}

// scheduleEventJob 注册事件触发任务
func (s *Scheduler) scheduleEventJob(job *types.Job, scheduleConfig *types.ScheduleConfig) error {
	// TODO: 实现事件监听器注册
	s.logger.Info("Scheduled event job", 
		zap.String("job_id", job.ID), 
		zap.String("job_name", job.Name))
	return nil
}

// scheduleDependencyJob 注册依赖调度任务
func (s *Scheduler) scheduleDependencyJob(job *types.Job, scheduleConfig *types.ScheduleConfig) error {
	// TODO: 实现依赖监听器注册
	s.logger.Info("Scheduled dependency job", 
		zap.String("job_id", job.ID), 
		zap.String("job_name", job.Name))
	return nil
}

// executeJob 执行任务
func (s *Scheduler) executeJob(job *types.Job) {
	s.logger.Info("Executing job", 
		zap.String("job_id", job.ID), 
		zap.String("job_name", job.Name))

	// 更新任务状态
	job.Status = types.JobStatusRunning
	job.LastExecutionTime = &time.Time{}
	*job.LastExecutionTime = time.Now()
	
	// TODO: 实际执行任务逻辑
	// 这里应该解析管道配置并执行管道

	// 模拟执行
	time.Sleep(1 * time.Second)

	// 更新任务状态
	job.Status = types.JobStatusSucceeded
	job.LastExecutionTime = &time.Time{}
	*job.LastExecutionTime = time.Now()

	s.logger.Info("Job executed successfully", 
		zap.String("job_id", job.ID), 
		zap.String("job_name", job.Name))
}

// AddJob 添加任务到调度器
func (s *Scheduler) AddJob(job *types.Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 如果任务已存在，先移除
	if entryID, exists := s.entries[job.ID]; exists {
		s.cron.Remove(entryID)
		delete(s.entries, job.ID)
	}

	// 注册新任务
	return s.scheduleJob(job)
}

// RemoveJob 从调度器移除任务
func (s *Scheduler) RemoveJob(jobID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entryID, exists := s.entries[jobID]; exists {
		s.cron.Remove(entryID)
		delete(s.entries, jobID)
		s.logger.Info("Removed job from scheduler", zap.String("job_id", jobID))
	}
}

// parseJSON 解析JSON字符串到结构体
func parseJSON(jsonStr string, v interface{}) error {
	return json.Unmarshal([]byte(jsonStr), v)
}