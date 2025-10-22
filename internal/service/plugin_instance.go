package service

import (
	"github.com/fustgo/fustgo2/internal/core/types"
	"github.com/fustgo/fustgo2/internal/repository"
	"github.com/google/uuid"
	"time"
)

// PluginInstanceService 插件实例服务
type PluginInstanceService struct {
	pluginInstanceRepo *repository.PluginInstanceRepository
}

// NewPluginInstanceService 创建插件实例服务实例
func NewPluginInstanceService(pluginInstanceRepo *repository.PluginInstanceRepository) *PluginInstanceService {
	return &PluginInstanceService{pluginInstanceRepo: pluginInstanceRepo}
}

// CreatePluginInstance 创建插件实例
func (s *PluginInstanceService) CreatePluginInstance(req *CreatePluginInstanceRequest) (*types.PluginInstance, error) {
	instance := &types.PluginInstance{
		ID:          uuid.New().String(),
		Name:        req.Name,
		PluginID:    req.PluginID,
		Config:      req.Config,
		Description: req.Description,
		Tags:        req.Tags,
		Enabled:     req.Enabled,
		CreatedBy:   req.CreatedBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.pluginInstanceRepo.Create(instance); err != nil {
		return nil, err
	}

	return instance, nil
}

// GetPluginInstanceByID 根据ID获取插件实例
func (s *PluginInstanceService) GetPluginInstanceByID(id string) (*types.PluginInstance, error) {
	return s.pluginInstanceRepo.GetByID(id)
}

// GetPluginInstanceByName 根据名称获取插件实例
func (s *PluginInstanceService) GetPluginInstanceByName(name string) (*types.PluginInstance, error) {
	return s.pluginInstanceRepo.GetByName(name)
}

// ListPluginInstances 获取插件实例列表
func (s *PluginInstanceService) ListPluginInstances(page, pageSize int) ([]*types.PluginInstance, int64, error) {
	return s.pluginInstanceRepo.List((page-1)*pageSize, pageSize)
}

// UpdatePluginInstance 更新插件实例
func (s *PluginInstanceService) UpdatePluginInstance(id string, req *UpdatePluginInstanceRequest) error {
	instance, err := s.pluginInstanceRepo.GetByID(id)
	if err != nil {
		return err
	}

	instance.Name = req.Name
	instance.PluginID = req.PluginID
	instance.Config = req.Config
	instance.Description = req.Description
	instance.Tags = req.Tags
	instance.Enabled = req.Enabled
	instance.UpdatedAt = time.Now()

	return s.pluginInstanceRepo.Update(instance)
}

// DeletePluginInstance 删除插件实例
func (s *PluginInstanceService) DeletePluginInstance(id string) error {
	return s.pluginInstanceRepo.Delete(id)
}

// CreatePluginInstanceRequest 创建插件实例请求
type CreatePluginInstanceRequest struct {
	Name        string `json:"name" validate:"required"`
	PluginID    string `json:"plugin_id" validate:"required"`
	Config      string `json:"config"`
	Description string `json:"description"`
	Tags        string `json:"tags"`
	Enabled     bool   `json:"enabled"`
	CreatedBy   string `json:"created_by"`
}

// UpdatePluginInstanceRequest 更新插件实例请求
type UpdatePluginInstanceRequest struct {
	Name        string `json:"name" validate:"required"`
	PluginID    string `json:"plugin_id" validate:"required"`
	Config      string `json:"config"`
	Description string `json:"description"`
	Tags        string `json:"tags"`
	Enabled     bool   `json:"enabled"`
}