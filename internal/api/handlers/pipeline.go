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

// PipelineHandler 管道API处理程序
type PipelineHandler struct {
	pipelineService   *service.PipelineService
	pipelineRepo      *repository.PipelineRepository
	pluginInstanceRepo *repository.PluginInstanceRepository
	logger            *zap.Logger
}

// NewPipelineHandler 创建管道API处理程序实例
func NewPipelineHandler(
	pipelineService *service.PipelineService,
	pipelineRepo *repository.PipelineRepository,
	pluginInstanceRepo *repository.PluginInstanceRepository,
	logger *zap.Logger,
) *PipelineHandler {
	return &PipelineHandler{
		pipelineService:   pipelineService,
		pipelineRepo:      pipelineRepo,
		pluginInstanceRepo: pluginInstanceRepo,
		logger:            logger,
	}
}

// Create 创建管道
// @Summary 创建管道
// @Description 创建新的数据管道
// @Tags pipelines
// @Accept json
// @Produce json
// @Param request body CreatePipelineRequest true "创建管道请求"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/pipelines [post]
func (h *PipelineHandler) Create(c *gin.Context) {
	var req CreatePipelineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 验证插件实例是否存在
	if req.Config.Source != nil {
		_, err := h.pluginInstanceRepo.GetByID(req.Config.Source.InstanceID)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Code:    400,
				Message: "数据源插件实例不存在: " + req.Config.Source.InstanceID,
			})
			return
		}
	}

	for _, processor := range req.Config.Processors {
		_, err := h.pluginInstanceRepo.GetByID(processor.InstanceID)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Code:    400,
				Message: "处理器插件实例不存在: " + processor.InstanceID,
			})
			return
		}
	}

	if req.Config.Sink != nil {
		_, err := h.pluginInstanceRepo.GetByID(req.Config.Sink.InstanceID)
		if err != nil {
			c.JSON(http.StatusBadRequest, Response{
				Code:    400,
				Message: "数据目标插件实例不存在: " + req.Config.Sink.InstanceID,
			})
			return
		}
	}

	pipeline := &types.Pipeline{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Status:      types.PipelineStatusDraft,
		Config:      "", // 实际实现中需要序列化配置
		Schedule:    "", // 实际实现中需要序列化调度配置
		Tags:        req.Tags,
		Enabled:     req.Enabled,
		CreatedBy:   req.CreatedBy,
	}

	if err := h.pipelineRepo.Create(pipeline); err != nil {
		h.logger.Error("创建管道失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "创建管道失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "创建管道成功",
		Data:    pipeline,
	})
}

// GetByID 根据ID获取管道
// @Summary 获取管道详情
// @Description 根据ID获取管道详情
// @Tags pipelines
// @Accept json
// @Produce json
// @Param id path string true "管道ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/pipelines/{id} [get]
func (h *PipelineHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "管道ID不能为空",
		})
		return
	}

	pipeline, err := h.pipelineRepo.GetByID(id)
	if err != nil {
		h.logger.Error("获取管道失败", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取管道失败: " + err.Error(),
		})
		return
	}

	if pipeline == nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "管道不存在",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取管道成功",
		Data:    pipeline,
	})
}

// List 获取管道列表
// @Summary 获取管道列表
// @Description 获取管道列表
// @Tags pipelines
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/pipelines [get]
func (h *PipelineHandler) List(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	pipelines, total, err := h.pipelineRepo.List((page-1)*pageSize, pageSize)
	if err != nil {
		h.logger.Error("获取管道列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取管道列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取管道列表成功",
		Data: map[string]interface{}{
			"items": pipelines,
			"total": total,
			"page":  page,
			"size":  pageSize,
		},
	})
}

// Update 更新管道
// @Summary 更新管道
// @Description 更新管道信息
// @Tags pipelines
// @Accept json
// @Produce json
// @Param id path string true "管道ID"
// @Param request body UpdatePipelineRequest true "更新管道请求"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/pipelines/{id} [put]
func (h *PipelineHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "管道ID不能为空",
		})
		return
	}

	var req UpdatePipelineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	pipeline, err := h.pipelineRepo.GetByID(id)
	if err != nil {
		h.logger.Error("获取管道失败", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取管道失败: " + err.Error(),
		})
		return
	}

	if pipeline == nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "管道不存在",
		})
		return
	}

	pipeline.Name = req.Name
	pipeline.Description = req.Description
	pipeline.Tags = req.Tags
	pipeline.Enabled = req.Enabled

	if err := h.pipelineRepo.Update(pipeline); err != nil {
		h.logger.Error("更新管道失败", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "更新管道失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "更新管道成功",
	})
}

// Delete 删除管道
// @Summary 删除管道
// @Description 删除管道
// @Tags pipelines
// @Accept json
// @Produce json
// @Param id path string true "管道ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/pipelines/{id} [delete]
func (h *PipelineHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "管道ID不能为空",
		})
		return
	}

	if err := h.pipelineRepo.Delete(id); err != nil {
		h.logger.Error("删除管道失败", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "删除管道失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "删除管道成功",
	})
}

// Execute 执行管道
// @Summary 执行管道
// @Description 执行指定的管道
// @Tags pipelines
// @Accept json
// @Produce json
// @Param id path string true "管道ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/pipelines/{id}/execute [post]
func (h *PipelineHandler) Execute(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "管道ID不能为空",
		})
		return
	}

	// TODO: 实现管道执行逻辑
	// 这里应该调用Pipeline引擎来执行管道

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "管道执行已启动",
		Data:    map[string]interface{}{"pipeline_id": id},
	})
}

// CreatePipelineRequest 创建管道请求
type CreatePipelineRequest struct {
	Name        string              `json:"name" binding:"required"`
	Description string              `json:"description"`
	Config      types.PipelineConfig `json:"config" binding:"required"`
	Tags        string              `json:"tags"`
	Enabled     bool                `json:"enabled"`
	CreatedBy   string              `json:"created_by"`
}

// UpdatePipelineRequest 更新管道请求
type UpdatePipelineRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Tags        string `json:"tags"`
	Enabled     bool   `json:"enabled"`
}