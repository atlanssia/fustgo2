package visual

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fustgo/fustgo2/internal/core/types"
	"github.com/fustgo/fustgo2/internal/repository"
	"github.com/fustgo/fustgo2/pkg/plugin"
	"go.uber.org/zap"
)

// ConfigVisualizer 可视化配置器
type ConfigVisualizer struct {
	pipelineRepo *repository.PipelineRepository
	pluginRepo   *repository.PluginRepository
	logger       *zap.Logger
}

// NewConfigVisualizer 创建可视化配置器实例
func NewConfigVisualizer(
	pipelineRepo *repository.PipelineRepository,
	pluginRepo *repository.PluginRepository,
	logger *zap.Logger,
) *ConfigVisualizer {
	return &ConfigVisualizer{
		pipelineRepo: pipelineRepo,
		pluginRepo:   pluginRepo,
		logger:       logger,
	}
}

// PipelineGraph 管道图结构
type PipelineGraph struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Nodes       []*GraphNode  `json:"nodes"`
	Edges       []*GraphEdge  `json:"edges"`
	Layout      string        `json:"layout"`
	CreatedAt   int64         `json:"created_at"`
	UpdatedAt   int64         `json:"updated_at"`
}

// GraphNode 图节点
type GraphNode struct {
	ID       string            `json:"id"`
	Type     string            `json:"type"` // source, processor, sink
	Name     string            `json:"name"`
	Plugin   string            `json:"plugin"`
	Config   map[string]interface{} `json:"config"`
	Position *NodePosition     `json:"position"`
	Status   string            `json:"status"`
}

// GraphEdge 图边
type GraphEdge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
	Label  string `json:"label"`
}

// NodePosition 节点位置
type NodePosition struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// GeneratePipelineGraph 生成管道图
func (v *ConfigVisualizer) GeneratePipelineGraph(ctx context.Context, pipelineID string) (*PipelineGraph, error) {
	// 获取管道配置
	pipeline, err := v.pipelineRepo.GetByID(pipelineID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pipeline: %w", err)
	}

	if pipeline == nil {
		return nil, fmt.Errorf("pipeline not found: %s", pipelineID)
	}

	// 解析管道配置
	var config types.PipelineConfig
	if err := json.Unmarshal([]byte(pipeline.Config), &config); err != nil {
		return nil, fmt.Errorf("failed to parse pipeline config: %w", err)
	}

	// 创建图结构
	graph := &PipelineGraph{
		ID:        pipeline.ID,
		Name:      pipeline.Name,
		Nodes:     make([]*GraphNode, 0),
		Edges:     make([]*GraphEdge, 0),
		Layout:    "horizontal",
		CreatedAt: pipeline.CreatedAt.Unix(),
		UpdatedAt: pipeline.UpdatedAt.Unix(),
	}

	// 添加数据源节点
	if config.Source != nil {
		node, err := v.createSourceNode(config.Source)
		if err != nil {
			return nil, fmt.Errorf("failed to create source node: %w", err)
		}
		graph.Nodes = append(graph.Nodes, node)
	}

	// 添加处理器节点
	for i, proc := range config.Processors {
		node, err := v.createProcessorNode(proc, i)
		if err != nil {
			return nil, fmt.Errorf("failed to create processor node: %w", err)
		}
		graph.Nodes = append(graph.Nodes, node)
	}

	// 添加数据目标节点
	if config.Sink != nil {
		node, err := v.createSinkNode(config.Sink)
		if err != nil {
			return nil, fmt.Errorf("failed to create sink node: %w", err)
		}
		graph.Nodes = append(graph.Nodes, node)
	}

	// 创建连接边
	if err := v.createEdges(graph); err != nil {
		return nil, fmt.Errorf("failed to create edges: %w", err)
	}

	// 设置节点位置
	v.layoutNodes(graph)

	return graph, nil
}

// createSourceNode 创建数据源节点
func (v *ConfigVisualizer) createSourceNode(ref *types.PluginInstanceRef) (*GraphNode, error) {
	// 获取插件实例信息
	// 这里简化处理，实际应该查询数据库获取详细信息
	node := &GraphNode{
		ID:     fmt.Sprintf("source_%s", ref.InstanceID),
		Type:   "source",
		Name:   "Data Source",
		Plugin: "unknown",
		Config: make(map[string]interface{}),
		Status: "active",
	}

	// 解析配置
	if ref.Config != "" {
		var config map[string]interface{}
		if err := json.Unmarshal([]byte(ref.Config), &config); err != nil {
			v.logger.Warn("Failed to parse source config", zap.Error(err))
		} else {
			node.Config = config
		}
	}

	return node, nil
}

// createProcessorNode 创建处理器节点
func (v *ConfigVisualizer) createProcessorNode(ref *types.PluginInstanceRef, index int) (*GraphNode, error) {
	// 获取插件实例信息
	// 这里简化处理，实际应该查询数据库获取详细信息
	node := &GraphNode{
		ID:     fmt.Sprintf("processor_%d", index),
		Type:   "processor",
		Name:   fmt.Sprintf("Processor %d", index+1),
		Plugin: "unknown",
		Config: make(map[string]interface{}),
		Status: "active",
	}

	// 解析配置
	if ref.Config != "" {
		var config map[string]interface{}
		if err := json.Unmarshal([]byte(ref.Config), &config); err != nil {
			v.logger.Warn("Failed to parse processor config", zap.Error(err))
		} else {
			node.Config = config
		}
	}

	return node, nil
}

