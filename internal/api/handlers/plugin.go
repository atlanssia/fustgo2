package handlers

import (
	"net/http"
	"strconv"

	"github.com/fustgo/fustgo2/internal/core/types"
	"github.com/fustgo/fustgo2/internal/repository"
	"github.com/fustgo/fustgo2/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// PluginHandler 插件API处理程序
type PluginHandler struct {
	pluginService        *service.PluginService
	pluginInstanceService *service.PluginInstanceService
	pluginRepo           *repository.PluginRepository
	pluginInstanceRepo   *repository.PluginInstanceRepository
	logger               *zap.Logger
}

// NewPluginHandler 创建插件API处理程序实例
func NewPluginHandler(
	pluginService *service.PluginService,
	pluginInstanceService *service.PluginInstanceService,
	pluginRepo *repository.PluginRepository,
	pluginInstanceRepo *repository.PluginInstanceRepository,
	logger *zap.Logger,
) *PluginHandler {
	return &PluginHandler{
		pluginService:        pluginService,
		pluginInstanceService: pluginInstanceService,
		pluginRepo:           pluginRepo,
		pluginInstanceRepo:   pluginInstanceRepo,
		logger:               logger,
	}
}

// ListPlugins 获取插件列表
// @Summary 获取插件列表
// @Description 获取所有可用插件列表
// @Tags plugins
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/plugins [get]
func (h *PluginHandler) ListPlugins(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	plugins, total, err := h.pluginRepo.List((page-1)*pageSize, pageSize)
	if err != nil {
		h.logger.Error("获取插件列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取插件列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取插件列表成功",
		Data: map[string]interface{}{
			"items": plugins,
			"total": total,
			"page":  page,
			"size":  pageSize,
		},
	})
}

// GetPluginByID 根据ID获取插件
// @Summary 获取插件详情
// @Description 根据ID获取插件详情
// @Tags plugins
// @Accept json
// @Produce json
// @Param id path string true "插件ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/plugins/{id} [get]
func (h *PluginHandler) GetPluginByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "插件ID不能为空",
		})
		return
	}

	plugin, err := h.pluginRepo.GetByID(id)
	if err != nil {
		h.logger.Error("获取插件失败", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取插件失败: " + err.Error(),
		})
		return
	}

	if plugin == nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "插件不存在",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取插件成功",
		Data:    plugin,
	})
}

// ListPluginInstances 获取插件实例列表
// @Summary 获取插件实例列表
// @Description 获取插件实例列表
// @Tags plugin_instances
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/plugin-instances [get]
func (h *PluginHandler) ListPluginInstances(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	instances, total, err := h.pluginInstanceRepo.List((page-1)*pageSize, pageSize)
	if err != nil {
		h.logger.Error("获取插件实例列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取插件实例列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取插件实例列表成功",
		Data: map[string]interface{}{
			"items": instances,
			"total": total,
			"page":  page,
			"size":  pageSize,
		},
	})
}

// CreatePluginInstance 创建插件实例
// @Summary 创建插件实例
// @Description 创建新的插件实例
// @Tags plugin_instances
// @Accept json
// @Produce json
// @Param request body CreatePluginInstanceRequest true "创建插件实例请求"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/plugin-instances [post]
func (h *PluginHandler) CreatePluginInstance(c *gin.Context) {
	var req CreatePluginInstanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	instance, err := h.pluginInstanceService.CreatePluginInstance(&req)
	if err != nil {
		h.logger.Error("创建插件实例失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "创建插件实例失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "创建插件实例成功",
		Data:    instance,
	})
}

// CreatePluginInstanceRequest 创建插件实例请求
type CreatePluginInstanceRequest struct {
	Name        string `json:"name" binding:"required"`
	PluginID    string `json:"plugin_id" binding:"required"`
	Config      string `json:"config"`
	Description string `json:"description"`
	Tags        string `json:"tags"`
	Enabled     bool   `json:"enabled"`
	CreatedBy   string `json:"created_by"`
}