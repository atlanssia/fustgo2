package repository

import (
	"github.com/fustgo/fustgo2/internal/core/types"
	"gorm.io/gorm"
)

// ExecutionRepository 执行记录仓库
type ExecutionRepository struct {
	db *gorm.DB
}

// NewExecutionRepository 创建执行记录仓库实例
func NewExecutionRepository(db *gorm.DB) *ExecutionRepository {
	return &ExecutionRepository{db: db}
}

// Create 创建执行记录
func (r *ExecutionRepository) Create(execution *types.Execution) error {
	return r.db.Create(execution).Error
}

// GetByID 根据ID获取执行记录
func (r *ExecutionRepository) GetByID(id string) (*types.Execution, error) {
	var execution types.Execution
	if err := r.db.Where("id = ?", id).First(&execution).Error; err != nil {
		return nil, err
	}
	return &execution, nil
}

// ListByJobID 根据任务ID获取执行记录列表
func (r *ExecutionRepository) ListByJobID(jobID string, offset, limit int) ([]*types.Execution, error) {
	var executions []*types.Execution
	if err := r.db.Where("job_id = ?", jobID).Offset(offset).Limit(limit).Order("created_at DESC").Find(&executions).Error; err != nil {
		return nil, err
	}
	return executions, nil
}

// List 获取执行记录列表
func (r *ExecutionRepository) List(offset, limit int) ([]*types.Execution, error) {
	var executions []*types.Execution
	if err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&executions).Error; err != nil {
		return nil, err
	}
	return executions, nil
}

// Update 更新执行记录
func (r *ExecutionRepository) Update(execution *types.Execution) error {
	return r.db.Save(execution).Error
}

// Delete 删除执行记录
func (r *ExecutionRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&types.Execution{}).Error
}

// Count 获取执行记录总数
func (r *ExecutionRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&types.Execution{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// CountByJobID 根据任务ID获取执行记录总数
func (r *ExecutionRepository) CountByJobID(jobID string) (int64, error) {
	var count int64
	if err := r.db.Model(&types.Execution{}).Where("job_id = ?", jobID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}