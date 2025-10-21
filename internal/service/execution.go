package service

import (
	"github.com/fustgo/fustgo2/internal/core/types"
	"github.com/fustgo/fustgo2/internal/repository"
	"github.com/google/uuid"
	"time"
)

// ExecutionService 执行记录服务
type ExecutionService struct {
	executionRepo *repository.ExecutionRepository
}

// NewExecutionService 创建执行记录服务实例
func NewExecutionService(executionRepo *repository.ExecutionRepository) *ExecutionService {
	return &ExecutionService{executionRepo: executionRepo}
}

// Create 创建执行记录
func (s *ExecutionService) Create(req *CreateExecutionRequest) (*types.Execution, error) {
	execution := &types.Execution{
		ID:          uuid.New().String(),
		JobID:       req.JobID,
		Status:      types.ExecutionStatusPending,
		StartTime:   time.Now(),
		TriggeredBy: req.TriggeredBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.executionRepo.Create(execution); err != nil {
		return nil, err
	}

	return execution, nil
}

// GetByID 根据ID获取执行记录
func (s *ExecutionService) GetByID(id string) (*types.Execution, error) {
	return s.executionRepo.GetByID(id)
}

// ListByJobID 根据任务ID获取执行记录列表
func (s *ExecutionService) ListByJobID(jobID string, page, pageSize int) ([]*types.Execution, int64, error) {
	total, err := s.executionRepo.CountByJobID(jobID)
	if err != nil {
		return nil, 0, err
	}

	executions, err := s.executionRepo.ListByJobID(jobID, (page-1)*pageSize, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return executions, total, nil
}

// List 获取执行记录列表
func (s *ExecutionService) List(page, pageSize int) ([]*types.Execution, int64, error) {
	total, err := s.executionRepo.Count()
	if err != nil {
		return nil, 0, err
	}

	executions, err := s.executionRepo.List((page-1)*pageSize, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return executions, total, nil
}

// Update 更新执行记录
func (s *ExecutionService) Update(id string, req *UpdateExecutionRequest) error {
	execution, err := s.executionRepo.GetByID(id)
	if err != nil {
		return err
	}

	execution.Status = req.Status
	execution.EndTime = req.EndTime
	execution.Duration = req.Duration
	execution.RecordsRead = req.RecordsRead
	execution.RecordsWritten = req.RecordsWritten
	execution.RecordsFiltered = req.RecordsFiltered
	execution.RecordsError = req.RecordsError
	execution.BytesRead = req.BytesRead
	execution.BytesWritten = req.BytesWritten
	execution.ErrorMessage = req.ErrorMessage
	execution.ErrorStack = req.ErrorStack
	execution.Context = req.Context
	execution.Metrics = req.Metrics
	execution.UpdatedAt = time.Now()

	return s.executionRepo.Update(execution)
}

// CreateExecutionRequest 创建执行记录请求
type CreateExecutionRequest struct {
	JobID       string `json:"job_id" validate:"required"`
	TriggeredBy string `json:"triggered_by"`
}

// UpdateExecutionRequest 更新执行记录请求
type UpdateExecutionRequest struct {
	Status         types.ExecutionStatus `json:"status"`
	EndTime        *time.Time            `json:"end_time"`
	Duration       int64                 `json:"duration"`
	RecordsRead    int64                 `json:"records_read"`
	RecordsWritten int64                 `json:"records_written"`
	RecordsFiltered int64                `json:"records_filtered"`
	RecordsError   int64                 `json:"records_error"`
	BytesRead      int64                 `json:"bytes_read"`
	BytesWritten   int64                 `json:"bytes_written"`
	ErrorMessage   string                `json:"error_message"`
	ErrorStack     string                `json:"error_stack"`
	Context        string                `json:"context"`
	Metrics        string                `json:"metrics"`
}