package distributed

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/fustgo/fustgo2/internal/core/types"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// Coordinator 分布式协调器
type Coordinator struct {
	id        string
	natsConn  *nats.Conn
	logger    *zap.Logger
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	workers   map[string]*WorkerInfo
	workerMu  sync.RWMutex
	pipelines map[string]*PipelineInfo
	pipelineMu sync.RWMutex
}

// WorkerInfo 工作节点信息
type WorkerInfo struct {
	ID        string    `json:"id"`
	LastHeartbeat time.Time `json:"last_heartbeat"`
	Status    string    `json:"status"`
}

// PipelineInfo 管道信息
type PipelineInfo struct {
	ID       string    `json:"id"`
	WorkerID string    `json:"worker_id"`
	Status   string    `json:"status"`
	Started  time.Time `json:"started"`
}

// NewCoordinator 创建协调器实例
func NewCoordinator(id string, natsURL string, logger *zap.Logger) (*Coordinator, error) {
	// 连接NATS
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Coordinator{
		id:        id,
		natsConn:  nc,
		logger:    logger,
		ctx:       ctx,
		cancel:    cancel,
		workers:   make(map[string]*WorkerInfo),
		pipelines: make(map[string]*PipelineInfo),
	}, nil
}

// Start 启动协调器
func (c *Coordinator) Start() error {
	c.logger.Info("Starting coordinator", zap.String("coordinator_id", c.id))

	// 订阅工作节点心跳响应
	_, err := c.natsConn.Subscribe("worker.heartbeat.response", c.handleWorkerHeartbeat)
	if err != nil {
		return fmt.Errorf("failed to subscribe to worker heartbeat topic: %w", err)
	}

	// 订阅管道执行结果
	_, err = c.natsConn.Subscribe("pipeline.execution.result", c.handlePipelineResult)
	if err != nil {
		return fmt.Errorf("failed to subscribe to pipeline execution result topic: %w", err)
	}

	// 启动工作节点监控
	c.wg.Add(1)
	go c.monitorWorkers()

	c.logger.Info("Coordinator started", zap.String("coordinator_id", c.id))
	return nil
}

// Stop 停止协调器
func (c *Coordinator) Stop() error {
	c.logger.Info("Stopping coordinator", zap.String("coordinator_id", c.id))

	// 取消上下文
	c.cancel()

	// 等待所有goroutine完成
	c.wg.Wait()

	// 关闭NATS连接
	if c.natsConn != nil {
		c.natsConn.Close()
	}

	c.logger.Info("Coordinator stopped", zap.String("coordinator_id", c.id))
	return nil
}

// handleWorkerHeartbeat 处理工作节点心跳
func (c *Coordinator) handleWorkerHeartbeat(msg *nats.Msg) {
	var heartbeat HeartbeatResponse
	if err := decodeMsg(msg.Data, &heartbeat); err != nil {
		c.logger.Error("Failed to decode worker heartbeat", zap.Error(err))
		return
	}

	// 更新工作节点信息
	c.workerMu.Lock()
	c.workers[heartbeat.WorkerID] = &WorkerInfo{
		ID:            heartbeat.WorkerID,
		LastHeartbeat: heartbeat.Timestamp,
		Status:        heartbeat.Status,
	}
	c.workerMu.Unlock()

	c.logger.Debug("Received worker heartbeat", 
		zap.String("worker_id", heartbeat.WorkerID),
		zap.String("status", heartbeat.Status))
}

// handlePipelineResult 处理管道执行结果
func (c *Coordinator) handlePipelineResult(msg *nats.Msg) {
	var result PipelineExecutionResponse
	if err := decodeMsg(msg.Data, &result); err != nil {
		c.logger.Error("Failed to decode pipeline execution result", zap.Error(err))
		return
	}

	// 更新管道信息
	c.pipelineMu.Lock()
	if pipelineInfo, exists := c.pipelines[result.PipelineID]; exists {
		pipelineInfo.Status = "completed"
		if !result.Success {
			pipelineInfo.Status = "failed"
		}
	}
	c.pipelineMu.Unlock()

	if result.Success {
		c.logger.Info("Pipeline executed successfully", 
			zap.String("pipeline_id", result.PipelineID),
			zap.String("worker_id", result.WorkerID))
	} else {
		c.logger.Error("Pipeline execution failed", 
			zap.String("pipeline_id", result.PipelineID),
			zap.String("worker_id", result.WorkerID),
			zap.String("error", result.Error))
	}
}

