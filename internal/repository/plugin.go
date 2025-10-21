package repository

import (
	"github.com/fustgo/fustgo2/internal/core/types"
	"gorm.io/gorm"
)

// PluginRepository 插件仓库
type PluginRepository struct {
	db *gorm.DB
}

// NewPluginRepository 创建插件仓库实例
func NewPluginRepository(db *gorm.DB) *PluginRepository {
	return &PluginRepository{db: db}
}

// Create 创建插件
func (r *PluginRepository) Create(plugin *types.Plugin) error {
	return r.db.Create(plugin).Error
}

// GetByID 根据ID获取插件
func (r *PluginRepository) GetByID(id string) (*types.Plugin, error) {
	var plugin types.Plugin
	if err := r.db.Where("id = ?", id).First(&plugin).Error; err != nil {
		return nil, err
	}
	return &plugin, nil
}

// GetByName 根据名称获取插件
func (r *PluginRepository) GetByName(name string) (*types.Plugin, error) {
	var plugin types.Plugin
	if err := r.db.Where("name = ?", name).First(&plugin).Error; err != nil {
		return nil, err
	}
	return &plugin, nil
}

// List 获取插件列表
func (r *PluginRepository) List(offset, limit int) ([]*types.Plugin, int64, error) {
	var plugins []*types.Plugin
	var total int64

	if err := r.db.Model(&types.Plugin{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&plugins).Error; err != nil {
		return nil, 0, err
	}

	return plugins, total, nil
}

// Update 更新插件
func (r *PluginRepository) Update(plugin *types.Plugin) error {
	return r.db.Save(plugin).Error
}

// Delete 删除插件
func (r *PluginRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&types.Plugin{}).Error
}