package plugin

import (
	"fmt"
	"sync"
)

// Factory 插件工厂函数类型
type Factory func() Plugin

// Registry 插件注册表
type Registry struct {
	mu        sync.RWMutex
	readers   map[string]Factory
	writers   map[string]Factory
	processors map[string]Factory
	infos     map[string]*Info
}

// globalRegistry 全局注册表
var globalRegistry = NewRegistry()

// NewRegistry 创建新的注册表
func NewRegistry() *Registry {
	return &Registry{
		readers:    make(map[string]Factory),
		writers:    make(map[string]Factory),
		processors: make(map[string]Factory),
		infos:      make(map[string]*Info),
	}
}

// Register 注册插件
func (r *Registry) Register(name string, factory Factory, info *Info) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// 验证插件信息
	if info == nil {
		return fmt.Errorf("plugin info is nil")
	}

	// 根据类型注册到相应的 map
	switch info.Type {
	case TypeReader:
		if _, exists := r.readers[name]; exists {
			return fmt.Errorf("reader plugin %s already registered", name)
		}
		r.readers[name] = factory
	case TypeWriter:
		if _, exists := r.writers[name]; exists {
			return fmt.Errorf("writer plugin %s already registered", name)
		}
		r.writers[name] = factory
	case TypeProcessor:
		if _, exists := r.processors[name]; exists {
			return fmt.Errorf("processor plugin %s already registered", name)
		}
		r.processors[name] = factory
	default:
		return fmt.Errorf("unknown plugin type: %s", info.Type)
	}

	// 保存插件信息
	r.infos[name] = info

	return nil
}

// CreateReader 创建 Reader 插件实例
func (r *Registry) CreateReader(name string) (Reader, error) {
	r.mu.RLock()
	factory, ok := r.readers[name]
	r.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("reader plugin %s not found", name)
	}

	plugin := factory()
	reader, ok := plugin.(Reader)
	if !ok {
		return nil, fmt.Errorf("plugin %s is not a Reader", name)
	}

	return reader, nil
}

// CreateWriter 创建 Writer 插件实例
func (r *Registry) CreateWriter(name string) (Writer, error) {
	r.mu.RLock()
	factory, ok := r.writers[name]
	r.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("writer plugin %s not found", name)
	}

	plugin := factory()
	writer, ok := plugin.(Writer)
	if !ok {
		return nil, fmt.Errorf("plugin %s is not a Writer", name)
	}

	return writer, nil
}

// CreateProcessor 创建 Processor 插件实例
func (r *Registry) CreateProcessor(name string) (Processor, error) {
	r.mu.RLock()
	factory, ok := r.processors[name]
	r.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("processor plugin %s not found", name)
	}

	plugin := factory()
	processor, ok := plugin.(Processor)
	if !ok {
		return nil, fmt.Errorf("plugin %s is not a Processor", name)
	}

	return processor, nil
}

// GetInfo 获取插件信息
func (r *Registry) GetInfo(name string) (*Info, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	info, ok := r.infos[name]
	if !ok {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	return info, nil
}

// ListReaders 列出所有 Reader 插件
func (r *Registry) ListReaders() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.readers))
	for name := range r.readers {
		names = append(names, name)
	}
	return names
}

// ListWriters 列出所有 Writer 插件
func (r *Registry) ListWriters() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.writers))
	for name := range r.writers {
		names = append(names, name)
	}
	return names
}

// ListProcessors 列出所有 Processor 插件
func (r *Registry) ListProcessors() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.processors))
	for name := range r.processors {
		names = append(names, name)
	}
	return names
}

// ListAll 列出所有插件信息
func (r *Registry) ListAll() []*Info {
	r.mu.RLock()
	defer r.mu.RUnlock()

	infos := make([]*Info, 0, len(r.infos))
	for _, info := range r.infos {
		infos = append(infos, info)
	}
	return infos
}

// Clear 清空注册表（主要用于测试）
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.readers = make(map[string]Factory)
	r.writers = make(map[string]Factory)
	r.processors = make(map[string]Factory)
	r.infos = make(map[string]*Info)
}

// 全局注册表的便捷方法

// Register 注册插件到全局注册表
func Register(name string, factory Factory, info *Info) error {
	return globalRegistry.Register(name, factory, info)
}

// CreateReader 从全局注册表创建 Reader
func CreateReader(name string) (Reader, error) {
	return globalRegistry.CreateReader(name)
}

// CreateWriter 从全局注册表创建 Writer
func CreateWriter(name string) (Writer, error) {
	return globalRegistry.CreateWriter(name)
}

// CreateProcessor 从全局注册表创建 Processor
func CreateProcessor(name string) (Processor, error) {
	return globalRegistry.CreateProcessor(name)
}

// GetInfo 从全局注册表获取插件信息
func GetInfo(name string) (*Info, error) {
	return globalRegistry.GetInfo(name)
}

// ListReaders 列出全局注册表中的所有 Reader
func ListReaders() []string {
	return globalRegistry.ListReaders()
}

// ListWriters 列出全局注册表中的所有 Writer
func ListWriters() []string {
	return globalRegistry.ListWriters()
}

// ListProcessors 列出全局注册表中的所有 Processor
func ListProcessors() []string {
	return globalRegistry.ListProcessors()
}

// ListAll 列出全局注册表中的所有插件
func ListAll() []*Info {
	return globalRegistry.ListAll()
}

// GetGlobalRegistry 获取全局注册表（主要用于测试）
func GetGlobalRegistry() *Registry {
	return globalRegistry
}
