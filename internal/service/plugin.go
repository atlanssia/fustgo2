package service

import (
	"github.com/fustgo/fustgo2/internal/core/types"
	"github.com/fustgo/fustgo2/internal/repository"
	"github.com/google/uuid"
	"time"
)

// PluginService 插件服务
type PluginService struct {
	pluginRepo *repository.PluginRepository
}

// NewPluginService 创建插件服务实例
func NewPluginService(pluginRepo *repository.PluginRepository) *PluginService {
	return &PluginService{pluginRepo: pluginRepo}
}

// RegisterPlugin 注册插件
func (s *PluginService) RegisterPlugin(req *RegisterPluginRequest) (*types.Plugin, error) {
	plugin := &types.Plugin{
		ID:           uuid.New().String(),
		Name:         req.Name,
		Type:         req.Type,
		Version:      req.Version,
		Description:  req.Description,
		Author:       req.Author,
		ConfigSchema: req.ConfigSchema,
		Enabled:      req.Enabled,
		BuiltIn:      req.BuiltIn,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.pluginRepo.Create(plugin); err != nil {
		return nil, err
	}

	return plugin, nil
}

// GetPluginByID 根据ID获取插件
func (s *PluginService) GetPluginByID(id string) (*types.Plugin, error) {
	return s.pluginRepo.GetByID(id)
}

// ListPlugins 获取插件列表
func (s *PluginService) ListPlugins(page, pageSize int) ([]*types.Plugin, int64, error) {
	return s.pluginRepo.List((page-1)*pageSize, pageSize)
}

// UpdatePlugin 更新插件
func (s *PluginService) UpdatePlugin(id string, req *UpdatePluginRequest) error {
	plugin, err := s.pluginRepo.GetByID(id)
	if err != nil {
		return err
	}

	plugin.Name = req.Name
	plugin.Type = req.Type
	plugin.Version = req.Version
	plugin.Description = req.Description
	plugin.Author = req.Author
	plugin.ConfigSchema = req.ConfigSchema
	plugin.Enabled = req.Enabled
	plugin.UpdatedAt = time.Now()

	return s.pluginRepo.Update(plugin)
}

// DeletePlugin 删除插件
func (s *PluginService) DeletePlugin(id string) error {
	return s.pluginRepo.Delete(id)
}

// RegisterPluginRequest 注册插件请求
type RegisterPluginRequest struct {
	Name         string            `json:"name" validate:"required"`
	Type         types.PluginType  `json:"type" validate:"required"`
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	Author       string            `json:"author"`
	ConfigSchema string            `json:"config_schema"`
	Enabled      bool              `json:"enabled"`
	BuiltIn      bool              `json:"built_in"`
}

// UpdatePluginRequest 更新插件请求
type UpdatePluginRequest struct {
	Name         string            `json:"name" validate:"required"`
	Type         types.PluginType  `json:"type" validate:"required"`
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	Author       string            `json:"author"`
	ConfigSchema string            `json:"config_schema"`
	Enabled      bool              `json:"enabled"`
}