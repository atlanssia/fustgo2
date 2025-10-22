package repository

import (
	"github.com/fustgo/fustgo2/internal/core/types"
	"gorm.io/gorm"
)

// JobRepository 任务仓库
type JobRepository struct {
	db *gorm.DB
}

// NewJobRepository 创建任务仓库实例
func NewJobRepository(db *gorm.DB) *JobRepository {
	return &JobRepository{db: db}
}

// Create 创建任务
func (r *JobRepository) Create(job *types.Job) error {
	return r.db.Create(job).Error
}

// GetByID 根据ID获取任务
func (r *JobRepository) GetByID(id string) (*types.Job, error) {
	var job types.Job
	if err := r.db.Where("id = ?", id).First(&job).Error; err != nil {
		return nil, err
	}
	return &job, nil
}

// List 获取任务列表
func (r *JobRepository) List(offset, limit int) ([]*types.Job, error) {
	var jobs []*types.Job
	if err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

// Update 更新任务
func (r *JobRepository) Update(job *types.Job) error {
	return r.db.Save(job).Error
}

// Delete 删除任务
func (r *JobRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&types.Job{}).Error
}

// Count 获取任务总数
func (r *JobRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&types.Job{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}