// monitorWorkers 监控工作节点
func (c *Coordinator) monitorWorkers() {
	defer c.wg.Done()

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return
		case <-ticker.C:
			c.checkWorkerHealth()
		}
	}
}

// checkWorkerHealth 检查工作节点健康状态
func (c *Coordinator) checkWorkerHealth() {
	c.workerMu.Lock()
	defer c.workerMu.Unlock()

	now := time.Now()
	for id, worker := range c.workers {
		// 如果超过2分钟没有心跳，认为工作节点失效
		if now.Sub(worker.LastHeartbeat) > 2*time.Minute {
			c.logger.Warn("Worker is offline", 
				zap.String("worker_id", id),
				zap.Time("last_heartbeat", worker.LastHeartbeat))
			
			// 移除失效的工作节点
			delete(c.workers, id)
			
			// 重新分配该工作节点负责的管道
			c.reassignPipelines(id)
		}
	}
}

// reassignPipelines 重新分配管道
func (c *Coordinator) reassignPipelines(workerID string) {
	c.pipelineMu.Lock()
	defer c.pipelineMu.Unlock()

	for _, pipeline := range c.pipelines {
		if pipeline.WorkerID == workerID {
			// 选择新的工作节点
			newWorkerID := c.selectWorker()
			if newWorkerID != "" {
				pipeline.WorkerID = newWorkerID
				pipeline.Status = "reassigned"
				
				c.logger.Info("Reassigned pipeline to new worker", 
					zap.String("pipeline_id", pipeline.ID),
					zap.String("old_worker_id", workerID),
					zap.String("new_worker_id", newWorkerID))
			} else {
				pipeline.Status = "pending"
				c.logger.Warn("No available worker to reassign pipeline", 
					zap.String("pipeline_id", pipeline.ID))
			}
		}
	}
}

// selectWorker 选择工作节点
func (c *Coordinator) selectWorker() string {
	c.workerMu.RLock()
	defer c.workerMu.RUnlock()

	// 简单的轮询选择算法
	for id, worker := range c.workers {
		if worker.Status == "active" {
			return id
		}
	}
	
	return ""
}

// ExecutePipeline 执行管道
func (c *Coordinator) ExecutePipeline(pipeline *types.Pipeline) error {
	c.logger.Info("Executing pipeline", zap.String("pipeline_id", pipeline.ID))

	// 选择工作节点
	workerID := c.selectWorker()
	if workerID == "" {
		return fmt.Errorf("no available worker to execute pipeline")
	}

	// 记录管道信息
	c.pipelineMu.Lock()
	c.pipelines[pipeline.ID] = &PipelineInfo{
		ID:       pipeline.ID,
		WorkerID: workerID,
		Status:   "assigned",
		Started:  time.Now(),
	}
	c.pipelineMu.Unlock()

	// 发送执行请求到工作节点
	req := &PipelineExecutionRequest{
		PipelineID: pipeline.ID,
		Config:     nil, // 简化处理
	}
	
	subject := fmt.Sprintf("pipeline.execute.%s", workerID)
	if err := c.sendMsg(subject, req); err != nil {
		return fmt.Errorf("failed to send pipeline execution request: %w", err)
	}

	c.logger.Info("Pipeline execution request sent", 
		zap.String("pipeline_id", pipeline.ID),
		zap.String("worker_id", workerID))

	return nil
}

// GetWorkerStatus 获取工作节点状态
func (c *Coordinator) GetWorkerStatus() map[string]*WorkerInfo {
	c.workerMu.RLock()
	defer c.workerMu.RUnlock()

	status := make(map[string]*WorkerInfo)
	for id, worker := range c.workers {
		status[id] = &WorkerInfo{
			ID:            worker.ID,
			LastHeartbeat: worker.LastHeartbeat,
			Status:        worker.Status,
		}
	}
	
	return status
}

// GetPipelineStatus 获取管道状态
func (c *Coordinator) GetPipelineStatus() map[string]*PipelineInfo {
	c.pipelineMu.RLock()
	defer c.pipelineMu.RUnlock()

	status := make(map[string]*PipelineInfo)
	for id, pipeline := range c.pipelines {
		status[id] = &PipelineInfo{
			ID:       pipeline.ID,
			WorkerID: pipeline.WorkerID,
			Status:   pipeline.Status,
			Started:  pipeline.Started,
		}
	}
	
	return status
}

// sendMsg 发送消息
func (c *Coordinator) sendMsg(subject string, data interface{}) error {
	encoded, err := encodeMsg(data)
	if err != nil {
		return fmt.Errorf("failed to encode message: %w", err)
	}
	
	return c.natsConn.Publish(subject, encoded)
}