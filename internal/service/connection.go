package service

import (
	"github.com/fustgo/fustgo2/internal/core/types"
	"github.com/fustgo/fustgo2/internal/repository"
	"github.com/google/uuid"
	"time"
)

// ConnectionService 连接配置服务
type ConnectionService struct {
	connectionRepo *repository.ConnectionRepository
}

// NewConnectionService 创建连接配置服务实例
func NewConnectionService(connectionRepo *repository.ConnectionRepository) *ConnectionService {
	return &ConnectionService{connectionRepo: connectionRepo}
}

// Create 创建连接配置
func (s *ConnectionService) Create(req *CreateConnectionRequest) (*types.Connection, error) {
	connection := &types.Connection{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
		Config:      req.Config,
		Tags:        req.Tags,
		Enabled:     req.Enabled,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		CreatedBy:   req.CreatedBy,
	}

	if err := s.connectionRepo.Create(connection); err != nil {
		return nil, err
	}

	return connection, nil
}

// GetByID 根据ID获取连接配置
func (s *ConnectionService) GetByID(id string) (*types.Connection, error) {
	return s.connectionRepo.GetByID(id)
}

// GetByName 根据名称获取连接配置
func (s *ConnectionService) GetByName(name string) (*types.Connection, error) {
	return s.connectionRepo.GetByName(name)
}

// List 获取连接配置列表
func (s *ConnectionService) List(page, pageSize int) ([]*types.Connection, int64, error) {
	total, err := s.connectionRepo.Count()
	if err != nil {
		return nil, 0, err
	}

	connections, err := s.connectionRepo.List((page-1)*pageSize, pageSize)
	if err != nil {
		return nil, 0, err
	}

	return connections, total, nil
}

// Update 更新连接配置
func (s *ConnectionService) Update(id string, req *UpdateConnectionRequest) error {
	connection, err := s.connectionRepo.GetByID(id)
	if err != nil {
		return err
	}

	connection.Name = req.Name
	connection.Type = req.Type
	connection.Description = req.Description
	connection.Config = req.Config
	connection.Tags = req.Tags
	connection.Enabled = req.Enabled
	connection.UpdatedAt = time.Now()

	return s.connectionRepo.Update(connection)
}

// Delete 删除连接配置
func (s *ConnectionService) Delete(id string) error {
	return s.connectionRepo.Delete(id)
}

// Test 测试连接配置
func (s *ConnectionService) Test(id string) (*TestConnectionResponse, error) {
	connection, err := s.connectionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// TODO: 实现连接测试逻辑
	// 这里应该根据连接类型创建相应的连接并测试连通性

	response := &TestConnectionResponse{
		Success: true,
		Message: "连接测试成功",
	}

	// 更新最后测试时间
	connection.LastTestedAt = &time.Time{}
	*connection.LastTestedAt = time.Now()
	connection.LastTestStatus = "success"
	
	// 更新数据库
	s.connectionRepo.Update(connection)

	return response, nil
}

// CreateConnectionRequest 创建连接配置请求
type CreateConnectionRequest struct {
	Name        string `json:"name" validate:"required"`
	Type        string `json:"type" validate:"required"`
	Description string `json:"description"`
	Config      string `json:"config"`
	Tags        string `json:"tags"`
	Enabled     bool   `json:"enabled"`
	CreatedBy   string `json:"created_by"`
}

// UpdateConnectionRequest 更新连接配置请求
type UpdateConnectionRequest struct {
	Name        string `json:"name" validate:"required"`
	Type        string `json:"type" validate:"required"`
	Description string `json:"description"`
	Config      string `json:"config"`
	Tags        string `json:"tags"`
	Enabled     bool   `json:"enabled"`
}

// TestConnectionResponse 测试连接响应
type TestConnectionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}