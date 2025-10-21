package distributed

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/fustgo/fustgo2/internal/core/types"
	"github.com/fustgo/fustgo2/pkg/pipeline"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

// Worker 分布式工作节点
type Worker struct {
	id         string
	natsConn   *nats.Conn
	logger     *zap.Logger
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	pipelines  map[string]*pipeline.Pipeline
	pipelineMu sync.RWMutex
}

// NewWorker 创建工作节点实例
func NewWorker(id string, natsURL string, logger *zap.Logger) (*Worker, error) {
	// 连接NATS
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Worker{
		id:        id,
		natsConn:  nc,
		logger:    logger,
		ctx:       ctx,
		cancel:    cancel,
		pipelines: make(map[string]*pipeline.Pipeline),
	}, nil
}

// Start 启动工作节点
func (w *Worker) Start() error {
	w.logger.Info("Starting worker", zap.String("worker_id", w.id))

	// 订阅任务分配主题
	_, err := w.natsConn.Subscribe("pipeline.execute", w.handlePipelineExecution)
	if err != nil {
		return fmt.Errorf("failed to subscribe to pipeline execution topic: %w", err)
	}

	// 订阅心跳主题
	_, err = w.natsConn.Subscribe("worker.heartbeat", w.handleHeartbeat)
	if err != nil {
		return fmt.Errorf("failed to subscribe to heartbeat topic: %w", err)
	}

	// 启动心跳发送
	w.wg.Add(1)
	go w.sendHeartbeat()

	w.logger.Info("Worker started", zap.String("worker_id", w.id))
	return nil
}

// Stop 停止工作节点
func (w *Worker) Stop() error {
	w.logger.Info("Stopping worker", zap.String("worker_id", w.id))

	// 取消上下文
	w.cancel()

	// 等待所有goroutine完成
	w.wg.Wait()

	// 关闭NATS连接
	if w.natsConn != nil {
		w.natsConn.Close()
	}

	w.logger.Info("Worker stopped", zap.String("worker_id", w.id))
	return nil
}

// handlePipelineExecution 处理管道执行请求
func (w *Worker) handlePipelineExecution(msg *nats.Msg) {
	w.logger.Debug("Received pipeline execution request")

	// 解析请求
	var req PipelineExecutionRequest
	if err := decodeMsg(msg.Data, &req); err != nil {
		w.logger.Error("Failed to decode pipeline execution request", zap.Error(err))
		return
	}

	// 执行管道
	if err := w.executePipeline(&req); err != nil {
		w.logger.Error("Failed to execute pipeline", 
			zap.String("pipeline_id", req.PipelineID), 
			zap.Error(err))
		
		// 发送执行失败响应
		resp := &PipelineExecutionResponse{
			PipelineID: req.PipelineID,
			WorkerID:   w.id,
			Success:    false,
			Error:      err.Error(),
		}
		
		if err := w.sendResponse("pipeline.execution.result", resp); err != nil {
			w.logger.Error("Failed to send execution result", zap.Error(err))
		}
		return
	}

	// 发送执行成功响应
	resp := &PipelineExecutionResponse{
		PipelineID: req.PipelineID,
		WorkerID:   w.id,
		Success:    true,
	}
	
	if err := w.sendResponse("pipeline.execution.result", resp); err != nil {
		w.logger.Error("Failed to send execution result", zap.Error(err))
	}
}

// executePipeline 执行管道
func (w *Worker) executePipeline(req *PipelineExecutionRequest) error {
	w.logger.Info("Executing pipeline", 
		zap.String("pipeline_id", req.PipelineID),
		zap.String("worker_id", w.id))

	// 创建管道执行器
	// 这里简化处理，实际应该从配置创建管道
	executor := pipeline.NewExecutor(nil)

	// 启动执行
	if err := executor.Start(); err != nil {
		return fmt.Errorf("failed to start pipeline execution: %w", err)
	}

	// 等待执行完成
	if err := executor.Wait(); err != nil {
		return fmt.Errorf("pipeline execution failed: %w", err)
	}

	w.logger.Info("Pipeline executed successfully", 
		zap.String("pipeline_id", req.PipelineID),
		zap.String("worker_id", w.id))

	return nil
}

// handleHeartbeat 处理心跳请求
func (w *Worker) handleHeartbeat(msg *nats.Msg) {
	w.logger.Debug("Received heartbeat request")

	// 发送心跳响应
	resp := &HeartbeatResponse{
		WorkerID:  w.id,
		Timestamp: time.Now(),
		Status:    "active",
	}
	
	if err := w.sendResponse("worker.heartbeat.response", resp); err != nil {
		w.logger.Error("Failed to send heartbeat response", zap.Error(err))
	}
}

// sendHeartbeat 发送心跳
func (w *Worker) sendHeartbeat() {
	defer w.wg.Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-w.ctx.Done():
			return
		case <-ticker.C:
			heartbeat := &Heartbeat{
				WorkerID:  w.id,
				Timestamp: time.Now(),
				Status:    "active",
			}
			
			if err := w.sendMsg("worker.heartbeat", heartbeat); err != nil {
				w.logger.Error("Failed to send heartbeat", zap.Error(err))
			}
		}
	}
}

// sendMsg 发送消息
func (w *Worker) sendMsg(subject string, data interface{}) error {
	encoded, err := encodeMsg(data)
	if err != nil {
		return fmt.Errorf("failed to encode message: %w", err)
	}
	
	return w.natsConn.Publish(subject, encoded)
}

// sendResponse 发送响应
func (w *Worker) sendResponse(subject string, data interface{}) error {
	encoded, err := encodeMsg(data)
	if err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}
	
	return w.natsConn.Publish(subject, encoded)
}

// encodeMsg 编码消息
func encodeMsg(data interface{}) ([]byte, error) {
	// TODO: 实现消息编码（JSON序列化等）
	return nil, nil
}

// decodeMsg 解码消息
func decodeMsg(data []byte, v interface{}) error {
	// TODO: 实现消息解码（JSON反序列化等）
	return nil
}

// PipelineExecutionRequest 管道执行请求
type PipelineExecutionRequest struct {
	PipelineID string        `json:"pipeline_id"`
	Config     *types.PipelineConfig `json:"config"`
}

// PipelineExecutionResponse 管道执行响应
type PipelineExecutionResponse struct {
	PipelineID string `json:"pipeline_id"`
	WorkerID   string `json:"worker_id"`
	Success    bool   `json:"success"`
	Error      string `json:"error,omitempty"`
}

// Heartbeat 心跳消息
type Heartbeat struct {
	WorkerID  string    `json:"worker_id"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}

// HeartbeatResponse 心跳响应
type HeartbeatResponse struct {
	WorkerID  string    `json:"worker_id"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"`
}