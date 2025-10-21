package repository

import (
	"github.com/fustgo/fustgo2/internal/core/types"
	"gorm.io/gorm"
)

// PipelineRepository 管道仓库
type PipelineRepository struct {
	db *gorm.DB
}

// NewPipelineRepository 创建管道仓库实例
func NewPipelineRepository(db *gorm.DB) *PipelineRepository {
	return &PipelineRepository{db: db}
}

// Create 创建管道
func (r *PipelineRepository) Create(pipeline *types.Pipeline) error {
	return r.db.Create(pipeline).Error
}

// GetByID 根据ID获取管道
func (r *PipelineRepository) GetByID(id string) (*types.Pipeline, error) {
	var pipeline types.Pipeline
	if err := r.db.Where("id = ?", id).First(&pipeline).Error; err != nil {
		return nil, err
	}
	return &pipeline, nil
}

// GetByName 根据名称获取管道
func (r *PipelineRepository) GetByName(name string) (*types.Pipeline, error) {
	var pipeline types.Pipeline
	if err := r.db.Where("name = ?", name).First(&pipeline).Error; err != nil {
		return nil, err
	}
	return &pipeline, nil
}

// List 获取管道列表
func (r *PipelineRepository) List(offset, limit int) ([]*types.Pipeline, int64, error) {
	var pipelines []*types.Pipeline
	var total int64

	if err := r.db.Model(&types.Pipeline{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&pipelines).Error; err != nil {
		return nil, 0, err
	}

	return pipelines, total, nil
}

// Update 更新管道
func (r *PipelineRepository) Update(pipeline *types.Pipeline) error {
	return r.db.Save(pipeline).Error
}

// Delete 删除管道
func (r *PipelineRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&types.Pipeline{}).Error
}