package service

import (
	"github.com/fustgo/fustgo2/internal/core/types"
	"github.com/fustgo/fustgo2/internal/repository"
	"github.com/google/uuid"
	"time"
)

// JobService 任务服务
type JobService struct {
	jobRepo *repository.JobRepository
}

// NewJobService 创建任务服务实例
func NewJobService(jobRepo *repository.JobRepository) *JobService {
	return &JobService{jobRepo: jobRepo}
}

// Create 创建任务
func (s *JobService) Create(req *CreateJobRequest) (*types.Job, error) {
	job := &types.Job{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Enabled:     req.Enabled,
		Status:      types.JobStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreatedBy:   req.CreatedBy,
	}

	if err := s.jobRepo.Create(job); err != nil {
		return nil, err
	}

	return job, nil
}

// GetByID 根据ID获取任务
func (s *JobService) GetByID(id string) (*types.Job, error) {
	return s.jobRepo.GetByID(id)
}

// List 获取任务列表
func (s *JobService) List(page, pageSize int) ([]*types.Job, int64, error) {
	total, err := s.jobRepo.Count()
	if err != nil {
		return nil, 0, err
	}

	jobs, err := s.jobRepo.List((page-1)*pageSize, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return jobs, total, nil
}

// Update 更新任务
func (s *JobService) Update(id string, req *UpdateJobRequest) error {
	job, err := s.jobRepo.GetByID(id)
	if err != nil {
		return err
	}

	job.Name = req.Name
	job.Description = req.Description
	job.Enabled = req.Enabled
	job.UpdatedAt = time.Now()

	return s.jobRepo.Update(job)
}

// Delete 删除任务
func (s *JobService) Delete(id string) error {
	return s.jobRepo.Delete(id)
}

// CreateJobRequest 创建任务请求
type CreateJobRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	CreatedBy   string `json:"created_by"`
}

// UpdateJobRequest 更新任务请求
type UpdateJobRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}