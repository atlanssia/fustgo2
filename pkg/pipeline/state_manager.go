package pipeline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// StateManager 状态管理器
type StateManager struct {
	config *Config
	mu     sync.RWMutex
}

// NewStateManager 创建新的状态管理器
func NewStateManager(config *Config) *StateManager {
	// 确保检查点目录存在
	if config.EnableCheckpoint {
		os.MkdirAll(config.CheckpointPath, 0755)
	}
	
	return &StateManager{
		config: config,
	}
}

// SaveCheckpoint 保存检查点
func (sm *StateManager) SaveCheckpoint(checkpoint *Checkpoint) error {
	if !sm.config.EnableCheckpoint {
		return nil
	}
	
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	// 创建检查点文件路径
	filename := filepath.Join(sm.config.CheckpointPath, fmt.Sprintf("checkpoint_%s.json", checkpoint.PipelineID))
	
	// 序列化检查点
	data, err := json.MarshalIndent(checkpoint, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal checkpoint: %w", err)
	}
	
	// 写入文件
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write checkpoint file: %w", err)
	}
	
	return nil
}

// LoadCheckpoint 加载检查点
func (sm *StateManager) LoadCheckpoint(pipelineID string) (*Checkpoint, error) {
	if !sm.config.EnableCheckpoint {
		return nil, nil
	}
	
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	
	// 检查点文件路径
	filename := filepath.Join(sm.config.CheckpointPath, fmt.Sprintf("checkpoint_%s.json", pipelineID))
	
	// 检查文件是否存在
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil, nil
	}
	
	// 读取文件
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read checkpoint file: %w", err)
	}
	
	// 反序列化
	var checkpoint Checkpoint
	if err := json.Unmarshal(data, &checkpoint); err != nil {
		return nil, fmt.Errorf("failed to unmarshal checkpoint: %w", err)
	}
	
	return &checkpoint, nil
}

// ClearCheckpoint 清除检查点
func (sm *StateManager) ClearCheckpoint(pipelineID string) error {
	if !sm.config.EnableCheckpoint {
		return nil
	}
	
	sm.mu.Lock()
	defer sm.mu.Unlock()
	
	// 检查点文件路径
	filename := filepath.Join(sm.config.CheckpointPath, fmt.Sprintf("checkpoint_%s.json", pipelineID))
	
	// 删除文件
	if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove checkpoint file: %w", err)
	}
	
	return nil
}

// CreateCheckpoint 创建检查点
func (sm *StateManager) CreateCheckpoint(pipeline *Pipeline) *Checkpoint {
	return &Checkpoint{
		PipelineID:   pipeline.ID,
		Timestamp:    time.Now(),
		RecordsRead:  pipeline.Metrics.RecordsRead,
		ReaderState:  make(map[string]interface{}),
		WriterState:  make(map[string]interface{}),
	}
}
