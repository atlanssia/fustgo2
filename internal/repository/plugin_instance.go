package repository

import (
	"github.com/fustgo/fustgo2/internal/core/types"
	"gorm.io/gorm"
)

// PluginInstanceRepository 插件实例仓库
type PluginInstanceRepository struct {
	db *gorm.DB
}

// NewPluginInstanceRepository 创建插件实例仓库实例
func NewPluginInstanceRepository(db *gorm.DB) *PluginInstanceRepository {
	return &PluginInstanceRepository{db: db}
}

// Create 创建插件实例
func (r *PluginInstanceRepository) Create(instance *types.PluginInstance) error {
	return r.db.Create(instance).Error
}

// GetByID 根据ID获取插件实例
func (r *PluginInstanceRepository) GetByID(id string) (*types.PluginInstance, error) {
	var instance types.PluginInstance
	if err := r.db.Where("id = ?", id).First(&instance).Error; err != nil {
		return nil, err
	}
	return &instance, nil
}

// GetByName 根据名称获取插件实例
func (r *PluginInstanceRepository) GetByName(name string) (*types.PluginInstance, error) {
	var instance types.PluginInstance
	if err := r.db.Where("name = ?", name).First(&instance).Error; err != nil {
		return nil, err
	}
	return &instance, nil
}

// List 获取插件实例列表
func (r *PluginInstanceRepository) List(offset, limit int) ([]*types.PluginInstance, int64, error) {
	var instances []*types.PluginInstance
	var total int64

	if err := r.db.Model(&types.PluginInstance{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&instances).Error; err != nil {
		return nil, 0, err
	}

	return instances, total, nil
}

// Update 更新插件实例
func (r *PluginInstanceRepository) Update(instance *types.PluginInstance) error {
	return r.db.Save(instance).Error
}

// Delete 删除插件实例
func (r *PluginInstanceRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&types.PluginInstance{}).Error
}