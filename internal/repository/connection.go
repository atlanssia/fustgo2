package repository

import (
	"github.com/fustgo/fustgo2/internal/core/types"
	"gorm.io/gorm"
)

// ConnectionRepository 连接配置仓库
type ConnectionRepository struct {
	db *gorm.DB
}

// NewConnectionRepository 创建连接配置仓库实例
func NewConnectionRepository(db *gorm.DB) *ConnectionRepository {
	return &ConnectionRepository{db: db}
}

// Create 创建连接配置
func (r *ConnectionRepository) Create(connection *types.Connection) error {
	return r.db.Create(connection).Error
}

// GetByID 根据ID获取连接配置
func (r *ConnectionRepository) GetByID(id string) (*types.Connection, error) {
	var connection types.Connection
	if err := r.db.Where("id = ?", id).First(&connection).Error; err != nil {
		return nil, err
	}
	return &connection, nil
}

// GetByName 根据名称获取连接配置
func (r *ConnectionRepository) GetByName(name string) (*types.Connection, error) {
	var connection types.Connection
	if err := r.db.Where("name = ?", name).First(&connection).Error; err != nil {
		return nil, err
	}
	return &connection, nil
}

// List 获取连接配置列表
func (r *ConnectionRepository) List(offset, limit int) ([]*types.Connection, error) {
	var connections []*types.Connection
	if err := r.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&connections).Error; err != nil {
		return nil, err
	}
	return connections, nil
}

// Update 更新连接配置
func (r *ConnectionRepository) Update(connection *types.Connection) error {
	return r.db.Save(connection).Error
}

// Delete 删除连接配置
func (r *ConnectionRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&types.Connection{}).Error
}

// Count 获取连接配置总数
func (r *ConnectionRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&types.Connection{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}