package handlers

import (
	"net/http"
	"strconv"

	"github.com/fustgo/fustgo2/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ExecutionHandler 执行记录API处理程序
type ExecutionHandler struct {
	executionService *service.ExecutionService
	logger           *zap.Logger
}

// NewExecutionHandler 创建执行记录API处理程序实例
func NewExecutionHandler(executionService *service.ExecutionService, logger *zap.Logger) *ExecutionHandler {
	return &ExecutionHandler{
		executionService: executionService,
		logger:           logger,
	}
}

// Create 创建执行记录
// @Summary 创建执行记录
// @Description 创建新的执行记录
// @Tags executions
// @Accept json
// @Produce json
// @Param request body service.CreateExecutionRequest true "创建执行记录请求"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/executions [post]
func (h *ExecutionHandler) Create(c *gin.Context) {
	var req service.CreateExecutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	execution, err := h.executionService.Create(&req)
	if err != nil {
		h.logger.Error("创建执行记录失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "创建执行记录失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "创建执行记录成功",
		Data:    execution,
	})
}

// GetByID 根据ID获取执行记录
// @Summary 获取执行记录详情
// @Description 根据ID获取执行记录详情
// @Tags executions
// @Accept json
// @Produce json
// @Param id path string true "执行记录ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/executions/{id} [get]
func (h *ExecutionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "执行记录ID不能为空",
		})
		return
	}

	execution, err := h.executionService.GetByID(id)
	if err != nil {
		h.logger.Error("获取执行记录失败", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取执行记录失败: " + err.Error(),
		})
		return
	}

	if execution == nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "执行记录不存在",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取执行记录成功",
		Data:    execution,
	})
}

// ListByJobID 根据任务ID获取执行记录列表
// @Summary 根据任务ID获取执行记录列表
// @Description 根据任务ID获取执行记录列表
// @Tags executions
// @Accept json
// @Produce json
// @Param job_id path string true "任务ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/executions/job/{job_id} [get]
func (h *ExecutionHandler) ListByJobID(c *gin.Context) {
	jobID := c.Param("job_id")
	if jobID == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "任务ID不能为空",
		})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	executions, total, err := h.executionService.ListByJobID(jobID, page, pageSize)
	if err != nil {
		h.logger.Error("获取执行记录列表失败", zap.String("job_id", jobID), zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取执行记录列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取执行记录列表成功",
		Data: map[string]interface{}{
			"items": executions,
			"total": total,
			"page":  page,
			"size":  pageSize,
		},
	})
}

// List 获取执行记录列表
// @Summary 获取执行记录列表
// @Description 获取执行记录列表
// @Tags executions
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/executions [get]
func (h *ExecutionHandler) List(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	executions, total, err := h.executionService.List(page, pageSize)
	if err != nil {
		h.logger.Error("获取执行记录列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取执行记录列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取执行记录列表成功",
		Data: map[string]interface{}{
			"items": executions,
			"total": total,
			"page":  page,
			"size":  pageSize,
		},
	})
}

// Update 更新执行记录
// @Summary 更新执行记录
// @Description 更新执行记录信息
// @Tags executions
// @Accept json
// @Produce json
// @Param id path string true "执行记录ID"
// @Param request body service.UpdateExecutionRequest true "更新执行记录请求"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/executions/{id} [put]
func (h *ExecutionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "执行记录ID不能为空",
		})
		return
	}

	var req service.UpdateExecutionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	if err := h.executionService.Update(id, &req); err != nil {
		h.logger.Error("更新执行记录失败", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "更新执行记录失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "更新执行记录成功",
	})
}