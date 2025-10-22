package service

import (
	"encoding/json"
	"github.com/fustgo/fustgo2/internal/core/types"
	"github.com/fustgo/fustgo2/internal/repository"
	"github.com/google/uuid"
	"time"
)

// PipelineService 管道服务
type PipelineService struct {
	pipelineRepo *repository.PipelineRepository
}

// NewPipelineService 创建管道服务实例
func NewPipelineService(pipelineRepo *repository.PipelineRepository) *PipelineService {
	return &PipelineService{pipelineRepo: pipelineRepo}
}

// CreatePipeline 创建管道
func (s *PipelineService) CreatePipeline(req *CreatePipelineRequest) (*types.Pipeline, error) {
	// 序列化配置
	configBytes, err := json.Marshal(req.Config)
	if err != nil {
		return nil, err
	}

	// 序列化调度配置
	scheduleBytes, err := json.Marshal(req.Schedule)
	if err != nil {
		return nil, err
	}

	pipeline := &types.Pipeline{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Status:      types.PipelineStatusDraft,
		Config:      string(configBytes),
		Schedule:    string(scheduleBytes),
		Tags:        req.Tags,
		Enabled:     req.Enabled,
		CreatedBy:   req.CreatedBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.pipelineRepo.Create(pipeline); err != nil {
		return nil, err
	}

	return pipeline, nil
}

// GetPipelineByID 根据ID获取管道
func (s *PipelineService) GetPipelineByID(id string) (*types.Pipeline, error) {
	return s.pipelineRepo.GetByID(id)
}

// ListPipelines 获取管道列表
func (s *PipelineService) ListPipelines(page, pageSize int) ([]*types.Pipeline, int64, error) {
	return s.pipelineRepo.List((page-1)*pageSize, pageSize)
}

// UpdatePipeline 更新管道
func (s *PipelineService) UpdatePipeline(id string, req *UpdatePipelineRequest) error {
	pipeline, err := s.pipelineRepo.GetByID(id)
	if err != nil {
		return err
	}

	pipeline.Name = req.Name
	pipeline.Description = req.Description
	pipeline.Tags = req.Tags
	pipeline.Enabled = req.Enabled
	pipeline.UpdatedAt = time.Now()

	return s.pipelineRepo.Update(pipeline)
}

// DeletePipeline 删除管道
func (s *PipelineService) DeletePipeline(id string) error {
	return s.pipelineRepo.Delete(id)
}

// UpdatePipelineStatus 更新管道状态
func (s *PipelineService) UpdatePipelineStatus(id string, status types.PipelineStatus) error {
	pipeline, err := s.pipelineRepo.GetByID(id)
	if err != nil {
		return err
	}

	pipeline.Status = status
	pipeline.UpdatedAt = time.Now()

	return s.pipelineRepo.Update(pipeline)
}

// CreatePipelineRequest 创建管道请求
type CreatePipelineRequest struct {
	Name        string                   `json:"name" validate:"required"`
	Description string                   `json:"description"`
	Config      types.PipelineConfig     `json:"config" validate:"required"`
	Schedule    interface{}              `json:"schedule"`
	Tags        string                   `json:"tags"`
	Enabled     bool                     `json:"enabled"`
	CreatedBy   string                   `json:"created_by"`
}

// UpdatePipelineRequest 更新管道请求
type UpdatePipelineRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Tags        string `json:"tags"`
	Enabled     bool   `json:"enabled"`
}