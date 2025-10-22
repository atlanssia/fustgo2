package enterprise

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/fustgo/fustgo2/internal/core/types"
	"github.com/fustgo/fustgo2/internal/repository"
	"go.uber.org/zap"
)

// TenantManager 多租户管理器
type TenantManager struct {
	tenantRepo    *repository.TenantRepository
	logger        *zap.Logger
	tenants       map[string]*Tenant
	tenantMu      sync.RWMutex
	defaultLimits *ResourceLimits
}

// Tenant 租户信息
type Tenant struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Enabled     bool            `json:"enabled"`
	Limits      *ResourceLimits `json:"limits"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// ResourceLimits 资源限制
type ResourceLimits struct {
	MaxJobs            int64 `json:"max_jobs"`
	MaxConcurrentJobs  int64 `json:"max_concurrent_jobs"`
	MaxPipelines       int64 `json:"max_pipelines"`
	MaxDataVolume      int64 `json:"max_data_volume"`      // MB
	MaxExecutionTime   int64 `json:"max_execution_time"`   // 秒
	MaxUsers           int64 `json:"max_users"`
	MaxConnections     int64 `json:"max_connections"`
}

// NewTenantManager 创建多租户管理器实例
func NewTenantManager(
	tenantRepo *repository.TenantRepository,
	defaultLimits *ResourceLimits,
	logger *zap.Logger,
) *TenantManager {
	return &TenantManager{
		tenantRepo:    tenantRepo,
		logger:        logger,
		tenants:       make(map[string]*Tenant),
		defaultLimits: defaultLimits,
	}
}

// CreateTenant 创建租户
func (tm *TenantManager) CreateTenant(ctx context.Context, tenant *Tenant) error {
	// 设置默认限制
	if tenant.Limits == nil {
		tenant.Limits = tm.defaultLimits
	}

	// 保存到数据库
	if err := tm.tenantRepo.Create(tenant); err != nil {
		return fmt.Errorf("failed to create tenant in database: %w", err)
	}

	// 缓存租户信息
	tm.tenantMu.Lock()
	tm.tenants[tenant.ID] = tenant
	tm.tenantMu.Unlock()

	tm.logger.Info("Tenant created", zap.String("tenant_id", tenant.ID), zap.String("tenant_name", tenant.Name))
	return nil
}

// GetTenant 获取租户信息
func (tm *TenantManager) GetTenant(ctx context.Context, tenantID string) (*Tenant, error) {
	// 先从缓存获取
	tm.tenantMu.RLock()
	tenant, exists := tm.tenants[tenantID]
	tm.tenantMu.RUnlock()

	if exists {
		return tenant, nil
	}

	// 从数据库获取
	tenant, err := tm.tenantRepo.GetByID(tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant from database: %w", err)
	}

	if tenant == nil {
		return nil, fmt.Errorf("tenant not found: %s", tenantID)
	}

	// 缓存租户信息
	tm.tenantMu.Lock()
	tm.tenants[tenantID] = tenant
	tm.tenantMu.Unlock()

	return tenant, nil
}

// UpdateTenant 更新租户信息
func (tm *TenantManager) UpdateTenant(ctx context.Context, tenantID string, updates map[string]interface{}) error {
	// 更新数据库
	if err := tm.tenantRepo.Update(tenantID, updates); err != nil {
		return fmt.Errorf("failed to update tenant in database: %w", err)
	}

	// 更新缓存
	tm.tenantMu.Lock()
	if tenant, exists := tm.tenants[tenantID]; exists {
		// 应用更新
		for key, value := range updates {
			switch key {
			case "name":
				if name, ok := value.(string); ok {
					tenant.Name = name
				}
			case "description":
				if desc, ok := value.(string); ok {
					tenant.Description = desc
				}
			case "enabled":
				if enabled, ok := value.(bool); ok {
					tenant.Enabled = enabled
				}
			case "limits":
				if limits, ok := value.(*ResourceLimits); ok {
					tenant.Limits = limits
				}
			}
		}
		tenant.UpdatedAt = time.Now()
	}
	tm.tenantMu.Unlock()

	tm.logger.Info("Tenant updated", zap.String("tenant_id", tenantID))
	return nil
}

// DeleteTenant 删除租户
func (tm *TenantManager) DeleteTenant(ctx context.Context, tenantID string) error {
	// 从数据库删除
	if err := tm.tenantRepo.Delete(tenantID); err != nil {
		return fmt.Errorf("failed to delete tenant from database: %w", err)
	}

	// 从缓存删除
	tm.tenantMu.Lock()
	delete(tm.tenants, tenantID)
	tm.tenantMu.Unlock()

	tm.logger.Info("Tenant deleted", zap.String("tenant_id", tenantID))
	return nil
}

// ListTenants 列出租户
func (tm *TenantManager) ListTenants(ctx context.Context, offset, limit int64) ([]*Tenant, int64, error) {
	return tm.tenantRepo.List(offset, limit)
}

// CheckResourceLimits 检查资源限制
func (tm *TenantManager) CheckResourceLimits(ctx context.Context, tenantID string, resourceType string) (bool, error) {
	tenant, err := tm.GetTenant(ctx, tenantID)
	if err != nil {
		return false, fmt.Errorf("failed to get tenant: %w", err)
	}

	if !tenant.Enabled {
		return false, fmt.Errorf("tenant is disabled")
	}

	// 检查具体资源限制
	switch resourceType {
	case "jobs":
		return tm.checkJobLimit(ctx, tenant)
	case "concurrent_jobs":
		return tm.checkConcurrentJobLimit(ctx, tenant)
	case "pipelines":
		return tm.checkPipelineLimit(ctx, tenant)
	case "data_volume":
		return tm.checkDataVolumeLimit(ctx, tenant)
	default:
		return true, nil
	}
}

// checkJobLimit 检查任务数量限制
func (tm *TenantManager) checkJobLimit(ctx context.Context, tenant *Tenant) (bool, error) {
	// 这里简化处理，实际应该查询数据库统计当前任务数量
	currentJobs := int64(0) // TODO: 实际查询
	return currentJobs < tenant.Limits.MaxJobs, nil
}

// checkConcurrentJobLimit 检查并发任务限制
func (tm *TenantManager) checkConcurrentJobLimit(ctx context.Context, tenant *Tenant) (bool, error) {
	// 这里简化处理，实际应该查询数据库统计当前并发任务数量
	currentConcurrentJobs := int64(0) // TODO: 实际查询
	return currentConcurrentJobs < tenant.Limits.MaxConcurrentJobs, nil
}

// checkPipelineLimit 检查管道数量限制
func (tm *TenantManager) checkPipelineLimit(ctx context.Context, tenant *Tenant) (bool, error) {
	// 这里简化处理，实际应该查询数据库统计当前管道数量
	currentPipelines := int64(0) // TODO: 实际查询
	return currentPipelines < tenant.Limits.MaxPipelines, nil
}

// checkDataVolumeLimit 检查数据量限制
func (tm *TenantManager) checkDataVolumeLimit(ctx context.Context, tenant *Tenant) (bool, error) {
	// 这里简化处理，实际应该查询数据库统计当前数据量
	currentDataVolume := int64(0) // TODO: 实际查询
	return currentDataVolume < tenant.Limits.MaxDataVolume, nil
}

// TenantContext 租户上下文
type TenantContext struct {
	TenantID  string
	UserID    string
	IsAdmin   bool
	Resources map[string]interface{}
}

// WithTenantContext 添加租户上下文
func (tm *TenantManager) WithTenantContext(ctx context.Context, tenantCtx *TenantContext) context.Context {
	return context.WithValue(ctx, "tenant_context", tenantCtx)
}

// GetTenantFromContext 从上下文获取租户信息
func (tm *TenantManager) GetTenantFromContext(ctx context.Context) (*TenantContext, bool) {
	tenantCtx, ok := ctx.Value("tenant_context").(*TenantContext)
	return tenantCtx, ok
}

// IsolateResources 隔离资源
func (tm *TenantManager) IsolateResources(ctx context.Context, tenantID string, resourceType string, resourceID string) error {
	// 确保资源属于指定租户
	// 这里简化处理，实际应该在数据库查询时添加租户ID过滤条件
	tm.logger.Debug("Resource isolation",
		zap.String("tenant_id", tenantID),
		zap.String("resource_type", resourceType),
		zap.String("resource_id", resourceID))
	return nil
}

// GetDefaultLimits 获取默认资源限制
func (tm *TenantManager) GetDefaultLimits() *ResourceLimits {
	return tm.defaultLimits
}

// SetDefaultLimits 设置默认资源限制
func (tm *TenantManager) SetDefaultLimits(limits *ResourceLimits) {
	tm.defaultLimits = limits
}