// createSinkNode 创建数据目标节点
func (v *ConfigVisualizer) createSinkNode(ref *types.PluginInstanceRef) (*GraphNode, error) {
	// 获取插件实例信息
	// 这里简化处理，实际应该查询数据库获取详细信息
	node := &GraphNode{
		ID:     fmt.Sprintf("sink_%s", ref.InstanceID),
		Type:   "sink",
		Name:   "Data Sink",
		Plugin: "unknown",
		Config: make(map[string]interface{}),
		Status: "active",
	}

	// 解析配置
	if ref.Config != "" {
		var config map[string]interface{}
		if err := json.Unmarshal([]byte(ref.Config), &config); err != nil {
			v.logger.Warn("Failed to parse sink config", zap.Error(err))
		} else {
			node.Config = config
		}
	}

	return node, nil
}

// createEdges 创建连接边
func (v *ConfigVisualizer) createEdges(graph *PipelineGraph) error {
	// 连接数据源到第一个处理器
	if len(graph.Nodes) > 1 {
		edge := &GraphEdge{
			ID:     fmt.Sprintf("edge_%s_to_%s", graph.Nodes[0].ID, graph.Nodes[1].ID),
			Source: graph.Nodes[0].ID,
			Target: graph.Nodes[1].ID,
			Label:  "data flow",
		}
		graph.Edges = append(graph.Edges, edge)
	}

	// 连接处理器之间
	for i := 1; i < len(graph.Nodes)-1; i++ {
		edge := &GraphEdge{
			ID:     fmt.Sprintf("edge_%s_to_%s", graph.Nodes[i].ID, graph.Nodes[i+1].ID),
			Source: graph.Nodes[i].ID,
			Target: graph.Nodes[i+1].ID,
			Label:  "data flow",
		}
		graph.Edges = append(graph.Edges, edge)
	}

	return nil
}

// layoutNodes 设置节点位置
func (v *ConfigVisualizer) layoutNodes(graph *PipelineGraph) {
	// 简单的水平布局
	nodeSpacing := 200.0
	for i, node := range graph.Nodes {
		node.Position = &NodePosition{
			X: float64(i) * nodeSpacing,
			Y: 100.0,
		}
	}
}

// GetPluginCatalog 获取插件目录
func (v *ConfigVisualizer) GetPluginCatalog(ctx context.Context) ([]*PluginInfo, error) {
	// 获取所有插件
	plugins, _, err := v.pluginRepo.List(0, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to list plugins: %w", err)
	}

	// 转换为插件信息
	catalog := make([]*PluginInfo, 0, len(plugins))
	for _, p := range plugins {
		info := &PluginInfo{
			ID:          p.ID,
			Name:        p.Name,
			Type:        string(p.Type),
			Version:     p.Version,
			Description: p.Description,
			Author:      p.Author,
			Enabled:     p.Enabled,
			BuiltIn:     p.BuiltIn,
		}
		catalog = append(catalog, info)
	}

	return catalog, nil
}

// PluginInfo 插件信息
type PluginInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Author      string `json:"author"`
	Enabled     bool   `json:"enabled"`
	BuiltIn     bool   `json:"built_in"`
}

// PipelineTemplate 管道模板
type PipelineTemplate struct {
	ID          string                   `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	Category    string                   `json:"category"`
	Config      *types.PipelineConfig    `json:"config"`
	Tags        []string                 `json:"tags"`
	CreatedAt   int64                    `json:"created_at"`
}

// GetPipelineTemplates 获取管道模板
func (v *ConfigVisualizer) GetPipelineTemplates(ctx context.Context) ([]*PipelineTemplate, error) {
	// 这里简化处理，实际应该从模板库获取
	templates := []*PipelineTemplate{
		{
			ID:          "mysql-to-postgresql",
			Name:        "MySQL to PostgreSQL",
			Description: "从MySQL同步数据到PostgreSQL",
			Category:    "database",
			Config: &types.PipelineConfig{
				Source: &types.PluginInstanceRef{
					InstanceID: "mysql-reader",
				},
				Sink: &types.PluginInstanceRef{
					InstanceID: "postgresql-writer",
				},
			},
			Tags: []string{"mysql", "postgresql", "sync"},
		},
		{
			ID:          "mongodb-to-elasticsearch",
			Name:        "MongoDB to Elasticsearch",
			Description: "从MongoDB同步数据到Elasticsearch",
			Category:    "search",
			Config: &types.PipelineConfig{
				Source: &types.PluginInstanceRef{
					InstanceID: "mongodb-reader",
				},
				Sink: &types.PluginInstanceRef{
					InstanceID: "elasticsearch-writer",
				},
			},
			Tags: []string{"mongodb", "elasticsearch", "search"},
		},
	}

	return templates, nil
}