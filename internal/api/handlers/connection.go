package handlers

import (
	"net/http"
	"strconv"

	"github.com/fustgo/fustgo2/internal/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ConnectionHandler 连接配置API处理程序
type ConnectionHandler struct {
	connectionService *service.ConnectionService
	logger            *zap.Logger
}

// NewConnectionHandler 创建连接配置API处理程序实例
func NewConnectionHandler(connectionService *service.ConnectionService, logger *zap.Logger) *ConnectionHandler {
	return &ConnectionHandler{
		connectionService: connectionService,
		logger:            logger,
	}
}

// Create 创建连接配置
// @Summary 创建连接配置
// @Description 创建新的数据源连接配置
// @Tags connections
// @Accept json
// @Produce json
// @Param request body service.CreateConnectionRequest true "创建连接配置请求"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/connections [post]
func (h *ConnectionHandler) Create(c *gin.Context) {
	var req service.CreateConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	connection, err := h.connectionService.Create(&req)
	if err != nil {
		h.logger.Error("创建连接配置失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "创建连接配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "创建连接配置成功",
		Data:    connection,
	})
}

// GetByID 根据ID获取连接配置
// @Summary 获取连接配置详情
// @Description 根据ID获取连接配置详情
// @Tags connections
// @Accept json
// @Produce json
// @Param id path string true "连接配置ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/connections/{id} [get]
func (h *ConnectionHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "连接配置ID不能为空",
		})
		return
	}

	connection, err := h.connectionService.GetByID(id)
	if err != nil {
		h.logger.Error("获取连接配置失败", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取连接配置失败: " + err.Error(),
		})
		return
	}

	if connection == nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "连接配置不存在",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取连接配置成功",
		Data:    connection,
	})
}

// GetByName 根据名称获取连接配置
// @Summary 根据名称获取连接配置
// @Description 根据名称获取连接配置
// @Tags connections
// @Accept json
// @Produce json
// @Param name path string true "连接配置名称"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/connections/name/{name} [get]
func (h *ConnectionHandler) GetByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "连接配置名称不能为空",
		})
		return
	}

	connection, err := h.connectionService.GetByName(name)
	if err != nil {
		h.logger.Error("获取连接配置失败", zap.String("name", name), zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取连接配置失败: " + err.Error(),
		})
		return
	}

	if connection == nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "连接配置不存在",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取连接配置成功",
		Data:    connection,
	})
}

// List 获取连接配置列表
// @Summary 获取连接配置列表
// @Description 获取连接配置列表
// @Tags connections
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/connections [get]
func (h *ConnectionHandler) List(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	connections, total, err := h.connectionService.List(page, pageSize)
	if err != nil {
		h.logger.Error("获取连接配置列表失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取连接配置列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取连接配置列表成功",
		Data: map[string]interface{}{
			"items": connections,
			"total": total,
			"page":  page,
			"size":  pageSize,
		},
	})
}

// Update 更新连接配置
// @Summary 更新连接配置
// @Description 更新连接配置信息
// @Tags connections
// @Accept json
// @Produce json
// @Param id path string true "连接配置ID"
// @Param request body service.UpdateConnectionRequest true "更新连接配置请求"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/connections/{id} [put]
func (h *ConnectionHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "连接配置ID不能为空",
		})
		return
	}

	var req service.UpdateConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	if err := h.connectionService.Update(id, &req); err != nil {
		h.logger.Error("更新连接配置失败", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "更新连接配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "更新连接配置成功",
	})
}

// Delete 删除连接配置
// @Summary 删除连接配置
// @Description 删除连接配置
// @Tags connections
// @Accept json
// @Produce json
// @Param id path string true "连接配置ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/connections/{id} [delete]
func (h *ConnectionHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "连接配置ID不能为空",
		})
		return
	}

	if err := h.connectionService.Delete(id); err != nil {
		h.logger.Error("删除连接配置失败", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "删除连接配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "删除连接配置成功",
	})
}

// Test 测试连接配置
// @Summary 测试连接配置
// @Description 测试连接配置是否有效
// @Tags connections
// @Accept json
// @Produce json
// @Param id path string true "连接配置ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/v1/connections/{id}/test [post]
func (h *ConnectionHandler) Test(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "连接配置ID不能为空",
		})
		return
	}

	response, err := h.connectionService.Test(id)
	if err != nil {
		h.logger.Error("测试连接配置失败", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "测试连接配置失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "测试连接配置完成",
		Data:    response,
	})
}