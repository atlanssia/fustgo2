package handlers

import (
	"net/http"
	"strconv"

	"github.com/fustgo/fustgo2/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// JobHandler 任务API处理程序
type JobHandler struct {
	jobService *service.JobService
	logger     *zap.Logger
}

// NewJobHandler 创建任务API处理程序实例
func NewJobHandler(jobService *service.JobService, logger *zap.Logger) *JobHandler {
	return &JobHandler{
		jobService: jobService,
		logger:     logger,
	}
}

// Create 创建任务
// @Summary 创建任务
// @Description 创建新的数据同步任务
// @Tags jobs
// @Accept json
// @Produce json
// @Param request body service.CreateJobRequest true "创建任务请求"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/jobs [post]
func (h *JobHandler) Create(c *gin.Context) {
	var req service.CreateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	job, err := h.jobService.Create(&req)
	if err != nil {
		h.logger.Error("创建任务失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "创建任务失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "创建任务成功",
		Data:    job,
	})
}

// GetByID 根据ID获取任务
// @Summary 获取任务详情
// @Description 根据ID获取任务详情
// @Tags jobs
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/jobs/{id} [get]
func (h *JobHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "任务ID不能为空",
		})
		return
	}

	job, err := h.jobService.GetByID(id)
	if err != nil {
		h.logger.Error("获取任务失败", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取任务失败: " + err.Error(),
		})
		return
	}

	if job == nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "任务不存在",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取任务成功",
		Data:    job,
	})
}

// List 获取任务列表
// @Summary 获取任务列表
// @Description 获取任务列表
// @Tags jobs
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/jobs [get]
func (h *JobHandler) List(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	jobs, total, err := h.jobService.List(page, pageSize)
	if err != nil {
		h.logger.Error("获取任务列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取任务列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取任务列表成功",
		Data: map[string]interface{}{
			"items": jobs,
			"total": total,
			"page":  page,
			"size":  pageSize,
		},
	})
}

// Update 更新任务
// @Summary 更新任务
// @Description 更新任务信息
// @Tags jobs
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Param request body service.UpdateJobRequest true "更新任务请求"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/jobs/{id} [put]
func (h *JobHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "任务ID不能为空",
		})
		return
	}

	var req service.UpdateJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	if err := h.jobService.Update(id, &req); err != nil {
		h.logger.Error("更新任务失败", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "更新任务失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "更新任务成功",
	})
}

// Delete 删除任务
// @Summary 删除任务
// @Description 删除任务
// @Tags jobs
// @Accept json
// @Produce json
// @Param id path string true "任务ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/jobs/{id} [delete]
func (h *JobHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "任务ID不能为空",
		})
		return
	}

	if err := h.jobService.Delete(id); err != nil {
		h.logger.Error("删除任务失败", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "删除任务失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "删除任务成功",
	})
}